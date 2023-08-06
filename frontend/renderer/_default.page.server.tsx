// See https://vite-plugin-ssr.com/data-fetching
export const passToClient = ['pageProps', 'urlPathname', 'urlOriginal']

import { renderToString } from 'react-dom/server'
import { createStylesServer, ServerStyles } from '@mantine/ssr'
import { dangerouslySkipEscape } from 'vite-plugin-ssr/server'

import { App } from './App'
import { Html } from './Html'
import type { PageContextServer } from './types'

const stylesServer = createStylesServer()

export const render = async (pageContext: PageContextServer) => {
  const { Page, pageProps } = pageContext
  // This render() hook only supports SSR, see https://vite-plugin-ssr.com/render-modes for how to modify render() to support SPA
  let pageHtml

  if (Page) {
    pageHtml = renderToString(
      <App pageContext={pageContext}>
        <Page {...pageProps} />
      </App>
    )
  } else {
    pageHtml = ''
  }

  // See https://vite-plugin-ssr.com/head
  const { documentProps } = pageContext.exports
  const title = (documentProps && documentProps.title) || 'Vite SSR app'

  // <Html /> is already sanitized by React
  const documentHtml = dangerouslySkipEscape(
    renderToString(
      <Html
        title={title}
        styles={<ServerStyles html={pageHtml} server={stylesServer} />}
      >
        {pageHtml}
      </Html>
    )
  )

  return {
    documentHtml,
    pageContext: {
      // We can add some `pageContext` here, which is useful if we want to do page redirection https://vite-plugin-ssr.com/page-redirection
    }
  }
}
