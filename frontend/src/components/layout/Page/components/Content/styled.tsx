import styled from '@emotion/styled'

export const Main = styled('main')`
  overflow: scroll;
  display: grid;
  width: 100%;
  height: 100%;
  background-color: ${({ theme }) =>
    theme.colorScheme === 'dark' ? theme.colors.dark[6] : theme.colors.gray[0]};
`
