const TRANSMISSION_ENDPOINT = 'http://transmission-endpoint:3000'

exports.handler = async function(event, context) {
  console.log('received order tranmission')

  return fetch(TRANSMISSION_ENDPOINT, {
    method: 'POST',
  })
}
