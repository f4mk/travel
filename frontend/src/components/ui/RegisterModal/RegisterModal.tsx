import {
  Anchor,
  Button,
  Group,
  Modal,
  PasswordInput,
  Space,
  TextInput
} from '@mantine/core'
import { useForm } from '@mantine/form'

import { FormValues, Props } from './types'

export const RegisterModal = ({ opened, onClose, onSwitch }: Props) => {
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

  return (
    <Modal opened={opened} onClose={onClose} title="Registration" centered>
      <form>
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
        </Group>
      </form>
      <Space h="xs" />
      <Anchor component="button" type="button" onClick={onSwitch}>
        Already have an account?
      </Anchor>
    </Modal>
  )
}
