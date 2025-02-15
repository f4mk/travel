import { FormattedMessage } from 'react-intl'
import { useNavigate } from 'react-router-dom'
import { Button, Text, Title } from '@mantine/core'

import { ROUTES } from '#/constants/routes'

import * as S from './styled'
export const ConfirmCreateUser = () => {
  const navigate = useNavigate()

  const handleRedirect = () => {
    const param = new URLSearchParams({
      auth: 'true',
    }).toString()
    navigate(`${ROUTES.ROOT}?${param}`)
  }

  return (
    <S.Div>
      <Title order={3}>
        <FormattedMessage
          description="User creation page need verification success message"
          defaultMessage="Almost there!"
          id="kjHp6z"
        />
      </Title>
      <Text fz={'lg'}>
        <FormattedMessage
          description="User creation page need verification message"
          defaultMessage={`Your account has been created.
            To activate it you need to follow the 
            link that we have sent to your email.
            If you didn't get the letter, try using password recovery option.`}
          id="DhpNph"
        />
      </Text>
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
