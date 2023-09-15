import { forwardRef } from 'react'
import { usePress } from 'react-aria'
import { useIntl } from 'react-intl'
import { Avatar } from '@mantine/core'

import avatar from '#/assets/coggers.png'
import { RoundButton } from '#/components/ui/RoundButton'

import { Props } from './types'

export const ProfileButton = forwardRef<HTMLButtonElement, Props>(
  ({ onPress }, ref) => {
    const { formatMessage } = useIntl()
    const { pressProps } = usePress({ onPress })
    return (
      <RoundButton ref={ref} {...pressProps}>
        <Avatar
          src={avatar}
          alt={formatMessage({
            description: 'Alt description on user avatar',
            defaultMessage: 'profile',
            id: 'jI32/e',
          })}
        />
      </RoundButton>
    )
  }
)
