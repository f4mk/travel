import { FormattedMessage } from 'react-intl'
import { Button, Group, PasswordInput, Space, TextInput } from '@mantine/core'
import { useForm } from '@mantine/form'

import { useLogin } from '#/api/auth'
import { useMessage } from '#/hooks'

import { FormValues, Props } from './types'

export const AuthForm = ({ onClose }: Props) => {
  const message = useMessage()
  const form = useForm<FormValues>({
    initialValues: {
      email: '',
      password: '',
    },
  })

  const { mutate } = useLogin()
  const handleSubmit = (values: FormValues) => {
    mutate({ email: values.email, password: values.password })

    onClose()
  }

  return (
    <form onSubmit={form.onSubmit(handleSubmit)}>
      <TextInput
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
        {...form.getInputProps('email')}
      />
      <PasswordInput
        placeholder="********"
        label={message({
          description: 'Auth form password label',
          defaultMessage: 'Password',
          id: '06sNqJ',
        })}
        withAsterisk
        {...form.getInputProps('password')}
      />
      <Space h="xs" />
      <Group position="center">
        <Button type="submit">
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
