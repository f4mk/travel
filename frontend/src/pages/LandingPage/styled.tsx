import styled from '@emotion/styled'

export const Div = styled.div`
  display: grid;
  grid-template-columns: 60% 40%;
  grid-template-rows: 100%;
`
export const Main = styled.main`
  position: relative;
`
export const Img = styled.img`
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
  position: relative;
  top: 60%;
  width: 20%;
  margin: auto;
  display: flex;
  flex-direction: column;
  align-items: stretch;
  justify-content: center;
  gap: 8px;
`
export const Aside = styled.aside``
