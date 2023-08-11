import styled from '@emotion/styled'

export const Div = styled.div<{ isPressable: boolean }>`
  cursor: ${(props) => (props.isPressable ? 'pointer' : 'auto')};
`
