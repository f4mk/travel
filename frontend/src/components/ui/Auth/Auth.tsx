import { FormattedMessage } from 'react-intl'
import { Modal, Tabs } from '@mantine/core'
import FocusTrap from 'focus-trap-react'
import { LogIn, Milestone } from 'lucide-react'

import { EFormView } from '#/components/ui/ProfileMenu'

import { AuthForm } from '../AuthForm'
import { RegisterForm } from '../RegisterForm'

import { Props } from './types'

export const Auth = ({ opened, activeTab, setActiveTab, onClose }: Props) => {
  return (
    <Modal
      opened={opened}
      onClose={onClose}
      withCloseButton={false}
      closeOnClickOutside
    >
      {/* NOTE: Focus trap by mantine doesnt work here */}
      <FocusTrap>
        <Tabs variant="outline" value={activeTab} onTabChange={setActiveTab}>
          <Tabs.List>
            <Tabs.Tab value={EFormView.SIGN_IN} icon={<LogIn />}>
              <FormattedMessage
                description="Authentication window title"
                defaultMessage="Authentication"
                id="0JzVNd"
              />
            </Tabs.Tab>
            <Tabs.Tab value={EFormView.SIGN_UP} icon={<Milestone />}>
              <FormattedMessage
                description="Registration window title"
                defaultMessage="Registration"
                id="70PPhl"
              />
            </Tabs.Tab>
          </Tabs.List>

          <Tabs.Panel value={EFormView.SIGN_IN} pt="xs">
            <AuthForm onClose={onClose} />
          </Tabs.Panel>
          <Tabs.Panel value={EFormView.SIGN_UP} pt="xs">
            <RegisterForm onClose={onClose} />
          </Tabs.Panel>
        </Tabs>
      </FocusTrap>
    </Modal>
  )
}
