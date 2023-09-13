import { FormattedMessage } from 'react-intl'
import { Button } from '@mantine/core'

import backgroundDefault from '#/assets/backgroundDefault.webp'

import * as S from './styled'
export const LandingPage = () => {
  return (
    <S.Div>
      <S.Main>
        <S.Img src={backgroundDefault} alt="background-picture" />
        <S.ButtonContainer>
          <Button variant="filled" onClick={() => console.log('register')}>
            <FormattedMessage
              description="Register button"
              defaultMessage="Join"
              id="qxTdKD"
            />
          </Button>
          <Button variant="default" onClick={() => console.log('login')}>
            <FormattedMessage
              description="Login button"
              defaultMessage="Sign In"
              id="Ww5Nr+"
            />
          </Button>
        </S.ButtonContainer>
      </S.Main>
      <S.Aside>ALLO</S.Aside>
    </S.Div>
  )
}
