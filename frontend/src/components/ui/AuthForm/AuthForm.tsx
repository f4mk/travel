import { FormattedMessage } from 'react-intl'
import { useNavigate } from 'react-router-dom'
import { Button, Group, PasswordInput, Space, TextInput } from '@mantine/core'
import { hasLength, isEmail, useForm } from '@mantine/form'

import { useLogin } from '#/api/auth'
import { ERoutes } from '#/constants/routes'
import { useMessage } from '#/hooks'

import { FormValues, Props } from './types'

export const AuthForm = ({ onClose }: Props) => {
  const navigate = useNavigate()
  const message = useMessage()
  const { onSubmit, getInputProps, isValid, setErrors } = useForm<FormValues>({
    initialValues: {
      email: '',
      password: '',
    },
    validate: {
      email: isEmail(),
      password: hasLength({ min: 1 }),
    },
  })
  const handleSuccess = () => {
    navigate(ERoutes.MAP)
    onClose()
  }
  const { mutate, isLoading } = useLogin({
    onSuccess: handleSuccess,
    onError: () => setErrors({ email: true, password: true }),
  })

  const handleSubmit = (values: FormValues) => {
    mutate(values)
  }

  return (
    <form onSubmit={onSubmit(handleSubmit)}>
      <TextInput
        {...getInputProps('email')}
        placeholder={message({
          description: 'Auth form email placeholder',
          defaultMessage: 'user@example.com',
          id: 'AohQXw',
        })}
        label={message({
          description: 'Auth form email label',
          defaultMessage: 'Email',
          id: '7lT95G',
        })}
        withAsterisk
      />
      <PasswordInput
        {...getInputProps('password')}
        placeholder="********"
        label={message({
          description: 'Auth form password label',
          defaultMessage: 'Password',
          id: '06sNqJ',
        })}
        withAsterisk
      />
      <Space h="xs" />
      <Group position="center">
        <Button type="submit" loading={isLoading} disabled={!isValid()}>
          <FormattedMessage
            description="Auth form Sign in button text"
            defaultMessage="Sign In"
            id="kqb0Va"
          />
        </Button>
        <Button variant="outline" onClick={() => onClose()}>
          <FormattedMessage
            description="Auth form Close button text"
            defaultMessage="Close"
            id="jRdArt"
          />
        </Button>
      </Group>
    </form>
  )
}
