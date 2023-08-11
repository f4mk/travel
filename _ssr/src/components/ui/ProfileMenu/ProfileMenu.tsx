import { useCallback, useState } from 'react'
import { Menu } from '@mantine/core'
import { useDisclosure } from '@mantine/hooks'
import { Cog, LogOut, MessageSquare, Trash } from 'lucide-react'

import { ProfileButton } from '#/components/ui/ProfileButton'

import { AuthModal } from '../AuthModal'
import { RegisterModal } from '../RegisterModal'

import { SignIn } from './components/SignIn'
import { SignUp } from './components/SignUp'
import { FormView, Props } from './types'
export const ProfileMenu = ({ isLoggedIn }: Props) => {
  const [menuOpened, { open: menuOpen, close: menuClose }] =
    useDisclosure(false)
  const [modalOpened, { open: modalOpen, close: modalClose }] =
    useDisclosure(false)
  const [view, setView] = useState<FormView>(FormView.SIGN_IN)

  const handleOpen = (view: FormView) => {
    setView(view)
    modalOpen()
  }
  const handleSwitch = useCallback(() => {
    setView((oldView) =>
      oldView === FormView.SIGN_IN ? FormView.SIGN_UP : FormView.SIGN_IN
    )
  }, [])

  return (
    <>
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
              <SignIn onClick={() => handleOpen(FormView.SIGN_IN)} />
              <SignUp onClick={() => handleOpen(FormView.SIGN_UP)} />
            </>
          )}
        </Menu.Dropdown>
      </Menu>

      {view === FormView.SIGN_IN ? (
        <AuthModal
          opened={modalOpened}
          onClose={modalClose}
          onSwitch={handleSwitch}
        />
      ) : (
        <RegisterModal
          opened={modalOpened}
          onClose={modalClose}
          onSwitch={handleSwitch}
        />
      )}
    </>
  )
}
