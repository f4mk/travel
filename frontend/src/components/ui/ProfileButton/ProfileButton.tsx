import { forwardRef } from 'react'
import { usePress } from 'react-aria'
import { Avatar } from '@mantine/core'

import avatar from '#/assets/coggers.png'
import { RoundButton } from '#/components/ui/RoundButton'
import { useMessage } from '#/hooks'

import { Props } from './types'

export const ProfileButton = forwardRef<HTMLButtonElement, Props>(
  ({ onPress }, ref) => {
    const message = useMessage()
    const { pressProps } = usePress({ onPress })
    return (
      <RoundButton ref={ref} {...pressProps}>
        <Avatar
          src={avatar}
          alt={message({
            description: 'Alt description on user avatar',
            defaultMessage: 'profile',
            id: 'jI32/e',
          })}
        />
      </RoundButton>
    )
  }
)
