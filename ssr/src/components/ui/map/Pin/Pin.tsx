import { usePress } from 'react-aria'
import { MapPin } from 'lucide-react'

import { Marker } from '#/components/ui/map/Marker'

import * as S from './styled'
import { Props } from './types'
export const Pin = ({ lat, lng, map, onPress }: Props) => {
  const { pressProps } = usePress({ onPress })
  return (
    <Marker lat={lat} lng={lng} map={map}>
      <S.Div {...(!!onPress && { pressProps })} isPressable={!!onPress}>
        {/* TODO: use palette */}
        <MapPin color="red" size={32} fill="beige" fillOpacity={20} />
      </S.Div>
    </Marker>
  )
}
