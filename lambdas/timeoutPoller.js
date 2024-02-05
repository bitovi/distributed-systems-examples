const { Pool } = require('pg')
const { SNSClient, PublishCommand } = require("@aws-sdk/client-sns")

const pool = new Pool({
  user: process.env.POSTGRES_USER,
  host: process.env.POSTGRES_HOST,
  password: process.env.POSTGRES_PASS,
  port: process.env.POSTGRES_PORT,
})
const snsClient = new SNSClient({
  endpoint: process.env.AWS_URL,
  tls: false
})

exports.handler = async function(event, context) {
  console.log('timeout poller running', {
    user: process.env.POSTGRES_USER,
    pass: process.env.POSTGRES_PASS
  })
  const client = await pool.connect()
  try {
    const res = await client.query(`
      SELECT * from orders
      WHERE created_at < NOW() - INTERVAL '5 minutes'
      AND status = 'PENDING';
    `)
    await Promise.all(res.rows.map(order => {
      publishMessage({ order })
    }))
    console.log('Success')
  } finally {
    client.release()
  }
}


async function publishMessage({ order }) {
  const params = {
    Message: JSON.stringify({
      order
    }),
    TopicArn: process.env.ORDER_TIMEOUT_TOPIC,
  }

  try {
    console.log(`Publishing message for order - ${order.order_id}`)
    const response = await snsClient.send(
      new PublishCommand(params),
    )
    console.log(`Message ${response.MessageId} sent to the topic ${params.TopicArn}`)
  } catch (err) {
    console.error(err)
    throw err
  }
}