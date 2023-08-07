import Fastify from 'fastify'
import path from 'path'
import sirv from 'sirv'
import { renderPage } from 'vite-plugin-ssr/server'

const root = path.resolve()
const isProduction = process.env.NODE_ENV === 'production'

const startServer = async () => {
  const fastify = Fastify({
    logger: false
  })

  await fastify.register(import('@fastify/compress'), { global: true })

  await fastify.register(import('@fastify/early-hints'), {
    warn: true
  })

  if (isProduction) {
    const assets = sirv(`${root}/dist/client`, {
      immutable: true,
      dev: !isProduction
    })

    fastify.addHook('onRequest', (req, reply, done) => {
      assets(req.raw, reply.raw, done)
    })
  } else {
    const vite = await import('vite')
    const viteDevMiddleware = (
      await vite.createServer({
        root,
        server: { middlewareMode: true }
      })
    ).middlewares
    fastify.addHook('onRequest', (req, reply, done) => {
      viteDevMiddleware(req.raw, reply.raw, done)
    })
  }

  fastify.get('/*', async function handler(request, reply) {
    {
      const pageContextInit = {
        urlOriginal: request.url
      }
      const pageContext = await renderPage(pageContextInit)
      const { httpResponse } = pageContext

      if (!httpResponse) {
        return
      }
      const { body, statusCode, earlyHints, contentType } = httpResponse
      await reply.writeEarlyHints({
        Link: earlyHints.map((e) => e.earlyHintLink)
      })
      return reply.status(statusCode).type(contentType).send(body)
    }
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
