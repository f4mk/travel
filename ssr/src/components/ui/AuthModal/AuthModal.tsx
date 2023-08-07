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

export const AuthModal = ({ opened, onClose, onSwitch }: Props) => {
  const form = useForm<FormValues>({
    initialValues: {
      email: '',
      password: ''
    }
  })

  return (
    <Modal opened={opened} onClose={onClose} title="Authentication" centered>
      <form onSubmit={form.onSubmit((values) => console.log(values))}>
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
        <Space h="xs" />
        <Group position="center">
          <Button type="submit">Sign In</Button>
        </Group>
      </form>
      <Space h="xs" />
      <Anchor component="button" type="button" onClick={onSwitch}>
        Don&apos;t have an account?
      </Anchor>
    </Modal>
  )
}
