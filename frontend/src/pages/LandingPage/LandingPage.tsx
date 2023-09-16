import { useEffect } from 'react'
import { FormattedMessage } from 'react-intl'
import { useSearchParams } from 'react-router-dom'
import { Button } from '@mantine/core'

import backgroundDefault from '#/assets/backgroundDefault.webp'
import { useModal } from '#/components/layout/ModalProvider'
import { EFormView } from '#/components/ui/Auth'
import { lazy } from '#/utils'

import * as S from './styled'

const { Auth } = lazy(() => import('#/components/ui/Auth'))

export const LandingPage = () => {
  const { showModal, hideModal } = useModal()
  const [searchParams, setSearchParams] = useSearchParams()
  const handleOpen = (view: EFormView) => {
    showModal(<Auth activeTab={view} onClose={hideModal} />)
  }

  useEffect(() => {
    const auth = searchParams.get('auth')
    if (auth) {
      handleOpen(EFormView.SIGN_IN)
      searchParams.delete('auth')
      setSearchParams(searchParams)
    }
    // eslint-disable-next-line
  }, [])

  return (
    <S.Div>
      <S.Main>
        <S.Sup>
          <FormattedMessage
            description="Landing title slogan"
            defaultMessage="Leave trails. Keep memories"
            id="+Pny/y"
          />
        </S.Sup>
        <S.Title>
          <FormattedMessage
            description="Landing title"
            defaultMessage="Traillyst"
            id="feLASf"
          />
        </S.Title>
        <S.Img src={backgroundDefault} alt="background-picture" />
        <S.ButtonContainer>
          <Button
            variant="filled"
            onClick={() => handleOpen(EFormView.SIGN_UP)}
          >
            <FormattedMessage
              description="Register button"
              defaultMessage="Join"
              id="qxTdKD"
            />
          </Button>
          <Button
            variant="default"
            onClick={() => handleOpen(EFormView.SIGN_IN)}
          >
            <FormattedMessage
              description="Login button"
              defaultMessage="Sign In"
              id="Ww5Nr+"
            />
          </Button>
        </S.ButtonContainer>
      </S.Main>
    </S.Div>
  )
}
