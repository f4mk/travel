import ReactDOM from 'react-dom/client'

import { App } from './App'
import type { PageContextClient } from './types'

// This render() hook only supports SSR, see https://vite-plugin-ssr.com/render-modes for how to modify render() to support SPA
let root: ReactDOM.Root
export const render = async (pageContext: PageContextClient) => {
  const { Page, pageProps } = pageContext

  const page = (
    <App pageContext={pageContext}>
      <Page {...pageProps} />
    </App>
  )

  const container = document.getElementById('react-root')
  if (!container) {
    throw new Error('No container found')
  }
  // SPA
  if (container.innerHTML === '' || !pageContext.isHydration) {
    if (!root) {
      root = ReactDOM.createRoot(container)
    }
    root.render(page)
    // SSR
  } else {
    root = ReactDOM.hydrateRoot(container, page)
  }
}

export const hydrationCanBeAborted = true
//To enable Client-side Routing:
export const clientRouting = true
