import backgroundDefault from '#/assets/backgroundDefault.webp'

import * as S from './styled'
export const LandingPage = () => {
  return (
    <S.Div>
      <S.Main>
        <S.Img src={backgroundDefault} alt="background-picture" />
      </S.Main>
      <S.Aside>ALLO</S.Aside>
    </S.Div>
  )
}
