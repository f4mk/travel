import { FormattedMessage, useIntl } from 'react-intl'
import { Button, Group, PasswordInput, Space, TextInput } from '@mantine/core'
import { useForm } from '@mantine/form'

import { FormValues, Props } from './types'

export const AuthForm = ({ onClose }: Props) => {
  const intl = useIntl()

  const form = useForm<FormValues>({
    initialValues: {
      email: '',
      password: ''
    }
  })

  const handleSubmit = (values: FormValues) => {
    // eslint-disable-next-line
    console.log(values)
    onClose()
  }

  return (
    <form onSubmit={form.onSubmit(handleSubmit)}>
      <TextInput
        placeholder={intl.formatMessage({
          description: 'Auth form email placeholder',
          defaultMessage: 'user@example.com'
        })}
        label={intl.formatMessage({
          description: 'Auth form email label',
          defaultMessage: 'Email'
        })}
        withAsterisk
        {...form.getInputProps('email')}
      />
      <PasswordInput
        placeholder="******"
        label={intl.formatMessage({
          description: 'Auth form password lable',
          defaultMessage: 'Password'
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
        <Button variant="outline" onClick={onClose}>
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
