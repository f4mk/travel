import { Suspense, useCallback } from 'react'
import { FormattedMessage } from 'react-intl'
import { useNavigate } from 'react-router-dom'
import { Button, Loader } from '@mantine/core'

import { HttpError } from '#/api/request'
import logo from '#/assets/coggers.png'
import { RoundButton } from '#/components/ui/RoundButton'
import { ERoutes } from '#/constants/routes'
import { lazy } from '#/utils/lazy'

import { useData } from './queries'
import * as S from './styled'

const { ProfileMenu } = lazy(() => import('#/components/ui/ProfileMenu'))

export const Header = () => {
  const navigate = useNavigate()
  const handleError = (e: HttpError) => {
    if (e.response.status === 403) {
      navigate('/')
    }
  }
  const { name, email } = useData({ onError: handleError })
  console.log(name, email)
  const handleTabChange = useCallback(
    (path: ERoutes) => {
      navigate(path)
    },
    [navigate]
  )

  const handleLogoClick = useCallback(() => {
    navigate(ERoutes.ROOT)
  }, [navigate])

  return (
    <S.Header>
      <RoundButton onClick={handleLogoClick}>
        <S.Img src={logo} alt="logo" loading="lazy" />
      </RoundButton>

      <S.Tabs>
        <Button variant="subtle" onClick={() => handleTabChange(ERoutes.ROOT)}>
          <FormattedMessage
            description="Home tab"
            defaultMessage="Home"
            id="BXXnPK"
          />
        </Button>
        <Button variant="subtle" onClick={() => handleTabChange(ERoutes.MAP)}>
          <FormattedMessage
            description="Map tab"
            defaultMessage="Map"
            id="JrIBeU"
          />
        </Button>
        <Button variant="subtle" onClick={() => handleTabChange(ERoutes.BLOG)}>
          <FormattedMessage
            description="Blog tab"
            defaultMessage="Blog"
            id="Ym3hvK"
          />
        </Button>
      </S.Tabs>
      <Suspense fallback={<Loader />}>
        <ProfileMenu />
      </Suspense>
    </S.Header>
  )
}
