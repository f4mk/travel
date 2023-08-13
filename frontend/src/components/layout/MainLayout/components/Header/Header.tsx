import { Suspense, useCallback } from 'react'
import { FormattedMessage } from 'react-intl'
import { useNavigate } from 'react-router-dom'
import { Button, Loader } from '@mantine/core'

import logo from '#/assets/coggers.png'
import { RoundButton } from '#/components/ui/RoundButton'
import { ERoutes } from '#/constants/routes'
import { lazy } from '#/utils/lazy'

import * as S from './styled'

const { ProfileMenu } = lazy(() => import('#/components/ui/ProfileMenu'))

export const Header = () => {
  const navigate = useNavigate()

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
          <FormattedMessage description="Home tab" defaultMessage="Home" />
        </Button>
        <Button variant="subtle" onClick={() => handleTabChange(ERoutes.MAP)}>
          <FormattedMessage description="Map tab" defaultMessage="Map" />
        </Button>
        <Button variant="subtle" onClick={() => handleTabChange(ERoutes.BLOG)}>
          <FormattedMessage description="Blog tab" defaultMessage="Blog" />
        </Button>
      </S.Tabs>
      <Suspense fallback={<Loader />}>
        <ProfileMenu isLoggedIn={false} />
      </Suspense>
    </S.Header>
  )
}
