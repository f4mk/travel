import styled from '@emotion/styled'

export const Div = styled.div<{ isPressable: boolean }>`
  transform: rotate(20deg);
  cursor: ${(props) => (props.isPressable ? 'pointer' : 'auto')};
`
