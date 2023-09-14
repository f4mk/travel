import { useState } from 'react'
import { FormattedMessage } from 'react-intl'
import { Tabs } from '@mantine/core'
import { LogIn, Milestone } from 'lucide-react'

import { AuthForm } from '../AuthForm'
import { RegisterForm } from '../RegisterForm'

import { EFormView, Props } from './types'

export const Auth = ({ activeTab, onClose }: Props) => {
  const [view, setView] = useState<EFormView>(activeTab)

  return (
    <Tabs
      variant="outline"
      value={view}
      onTabChange={(tab) => setView(tab as EFormView)}
    >
      <Tabs.List>
        <Tabs.Tab value={EFormView.SIGN_UP} icon={<Milestone />}>
          <FormattedMessage
            description="Registration window title"
            defaultMessage="Registration"
            id="70PPhl"
          />
        </Tabs.Tab>
        <Tabs.Tab value={EFormView.SIGN_IN} icon={<LogIn />}>
          <FormattedMessage
            description="Authentication window title"
            defaultMessage="Authentication"
            id="0JzVNd"
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
  )
}
