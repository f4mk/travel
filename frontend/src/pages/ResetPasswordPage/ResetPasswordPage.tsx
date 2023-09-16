import { FormattedMessage, useIntl } from 'react-intl'
import { useNavigate } from 'react-router-dom'
import { Button, Group, PasswordInput, Space, Title } from '@mantine/core'
import { useForm } from '@mantine/form'

import { usePasswordResetSubmit } from '#/api/auth'
import { PASSWORD_MIN_LENGTH } from '#/components/ui/RegisterForm'
import { ERoutes } from '#/constants/routes'
import { getUrlParams } from '#/utils'

import * as S from './styled'
import { FormValues } from './types'

export const ResetPasswordPage = () => {
  const navigate = useNavigate()
  const { formatMessage } = useIntl()
  const { onSubmit, getInputProps, isValid } = useForm<FormValues>({
    initialValues: {
      password: '',
      password_confirm: '',
    },
    validate: {
      password: (value) =>
        value.trim().length < PASSWORD_MIN_LENGTH ||
        value.trim().length !== value.length
          ? formatMessage(
              {
                description: 'Register form password error message',
                defaultMessage: `Password must include at 
                least {value} characters and 
                not start or end with space`,
                id: 'xqSIWB',
              },
              {
                value: PASSWORD_MIN_LENGTH,
              }
            )
          : null,
      password_confirm: (value, values) => {
        return value.trim() !== values.password.trim()
          ? formatMessage({
              description: 'Register form password repeat error message',
              defaultMessage: 'Field must be equal to password',
              id: 'C9oQvf',
            })
          : null
      },
    },
  })

  const handleRedirect = () => {
    const param = new URLSearchParams({
      auth: 'true',
    }).toString()
    navigate(`${ERoutes.ROOT}?${param}`)
  }

  const { mutate, isLoading, isError, isSuccess } = usePasswordResetSubmit()

  const handleSubmit = (values: FormValues) => {
    const { token } = getUrlParams('token')

    mutate({ ...values, token })
  }

  if (isError) {
    navigate(ERoutes.NOT_FOUND)
  }

  if (isSuccess) {
    return (
      <S.Div>
        <Title order={3}>
          <FormattedMessage
            description="Password reset page success message"
            defaultMessage="Success! Your password has been changed"
            id="CMYTok"
          />
        </Title>
        <Button variant="subtle" onClick={handleRedirect}>
          <FormattedMessage
            description="Password reset page redirect button text"
            defaultMessage="Go to home page"
            id="SHpxRQ"
          />
        </Button>
      </S.Div>
    )
  }
  return (
    <S.Div>
      <form onSubmit={onSubmit(handleSubmit)}>
        <PasswordInput
          placeholder={formatMessage(
            {
              description: 'Reset password form password field placeholder',
              defaultMessage: 'Must include at least {value} characters',
              id: '0uIUGn',
            },
            { value: PASSWORD_MIN_LENGTH }
          )}
          label={formatMessage({
            description: 'Reset password form password field label',
            defaultMessage: 'New password',
            id: 'HfGoBe',
          })}
          withAsterisk
          {...getInputProps('password')}
        />
        <PasswordInput
          placeholder="Must match password field"
          label={formatMessage({
            description: 'Reset password form password repeat field label',
            defaultMessage: 'Repeat password',
            id: 'eBO34P',
          })}
          withAsterisk
          {...getInputProps('password_confirm')}
        />
        <Space h="xs" />
        <Group position="center">
          <Button type="submit" loading={isLoading} disabled={!isValid()}>
            <FormattedMessage
              description="Reset password form Submit button text"
              defaultMessage="Submit"
              id="Tb4gT9"
            />
          </Button>
        </Group>
      </form>
    </S.Div>
  )
}
