import { useEffect } from 'react'
import { FormattedMessage } from 'react-intl'
import { useNavigate } from 'react-router-dom'
import { Button, Title } from '@mantine/core'

import { useVerifyUser } from '#/api/user'
import { CenteredLoader } from '#/components/ui/CenteredLoader'
import { ROUTES } from '#/constants/routes'
import { getUrlParams } from '#/utils'

import * as S from './styled'
export const VerifyAccount = () => {
  const navigate = useNavigate()
  const { email, token } = getUrlParams('email', 'token')
  const { mutate, isLoading, isError } = useVerifyUser()

  const handleRedirect = () => {
    const param = new URLSearchParams({
      auth: 'true',
    }).toString()
    navigate(`${ROUTES.ROOT}?${param}`)
  }

  useEffect(() => {
    mutate({ email, token })
    // eslint-disable-next-line
  }, [])
  if (isLoading) {
    return <CenteredLoader />
  }
  if (isError) {
    navigate(ROUTES.NOT_FOUND)
  }
  return (
    <S.Div>
      <Title order={3}>
        <FormattedMessage
          description="User verification page success message"
          defaultMessage="Success! Your account has been verified"
          id="yq6ZoR"
        />
      </Title>
      <Button variant="subtle" onClick={handleRedirect}>
        <FormattedMessage
          description="User verification page redirect button text"
          defaultMessage="Go to home page"
          id="o44TG+"
        />
      </Button>
    </S.Div>
  )
}
