import { usePress } from 'react-aria'
import { Pin as PinIcon } from 'lucide-react'

import { Marker } from '#/components/ui/map/Marker'

import * as S from './styled'
import { Props } from './types'
export const Pin = ({ lat, lng, onPress }: Props) => {
  const { pressProps } = usePress({ onPress })
  return (
    <Marker lat={lat} lng={lng}>
      <S.Div {...pressProps} isPressable={!!onPress}>
        {/* TODO: use palette */}
        <PinIcon color="red" size={32} fill="red" fillOpacity={20} />
      </S.Div>
    </Marker>
  )
}
