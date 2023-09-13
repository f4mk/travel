import { Helmet } from 'react-helmet-async'

export const Meta = () => {
  return (
    <Helmet>
      <meta charSet="utf-8" />
      <meta name="viewport" content="initial-scale=1, width=device-width" />
      <link
        rel="stylesheet"
        href="https://fonts.googleapis.com/css?family=Roboto:400,500,700&display=swap"
      />
      <link rel="icon" type="image/x-icon" href="/favicon.ico"></link>
      <title>Traillyst</title>
    </Helmet>
  )
}
