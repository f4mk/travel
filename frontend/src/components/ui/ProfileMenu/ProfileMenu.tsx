import { useState } from 'react'
import { Menu } from '@mantine/core'
import { useDisclosure } from '@mantine/hooks'
import { Cog, LogOut, MessageSquare, Trash } from 'lucide-react'

import { Auth } from '#/components/ui/Auth'
import { ProfileButton } from '#/components/ui/ProfileButton'

import { SignIn } from './components/SignIn'
import { SignUp } from './components/SignUp'
import { EFormView, Props } from './types'

export const ProfileMenu = ({ isLoggedIn }: Props) => {
  const [menuOpened, { open: menuOpen, close: menuClose }] =
    useDisclosure(false)
  const [modalOpened, { open: modalOpen, close: modalClose }] =
    useDisclosure(false)
  const [view, setView] = useState<EFormView>(EFormView.SIGN_IN)

  const handleOpen = (view: EFormView) => {
    setView(view)
    modalOpen()
  }

  return (
    <Menu opened={menuOpened} onClose={menuClose} shadow="md" width={200}>
      <Menu.Target>
        <ProfileButton onPress={menuOpen} />
      </Menu.Target>

      <Menu.Dropdown>
        {isLoggedIn ? (
          <>
            <Menu.Item icon={<Cog />}>Settings</Menu.Item>
            <Menu.Item icon={<MessageSquare />}>Messages</Menu.Item>
            <Menu.Divider />
            <Menu.Item color="red" icon={<Trash />}>
              Delete my account
            </Menu.Item>
            <Menu.Divider />
            <Menu.Item icon={<LogOut />}>Sign Out</Menu.Item>
          </>
        ) : (
          <>
            <SignIn onClick={() => handleOpen(EFormView.SIGN_IN)} />
            <SignUp onClick={() => handleOpen(EFormView.SIGN_UP)} />
          </>
        )}
      </Menu.Dropdown>
      <Auth
        opened={modalOpened}
        activeTab={view}
        setActiveTab={setView}
        onClose={modalClose}
      />
    </Menu>
  )
}
