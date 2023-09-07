import { FormattedMessage } from 'react-intl'
import { Menu } from '@mantine/core'
import { useDisclosure } from '@mantine/hooks'
import { Cog, LogOut, MessageSquare, Trash } from 'lucide-react'

import { useModal } from '#/components/layout/ModalProvider'
import { ProfileButton } from '#/components/ui/ProfileButton'
import { lazy } from '#/utils/lazy'

import { SignIn } from './components/SignIn'
import { SignUp } from './components/SignUp'
import { EFormView, Props } from './types'

const { Auth } = lazy(() => import('#/components/ui/Auth'))

export const ProfileMenu = ({ isLoggedIn }: Props) => {
  const { showModal, hideModal } = useModal()

  const [menuOpened, { open: menuOpen, close: menuClose }] =
    useDisclosure(false)

  const handleOpen = (view: EFormView) => {
    showModal(<Auth activeTab={view} onClose={hideModal} />)
  }

  return (
    <Menu opened={menuOpened} onClose={menuClose} shadow="md" width={200}>
      <Menu.Target>
        <ProfileButton onPress={menuOpen} />
      </Menu.Target>

      <Menu.Dropdown>
        {isLoggedIn ? (
          <>
            <Menu.Item icon={<Cog />}>
              <FormattedMessage
                description="Profile menu Settings button text"
                defaultMessage="Settings"
                id="RpRPPn"
              />
            </Menu.Item>
            <Menu.Item icon={<MessageSquare />}>
              <FormattedMessage
                description="Profile menu Messages button text"
                defaultMessage="Messages"
                id="OrXAfT"
              />
            </Menu.Item>
            <Menu.Divider />
            <Menu.Item color="red" icon={<Trash />}>
              <FormattedMessage
                description="Profile menu Delete my account button text"
                defaultMessage="Delete my account"
                id="XaZp23"
              />
            </Menu.Item>
            <Menu.Divider />
            <Menu.Item icon={<LogOut />}>
              <FormattedMessage
                description="Profile menu Sign Out button text"
                defaultMessage="Sign Out"
                id="TOEbd7"
              />
            </Menu.Item>
          </>
        ) : (
          <>
            <SignIn onClick={() => handleOpen(EFormView.SIGN_IN)} />
            <SignUp onClick={() => handleOpen(EFormView.SIGN_UP)} />
          </>
        )}
      </Menu.Dropdown>
    </Menu>
  )
}
