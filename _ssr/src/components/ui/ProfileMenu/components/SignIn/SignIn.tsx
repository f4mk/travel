import { Menu } from '@mantine/core'
import { LogIn } from 'lucide-react'

import { Props } from './types'
export const SignIn = ({ onClick }: Props) => {
  return (
    <>
      <Menu.Item onClick={onClick} icon={<LogIn />}>
        Sign In
      </Menu.Item>
    </>
  )
}
