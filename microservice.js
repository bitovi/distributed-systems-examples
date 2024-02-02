const Fastify = require('fastify')
const fastifyPostgres = require('@fastify/postgres')
const { SNSClient, PublishCommand } = require("@aws-sdk/client-sns")
const { orderSchema } = require('./schemas/orderSchema')

const PORT = process.env.PORT || 3000

const fastify = Fastify({
  logger: true
})
// Replace these with your actual database connection details
fastify.register(fastifyPostgres, {
  connectionString: process.env.DATABASE_URL
})

const snsClient = new SNSClient({ endpoint: "http://localstack:4566", tls: false })

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

async function publishTransmission({ order }) {
  const params = {
    Message: JSON.stringify({
      order
    }),
    TopicArn: 'arn:aws:sns:us-east-1:000000000000:orderTransmission',
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