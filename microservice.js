const Fastify = require('fastify')
const fastifyPostgres = require('@fastify/postgres')
const { SNSClient, PublishCommand } = require("@aws-sdk/client-sns")
const { orderSchema, orderPatchSchema } = require('./schemas/orderSchema')

const PORT = process.env.PORT || 3000

const fastify = Fastify({
  logger: true
})

fastify.register(fastifyPostgres, {
  connectionString: process.env.DATABASE_URL
})

const snsClient = new SNSClient({ endpoint: process.env.AWS_URL, tls: false })

fastify.post('/order', {
  schema: orderSchema,
  handler: async function (request, reply) {
    const { customer, total, products } = request.body

    try {
      // Save Order
      const order = await saveOrder({ customer, total, products })

      // Publish Transmission
      await publishTransmission({ order })

      reply.send({ message: 'Order Created', order_id: order.order_id })
    } catch (err) {
      request.log.error(err)
      reply.status(500).send('Unable to process order')
    }
  }
})

fastify.patch('/order', {
  schema: orderPatchSchema,
  handler: async function (request, reply) {
    const { orderId, status } = request.body

    try {
      // Update Order
      const order = await updateOrder({ orderId, status })

      reply.send({ message: 'Order Updated', order_id: order.order_id })
    } catch (err) {
      request.log.error(err)
      reply.status(500).send('Unable to update order')
    }
  }
})

fastify.listen({ host: '0.0.0.0', port: PORT }, function (err, address) {
  if (err) {
    fastify.log.error(err)
    process.exit(1)
  }
})

async function saveOrder({ customer, products, total }) {
  const query = `
    INSERT INTO orders (customer, products, status, total)
    VALUES ($1, $2, $3, $4)
    RETURNING *
  `
  const values = [customer, JSON.stringify(products), "PENDING", total]

  const { rows } = await fastify.pg.query(query, values)

  return rows[0]
}

async function updateOrder({ orderId, status }) {
  const query = `
    UPDATE orders SET status = $1
    where order_id = $2
    RETURNING *
  `
  const values = [status, orderId]

  const result = await fastify.pg.query(query, values)

  if (result.rowCount === 0) {
    throw new Error('order not found')
  }

  return result.rows[0]
}

async function publishTransmission({ order }) {
  const params = {
    Message: JSON.stringify({
      order
    }),
    TopicArn: process.env.ORDER_TRANSMISSION_TOPIC,
  }

  try {
    const response = await snsClient.send(
      new PublishCommand(params),
    )
    fastify.log.info(`Message ${response.MessageId} sent to the topic ${params.TopicArn}`)
  } catch (err) {
    fastify.log.error(err)
    throw err
  }
}