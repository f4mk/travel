import { Button, Group, PasswordInput, Space, TextInput } from '@mantine/core'
import { useForm } from '@mantine/form'

import { FormValues, Props } from './types'

export const RegisterForm = ({ onClose }: Props) => {
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
            ? 'Username must include at least 6 characters'
            : null,
        password:
          values.password.length < 6
            ? 'Password must include at least 6 characters'
            : null,
        passwordRepeat:
          values.password !== values.passwordRepeat
            ? 'Field should be equal to password'
            : null,
        name:
          values.name.trim().length < 2
            ? 'Name must include at least 2 characters'
            : null,
        email: /^\w+([.-]?\w+)*@\w+([.-]?\w+)*(\.\w{2,3})+$/.test(values.email)
          ? null
          : 'Invalid email'
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
        placeholder="SuperJohn3000"
        label="Username"
        withAsterisk
        {...form.getInputProps('username')}
      />

      <TextInput
        placeholder="John"
        label="First name"
        withAsterisk
        {...form.getInputProps('name')}
      />
      <TextInput
        placeholder="Doe"
        label="Last name"
        {...form.getInputProps('lastname')}
      />
      <TextInput
        placeholder="user@example.com"
        label="Email"
        withAsterisk
        {...form.getInputProps('email')}
      />
      <PasswordInput
        placeholder="******"
        label="Password"
        withAsterisk
        {...form.getInputProps('password')}
      />
      <PasswordInput
        placeholder="******"
        label="Repeat password"
        withAsterisk
        {...form.getInputProps('passwordRepeat')}
      />
      <Space h="xs" />
      <Group position="center">
        <Button type="submit">Sign Up</Button>
        <Button variant="outline" onClick={onClose}>
          Close
        </Button>
      </Group>
    </form>
  )
}
