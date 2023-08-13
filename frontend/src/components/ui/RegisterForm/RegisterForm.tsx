import { FormattedMessage, useIntl } from 'react-intl'
import { Button, Group, PasswordInput, Space, TextInput } from '@mantine/core'
import { useForm } from '@mantine/form'

import { FormValues, Props } from './types'

export const RegisterForm = ({ onClose }: Props) => {
  const intl = useIntl()
  const form = useForm<FormValues>({
    initialValues: {
      username: '',
      password: '',
      passwordRepeat: '',
      name: '',
      lastname: '',
      email: ''
    },

    validate: (values) => {
      return {
        username:
          values.username.trim().length < 6
            ? intl.formatMessage({
                description: 'Register form username error message',
                defaultMessage: 'Username must include at least 6 characters'
              })
            : null,
        password:
          values.password.length < 6
            ? intl.formatMessage({
                description: 'Register form password error message',
                defaultMessage: 'Password must include at least 6 characters'
              })
            : null,
        passwordRepeat:
          values.password !== values.passwordRepeat
            ? intl.formatMessage({
                description: 'Register form password repeat error message',
                defaultMessage: 'Field should be equal to password'
              })
            : null,
        name:
          values.name.trim().length < 2
            ? intl.formatMessage({
                description: 'Register form name error message',
                defaultMessage: 'Name must include at least 2 characters'
              })
            : null,
        email: /^\w+([.-]?\w+)*@\w+([.-]?\w+)*(\.\w{2,3})+$/.test(values.email)
          ? null
          : intl.formatMessage({
              description: 'Register form email error message',
              defaultMessage: 'Invalid email'
            })
      }
    }
  })

  const handleSubmit = (values: FormValues) => {
    // eslint-disable-next-line
    console.log(values)
    onClose?.()
  }

  return (
    <form onSubmit={form.onSubmit(handleSubmit)}>
      <TextInput
        placeholder={intl.formatMessage({
          description: 'Register form username field placeholder',
          defaultMessage: 'SuperJohn3000'
        })}
        label={intl.formatMessage({
          description: 'Register form username field label',
          defaultMessage: 'Username'
        })}
        withAsterisk
        {...form.getInputProps('username')}
      />

      <TextInput
        placeholder={intl.formatMessage({
          description: 'Register form name field placeholder',
          defaultMessage: 'John'
        })}
        label={intl.formatMessage({
          description: 'Register form name field label',
          defaultMessage: 'First name'
        })}
        withAsterisk
        {...form.getInputProps('name')}
      />
      <TextInput
        placeholder={intl.formatMessage({
          description: 'Register form last name field placeholder',
          defaultMessage: 'Username'
        })}
        label={intl.formatMessage({
          description: 'Register form last name field label',
          defaultMessage: 'Last name'
        })}
        {...form.getInputProps('lastname')}
      />
      <TextInput
        placeholder={intl.formatMessage({
          description: 'Register form email field placeholder',
          defaultMessage: 'user@example.com'
        })}
        label={intl.formatMessage({
          description: 'Register form email field label',
          defaultMessage: 'Email'
        })}
        withAsterisk
        {...form.getInputProps('email')}
      />
      <PasswordInput
        placeholder="******"
        label={intl.formatMessage({
          description: 'Register form password field label',
          defaultMessage: 'Password'
        })}
        withAsterisk
        {...form.getInputProps('password')}
      />
      <PasswordInput
        placeholder="******"
        label={intl.formatMessage({
          description: 'Register form password repeat field label',
          defaultMessage: 'Repeat password'
        })}
        withAsterisk
        {...form.getInputProps('passwordRepeat')}
      />
      <Space h="xs" />
      <Group position="center">
        <Button type="submit">
          <FormattedMessage
            description="Register form Sign Up button text"
            defaultMessage="Sign Up"
          />
        </Button>
        <Button variant="outline" onClick={onClose}>
          <FormattedMessage
            description="Register form Close button text"
            defaultMessage="Close"
          />
        </Button>
      </Group>
    </form>
  )
}
