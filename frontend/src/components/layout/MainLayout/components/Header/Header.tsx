import { useCallback } from 'react'
import { Button } from '@mantine/core'
import { navigate } from 'vite-plugin-ssr/client/router'

import logo from '#/assets/coggers.png'
import { ProfileMenu } from '#/components/ui/ProfileMenu'
import { RoundButton } from '#/components/ui/RoundButton'
import { ERoutes } from '#/constants/routes'

import * as S from './styled'

export const Header = () => {
  const handleTabChange = useCallback((path: ERoutes) => {
    navigate(path)
  }, [])

  const handleLogoClick = useCallback(() => {
    navigate(ERoutes.ROOT)
  }, [])

  return (
    <S.Header>
      <RoundButton onClick={handleLogoClick}>
        <S.Img src={logo} alt="logo" loading="lazy" />
      </RoundButton>

      <S.Tabs>
        <Button variant="subtle" onClick={() => handleTabChange(ERoutes.ROOT)}>
          Home
        </Button>
        <Button variant="subtle" onClick={() => handleTabChange(ERoutes.MAP)}>
          Map
        </Button>
        <Button variant="subtle" onClick={() => handleTabChange(ERoutes.BLOG)}>
          Blog
        </Button>
      </S.Tabs>
      <ProfileMenu isLoggedIn={false} />
    </S.Header>
  )
}
