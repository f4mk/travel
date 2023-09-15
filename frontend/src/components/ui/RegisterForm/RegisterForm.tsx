import { FormattedMessage, useIntl } from 'react-intl'
import { Button, Group, PasswordInput, Space, TextInput } from '@mantine/core'
import { useForm } from '@mantine/form'

import { useCreateUser } from '#/api/user'
import { CreateUserRequest } from '#/api/user/types'

import {
  PASSWROD_MIN_LENGTH,
  USERNAME_MAX_LENGTH,
  USERNAME_MIN_LENGTH,
} from './consts'
import { FormValues, Props } from './types'

export const RegisterForm = ({ onClose }: Props) => {
  const { formatMessage } = useIntl()
  const form = useForm<FormValues>({
    initialValues: {
      username: '',
      password: '',
      passwordRepeat: '',
      email: '',
    },

    validate: {
      username: (value) =>
        value.trim().length < USERNAME_MIN_LENGTH ||
        value.trim().length > USERNAME_MAX_LENGTH
          ? formatMessage(
              {
                description: 'Register form username error message',
                defaultMessage: `Username must be {minValue} - {maxValue} characters and 
                not start or end with space`,
                id: 'xsMslr',
              },
              {
                minValue: USERNAME_MIN_LENGTH,
                maxValue: USERNAME_MAX_LENGTH,
              }
            )
          : null,
      password: (value) =>
        value.trim().length < PASSWROD_MIN_LENGTH ||
        value.trim().length !== value.length
          ? formatMessage(
              {
                description: 'Register form password error message',
                defaultMessage: `Password must include at 
                least {value} characters and 
                not start or end with space`,
                id: 'xqSIWB',
              },
              {
                value: PASSWROD_MIN_LENGTH,
              }
            )
          : null,
      passwordRepeat: (value, values) =>
        value.trim() !== values.passwordRepeat.trim()
          ? formatMessage({
              description: 'Register form password repeat error message',
              defaultMessage: 'Field must be equal to password',
              id: 'C9oQvf',
            })
          : null,
      email: (value) =>
        /^\w+([.-]?\w+)*@\w+([.-]?\w+)*(\.\w{2,3})+$/.test(value.trim())
          ? null
          : formatMessage({
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
        placeholder={formatMessage(
          {
            description: 'Register form username field placeholder',
            defaultMessage: 'Must be {minValue} - {maxValue} characters long',
            id: 'oq79wm',
          },
          {
            minValue: USERNAME_MIN_LENGTH,
            maxValue: USERNAME_MAX_LENGTH,
          }
        )}
        label={formatMessage({
          description: 'Register form username field label',
          defaultMessage: 'Username',
          id: 'oWGSmJ',
        })}
        withAsterisk
        {...form.getInputProps('username')}
      />
      <TextInput
        placeholder={formatMessage({
          description: 'Register form email field placeholder',
          defaultMessage: 'Must be valid email address',
          id: 'FX4z2/',
        })}
        label={formatMessage({
          description: 'Register form email field label',
          defaultMessage: 'Email',
          id: 'AyMd2C',
        })}
        withAsterisk
        {...form.getInputProps('email')}
      />
      <PasswordInput
        placeholder={formatMessage(
          {
            description: 'Register form password field placeholder',
            defaultMessage: 'Must include at least {value} characters',
            id: 'M4y+WF',
          },
          { value: PASSWROD_MIN_LENGTH }
        )}
        label={formatMessage({
          description: 'Register form password field label',
          defaultMessage: 'Password',
          id: 'TQEu8X',
        })}
        withAsterisk
        {...form.getInputProps('password')}
      />
      <PasswordInput
        placeholder="Must match Password field"
        label={formatMessage({
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
