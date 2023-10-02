import { FormattedMessage } from 'react-intl'
import { useNavigate } from 'react-router-dom'
import { Menu } from '@mantine/core'
import { useDisclosure } from '@mantine/hooks'
import { Cog, LogOut, MessageSquare, Trash } from 'lucide-react'

import { useLogout } from '#/api/auth'
import { ProfileButton } from '#/components/ui/ProfileButton'
import { ROUTES } from '#/constants/routes'

export const ProfileMenu = () => {
  const navigate = useNavigate()
  const [menuOpened, { open: menuOpen, close: menuClose }] =
    useDisclosure(false)

  const { refetch: logout } = useLogout({
    enabled: false,
    onSuccess: () => {
      navigate(ROUTES.ROOT)
    },
  })

  return (
    <Menu opened={menuOpened} onClose={menuClose} shadow="md" width={200}>
      <Menu.Target>
        <ProfileButton onPress={menuOpen} />
      </Menu.Target>

      <Menu.Dropdown>
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
        <Menu.Item icon={<LogOut />} onClick={() => logout()}>
          <FormattedMessage
            description="Profile menu Sign Out button text"
            defaultMessage="Sign Out"
            id="TOEbd7"
          />
        </Menu.Item>
      </Menu.Dropdown>
    </Menu>
  )
}
