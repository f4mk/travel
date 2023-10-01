import { Suspense, useCallback } from 'react'
import { FormattedMessage } from 'react-intl'
import { useNavigate } from 'react-router-dom'
import { Button, Loader } from '@mantine/core'

import { useGetMe } from '#/api/user'
import logo from '#/assets/coggers.png'
import { CenteredLoader } from '#/components/ui/CenteredLoader'
import { RoundButton } from '#/components/ui/RoundButton'
import { ROUTES } from '#/constants/routes'
import { lazy } from '#/utils'

import * as S from './styled'

const { ProfileMenu } = lazy(() => import('#/components/ui/ProfileMenu'))

export const Header = () => {
  const navigate = useNavigate()
  const { data, isFetching } = useGetMe({
    onError: () => {
      navigate(ROUTES.ROOT)
    },
    retry: false,
  })

  console.log(data)
  const handleTabChange = useCallback(
    (path: string) => {
      navigate(path)
    },
    [navigate]
  )

  const handleLogoClick = useCallback(() => {
    navigate(ROUTES.APP.ROOT)
  }, [navigate])

  return isFetching ? (
    <CenteredLoader />
  ) : (
    <S.Header>
      <RoundButton onClick={handleLogoClick}>
        <S.Img src={logo} alt="logo" loading="lazy" />
      </RoundButton>

      <S.Tabs>
        <Button
          variant="subtle"
          onClick={() => handleTabChange(ROUTES.APP.ROOT)}
        >
          <FormattedMessage
            description="Home tab"
            defaultMessage="Home"
            id="BXXnPK"
          />
        </Button>
        <Button
          variant="subtle"
          onClick={() => handleTabChange(ROUTES.APP.MAP)}
        >
          <FormattedMessage
            description="Map tab"
            defaultMessage="Map"
            id="JrIBeU"
          />
        </Button>
        <Button
          variant="subtle"
          onClick={() => handleTabChange(ROUTES.APP.BLOG)}
        >
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
