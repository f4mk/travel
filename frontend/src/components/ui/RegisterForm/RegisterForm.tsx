import { FormattedMessage } from 'react-intl'
import { Button, Group, PasswordInput, Space, TextInput } from '@mantine/core'
import { useForm } from '@mantine/form'

import { useCreateUser } from '#/api/user'
import { CreateUserRequest } from '#/api/user/types'
import { useMessage } from '#/hooks'

import { FormValues, Props } from './types'

export const RegisterForm = ({ onClose }: Props) => {
  const message = useMessage()
  const form = useForm<FormValues>({
    initialValues: {
      username: '',
      password: '',
      passwordRepeat: '',
      email: '',
    },

    validate: {
      username: (value) =>
        value.trim().length < 2
          ? message({
              description: 'Register form username error message',
              defaultMessage:
                'Username must include at least 2 characters and not start or end with space',
              id: 'v3oysd',
            })
          : null,
      password: (value) =>
        value.trim().length < 8 || value.trim().length !== value.length
          ? message({
              description: 'Register form password error message',
              defaultMessage:
                'Password must include at least 8 characters and not start or end with space',
              id: 'JtqjO1',
            })
          : null,
      passwordRepeat: (value, values) =>
        value.trim() !== values.passwordRepeat.trim()
          ? message({
              description: 'Register form password repeat error message',
              defaultMessage: 'Field must be equal to password',
              id: 'MPtkD8',
            })
          : null,
      email: (value) =>
        /^\w+([.-]?\w+)*@\w+([.-]?\w+)*(\.\w{2,3})+$/.test(value.trim())
          ? null
          : message({
              description: 'Register form email error message',
              defaultMessage: 'Invalid email',
              id: '75g0FC',
            }),
    },
  })

  const handleSuccess = () => {
    // TODO: add redirect to confirm page
    // navigate(ERoutes.MAP)
    onClose()
  }

  const { mutate, isLoading } = useCreateUser({
    onSuccess: handleSuccess,
    // TODO: handle errors
    // onError: () => setHasErrors(true),
  })

  const handleSubmit = (values: FormValues) => {
    const data: CreateUserRequest = {
      name: values.username,
      password_confirm: values.passwordRepeat,
      password: values.password,
      email: values.email,
    }
    mutate(data)
  }

  return (
    <form onSubmit={form.onSubmit(handleSubmit)}>
      <TextInput
        placeholder={message({
          description: 'Register form username field placeholder',
          defaultMessage: 'SuperJohn3000',
          id: 'SnPjw3',
        })}
        label={message({
          description: 'Register form username field label',
          defaultMessage: 'Username',
          id: 'oWGSmJ',
        })}
        withAsterisk
        {...form.getInputProps('username')}
      />
      <TextInput
        placeholder={message({
          description: 'Register form email field placeholder',
          defaultMessage: 'user@example.com',
          id: 'R14xEd',
        })}
        label={message({
          description: 'Register form email field label',
          defaultMessage: 'Email',
          id: 'AyMd2C',
        })}
        withAsterisk
        {...form.getInputProps('email')}
      />
      <PasswordInput
        placeholder="********"
        label={message({
          description: 'Register form password field label',
          defaultMessage: 'Password',
          id: 'TQEu8X',
        })}
        withAsterisk
        {...form.getInputProps('password')}
      />
      <PasswordInput
        placeholder="********"
        label={message({
          description: 'Register form password repeat field label',
          defaultMessage: 'Repeat password',
          id: '+BKGrr',
        })}
        withAsterisk
        {...form.getInputProps('passwordRepeat')}
      />
      <Space h="xs" />
      <Group position="center">
        <Button type="submit" loading={isLoading}>
          <FormattedMessage
            description="Register form Sign Up button text"
            defaultMessage="Sign Up"
            id="qEt3X4"
          />
        </Button>
        <Button variant="outline" onClick={() => onClose()}>
          <FormattedMessage
            description="Register form Close button text"
            defaultMessage="Close"
            id="tOTXPP"
          />
        </Button>
      </Group>
    </form>
  )
}
