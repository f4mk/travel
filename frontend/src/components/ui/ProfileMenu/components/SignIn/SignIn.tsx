import { FormattedMessage } from 'react-intl'
import { Menu } from '@mantine/core'
import { LogIn } from 'lucide-react'

import { Props } from './types'
export const SignIn = ({ onClick }: Props) => {
  return (
    <Menu.Item onClick={onClick} icon={<LogIn />}>
      <FormattedMessage
        description="Profile menu Sign in button text"
        defaultMessage="Sign In"
        id="heIN4y"
      />
    </Menu.Item>
  )
}
