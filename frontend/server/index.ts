import fastifyStatic from '@fastify/static'
import Fastify from 'fastify'
import fs from 'fs'
import path from 'path'
const startServer = async () => {
  const fastify = Fastify({
    logger: false
  })

  const indexPath = path.resolve('./dist/index.html')
  const index = fs.readFileSync(indexPath, 'utf8')

  fastify.register(fastifyStatic, {
    root: path.resolve('./dist/assets'),
    prefix: '/assets/'
  })

  fastify.get('*', (_, reply) => {
    reply.type('text/html').send(index)
  })

  try {
    // eslint-disable-next-line
    console.log('starting server on port: 3000')
    await fastify.listen({ port: 3000, host: '0.0.0.0' })
  } catch (err) {
    fastify.log.error(err)
    process.exit(1)
  }
}
startServer()
