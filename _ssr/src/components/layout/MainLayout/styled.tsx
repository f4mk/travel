import styled from '@emotion/styled'

export const Section = styled('section')`
  display: grid;
  width: 100vw;
  height: 100vh;
  min-width: 768px;
  max-width: 1440px;
  grid-template-areas:
    'header  header  header  header'
    'content content content content'
    'footer  footer  footer  footer';
  grid-template-columns: 150px auto auto auto;
  grid-template-rows: min-content auto 24px;
`
