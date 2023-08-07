import { useCallback } from 'react'
import { useNavigate } from 'react-router-dom'
import { Button } from '@mantine/core'

import logo from '#/assets/coggers.png'
import { ProfileMenu } from '#/components/ui/ProfileMenu'
import { RoundButton } from '#/components/ui/RoundButton'
import { ERoutes } from '#/constants/routes'

import * as S from './styled'

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
