import { FormattedMessage } from 'react-intl'
import { useNavigate } from 'react-router-dom'
import { Button, Title } from '@mantine/core'

import { ERoutes } from '#/constants/routes'

import * as S from './styled'
export const NotFoundPage = () => {
  const navigate = useNavigate()
  return (
    <S.Div>
      <Title order={1}>
        <FormattedMessage
          description="Not found page"
          defaultMessage="Page not found"
          id="NAkyhe"
        />
      </Title>
      <Button variant="subtle" onClick={() => navigate(ERoutes.ROOT)}>
        <FormattedMessage
          defaultMessage="Go home"
          description="Not found page redirect home button"
          id="Wu1iXO"
        />
      </Button>
    </S.Div>
  )
}
