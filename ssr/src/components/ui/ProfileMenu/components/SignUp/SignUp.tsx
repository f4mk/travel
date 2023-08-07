import { Menu } from '@mantine/core'
import { Milestone } from 'lucide-react'

import { Props } from './types'
export const SignUp = ({ onClick }: Props) => {
  return (
    <>
      <Menu.Item onClick={onClick} icon={<Milestone />}>
        Sign Up
      </Menu.Item>
    </>
  )
}
