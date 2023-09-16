import { FormattedMessage, useIntl } from 'react-intl'
import { useNavigate } from 'react-router-dom'
import { Button, Group, PasswordInput, Space, TextInput } from '@mantine/core'
import { useForm } from '@mantine/form'

import { useCreateUser } from '#/api/user'
import { ERoutes } from '#/constants/routes'

import {
  PASSWORD_MIN_LENGTH,
  USERNAME_MAX_LENGTH,
  USERNAME_MIN_LENGTH,
} from './consts'
import { FormValues, Props } from './types'

export const RegisterForm = ({ onClose }: Props) => {
  const navigate = useNavigate()
  const { formatMessage } = useIntl()
  const { getInputProps, onSubmit, isValid, setErrors, values } =
    useForm<FormValues>({
      initialValues: {
        name: '',
        password: '',
        password_confirm: '',
        email: '',
      },

      validate: {
        name: (value) =>
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
          value.trim().length < PASSWORD_MIN_LENGTH ||
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
                  value: PASSWORD_MIN_LENGTH,
                }
              )
            : null,
        password_confirm: (value, values) =>
          value.trim() !== values.password.trim()
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
    navigate(ERoutes.USER_CREATE)
    onClose()
  }

  const { mutate, isLoading } = useCreateUser({
    onSuccess: handleSuccess,
    onError: (error) => {
      if (error.response.status === 400) {
        const errors: { [K in keyof FormValues]?: true } = {}
        Object.keys(values).forEach(
          (field) => (errors[field as keyof FormValues] = true)
        )
        setErrors(errors)
      }
      if (error.response.status === 409) {
        setErrors({ email: true, name: true })
      }
    },
  })

  const handleSubmit = (values: FormValues) => {
    mutate(values)
  }

  return (
    <form onSubmit={onSubmit(handleSubmit)}>
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
        {...getInputProps('name')}
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
        {...getInputProps('email')}
      />
      <PasswordInput
        placeholder={formatMessage(
          {
            description: 'Register form password field placeholder',
            defaultMessage: 'Must include at least {value} characters',
            id: 'M4y+WF',
          },
          { value: PASSWORD_MIN_LENGTH }
        )}
        label={formatMessage({
          description: 'Register form password field label',
          defaultMessage: 'Password',
          id: 'TQEu8X',
        })}
        withAsterisk
        {...getInputProps('password')}
      />
      <PasswordInput
        placeholder="Must match Password field"
        label={formatMessage({
          description: 'Register form password repeat field label',
          defaultMessage: 'Repeat password',
          id: '+BKGrr',
        })}
        withAsterisk
        {...getInputProps('password_confirm')}
      />
      <Space h="xs" />
      <Group position="center">
        <Button type="submit" loading={isLoading} disabled={!isValid()}>
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
