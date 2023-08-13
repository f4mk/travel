import { FormattedMessage } from 'react-intl'
import { Button, Group, PasswordInput, Space, TextInput } from '@mantine/core'
import { useForm } from '@mantine/form'

import { useMessage } from '#/hooks'

import { FormValues, Props } from './types'

export const RegisterForm = ({ onClose }: Props) => {
  const message = useMessage()
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
            ? message({
                description: 'Register form username error message',
                defaultMessage: 'Username must include at least 6 characters',
                id: 'v3oysd'
              })
            : null,
        password:
          values.password.length < 6
            ? message({
                description: 'Register form password error message',
                defaultMessage: 'Password must include at least 6 characters',
                id: 'JtqjO1'
              })
            : null,
        passwordRepeat:
          values.password !== values.passwordRepeat
            ? message({
                description: 'Register form password repeat error message',
                defaultMessage: 'Field should be equal to password',
                id: 'MPtkD8'
              })
            : null,
        name:
          values.name.trim().length < 2
            ? message({
                description: 'Register form name error message',
                defaultMessage: 'Name must include at least 2 characters',
                id: 'sHDcWI'
              })
            : null,
        email: /^\w+([.-]?\w+)*@\w+([.-]?\w+)*(\.\w{2,3})+$/.test(values.email)
          ? null
          : message({
              description: 'Register form email error message',
              defaultMessage: 'Invalid email',
              id: '75g0FC'
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
        placeholder={message({
          description: 'Register form username field placeholder',
          defaultMessage: 'SuperJohn3000',
          id: 'SnPjw3'
        })}
        label={message({
          description: 'Register form username field label',
          defaultMessage: 'Username',
          id: 'oWGSmJ'
        })}
        withAsterisk
        {...form.getInputProps('username')}
      />

      <TextInput
        placeholder={message({
          description: 'Register form name field placeholder',
          defaultMessage: 'John',
          id: 'zw4spD'
        })}
        label={message({
          description: 'Register form name field label',
          defaultMessage: 'First name',
          id: 'iZRk5H'
        })}
        withAsterisk
        {...form.getInputProps('name')}
      />
      <TextInput
        placeholder={message({
          description: 'Register form last name field placeholder',
          defaultMessage: 'Username',
          id: '5V6R/N'
        })}
        label={message({
          description: 'Register form last name field label',
          defaultMessage: 'Last name',
          id: 'mp6K2g'
        })}
        {...form.getInputProps('lastname')}
      />
      <TextInput
        placeholder={message({
          description: 'Register form email field placeholder',
          defaultMessage: 'user@example.com',
          id: 'R14xEd'
        })}
        label={message({
          description: 'Register form email field label',
          defaultMessage: 'Email',
          id: 'AyMd2C'
        })}
        withAsterisk
        {...form.getInputProps('email')}
      />
      <PasswordInput
        placeholder="******"
        label={message({
          description: 'Register form password field label',
          defaultMessage: 'Password',
          id: 'TQEu8X'
        })}
        withAsterisk
        {...form.getInputProps('password')}
      />
      <PasswordInput
        placeholder="******"
        label={message({
          description: 'Register form password repeat field label',
          defaultMessage: 'Repeat password',
          id: '+BKGrr'
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
            id="qEt3X4"
          />
        </Button>
        <Button variant="outline" onClick={onClose}>
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
