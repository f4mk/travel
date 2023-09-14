import styled from '@emotion/styled'

export const Div = styled.div`
  display: grid;
  overflow-x: scroll;
  grid-template-rows: 100%;
  min-width: 1080px;
`
export const Main = styled.main`
  position: relative;
  display: grid;
  grid-template-rows: min-content min-content auto;
  justify-content: center;
  align-items: center;
  min-width: 768px;
`
export const Img = styled.img`
  z-index: -1;
  width: 100%;
  height: 100%;
  object-fit: cover;
  position: absolute;
  left: 0;
  right: 0;
  top: 0;
  bottom: 0;
`
export const ButtonContainer = styled.div`
  width: 20%;
  margin: auto;
  display: flex;
  flex-direction: column;
  align-items: stretch;
  justify-content: center;
  gap: 8px;
`
export const Sup = styled.h3`
  margin: 0;
  font-size: 36px;
  color: ${({ theme }) => theme.colors.gray[0]};
`
export const Title = styled.h1`
  margin: 0;
  font-size: 160px;
  color: ${({ theme }) => theme.colors.gray[0]};
`
export const Aside = styled.aside``
