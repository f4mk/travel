import { Button, Group, PasswordInput, Space, TextInput } from '@mantine/core'
import { useForm } from '@mantine/form'

import { FormValues, Props } from './types'

export const AuthForm = ({ onClose }: Props) => {
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
        <Button variant="outline" onClick={onClose}>
          Close
        </Button>
      </Group>
    </form>
  )
}
