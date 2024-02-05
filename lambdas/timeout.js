const API_ENDPOINT = `${process.env.API_URL}/order`

exports.handler = async function(event, context) {
  console.log('received order timeout')
  
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

  try {
    const { order } = payload
    await sendUpdate({
      orderId: order.order_id,
      status: 'FAILED'
    })
    console.log('timeout complete')
  } catch (err) {
    console.error('failed to update order')
  }
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