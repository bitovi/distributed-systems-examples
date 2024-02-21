const Fastify = require('fastify')

const PORT = process.env.PORT || 3000

const fastify = Fastify({
  logger: true
})

fastify.register(require('fastify-graceful-shutdown'))

const transmissions = []

fastify.post('/', function (request, reply) {
  transmissions.push(request.body)
  console.log('received transmission', request.body)
  reply.send({ status: 'success' })
})

fastify.get('/', function (request, reply) {
  reply.send({ transmissions })
})

fastify.after(() => {
  fastify.gracefulShutdown((signal, next) => {
    console.log('Upps!')
    next()
  })
})

fastify.listen({ host: '0.0.0.0', port: PORT }, function (err, address) {
  if (err) {
    fastify.log.error(err)
    process.exit(1)
  }
})
