import logo from './logo.svg'
import type { HTMLProps } from './types'

export function Html({ children, styles, title }: HTMLProps) {
  return (
    <html lang="en">
      <head>
        <meta charSet="utf-8" />
        <meta name="viewport" content="width=device-width, initial-scale=1" />
        <link rel="icon" href={logo} />
        <title>{title}</title>
        {styles}
      </head>

      <body>
        <div id="react-root" dangerouslySetInnerHTML={{ __html: children }} />
      </body>
    </html>
  )
}
