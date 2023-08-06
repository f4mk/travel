import styled from '@emotion/styled'
import { keyframes } from '@mantine/core'

const spin = keyframes`
from {
  transform: rotate(0deg);
}
to {
  transform: rotate(360deg);
}
`
export const Header = styled('header')`
  padding: 8px 16px;
  grid-area: header;
  display: grid;
  align-items: center;
  grid-template-columns: min-content auto min-content;
`
export const Img = styled('img')`
  object-fit: contain;
  height: 40px;
  animation: ${spin} 24s linear infinite;
  cursor: pointer;
`
export const Tabs = styled('div')`
  display: flex;
  justify-content: center;
  column-gap: 16px;
`
