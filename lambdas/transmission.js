const { SQSClient, ChangeMessageVisibilityCommand } = require("@aws-sdk/client-sqs")
const TRANSMISSION_ENDPOINT = 'http://transmission-endpoint:3000'
const API_ENDPOINT = `${process.env.API_URL}/order`
const MAX_RETRIES = 3

const client = new SQSClient({ endpoint: "http://localstack:4566", tls: false })

exports.handler = async function(event, context) {
  console.log('received order tranmission')
  const sqsMessage = event.Records[0]
  const body = JSON.parse(sqsMessage?.body)

  let payload
  try {
    payload = JSON.parse(body?.Message)
  }
  catch (err){
    console.error(`Unable to parse message body`, err)
    throw err
  }

  let retryCount = (sqsMessage.attributes || {}).ApproximateReceiveCount
  console.log({ retryCount })

  if (retryCount >= MAX_RETRIES) {
    console.log('Max Retries reached')
    return
  }

  try {
    const transmission = await sendTransmission(payload)

    switch (transmission.status) {
      case 'SUCCESS':
        console.log('Transmission success')
        const { order } = payload
        await sendUpdate({
          orderId: order.order_id,
          status: 'CONFIRMED'
        })
        break
      case 'ERROR':
        await reQueueMessageWithDelay(sqsMessage.receiptHandle, 240)
        break
    }
  } catch (err) {
    if (!(err instanceof DelayError)) {
      await reQueueMessageWithDelay(sqsMessage.receiptHandle, 240)
    }
  }
  console.log('END')
}

async function sendTransmission(body) {
  try {
    const response = await fetch(TRANSMISSION_ENDPOINT, {
      method: 'POST',
      headers: {
        'content-type': 'application/json'
      },
      body: JSON.stringify(body)
    })
    if (!response.ok) {
      return {
        status: 'ERROR',
        message: response.statusText
      }
    }

    return {
      status: 'SUCCESS',
      data: await response.json()
    }
  } catch (err) {
    console.error(err)
    return {
      status: 'ERROR',
      message: err.message
    }
  }
}

class DelayError extends Error {}

async function reQueueMessageWithDelay(handle, delaySeconds) {
  console.log('reQueueMessageWithDelay', { handle })
  const params = {
      QueueUrl: 'http://sqs.us-east-1.localhost.localstack.cloud:4566/000000000000/orderTransmission',
      ReceiptHandle: handle,
      VisibilityTimeout: delaySeconds,
  }
  try {
    const command = new ChangeMessageVisibilityCommand(params)
    const res = await client.send(command)
    console.log('Message re-queued with delay:', delaySeconds)
    console.log(JSON.stringify(res))
  } catch (err) {
    console.error(`failed to change visibility - ${err}`)
    throw err
  }
  throw new DelayError('Lets go fly a kite')
}

async function sendUpdate(body) {
  try {
    console.log('sendUpdate', JSON.stringify(body))
    const response = await fetch(API_ENDPOINT, {
      method: 'PATCH',
      headers: {
        'content-type': 'application/json'
      },
      body: JSON.stringify(body)
    })
    if (!response.ok) {
      return {
        status: 'ERROR',
        message: response.statusText
      }
    }

    return {
      status: 'SUCCESS',
      data: await response.json()
    }
  } catch (err) {
    console.error(err)
    return {
      status: 'ERROR',
      message: err.message
    }
  }
}