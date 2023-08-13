import { FormattedMessage } from 'react-intl'
import { Menu } from '@mantine/core'
import { Milestone } from 'lucide-react'

import { Props } from './types'
export const SignUp = ({ onClick }: Props) => {
  return (
    <Menu.Item onClick={onClick} icon={<Milestone />}>
      <FormattedMessage
        description="Profile menu Sign up button text"
        defaultMessage="Sign Up"
      />
    </Menu.Item>
  )
}
