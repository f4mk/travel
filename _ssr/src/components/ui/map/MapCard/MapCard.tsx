import { Card } from '#/components/ui/Card'
import { Marker } from '#/components/ui/map/Marker'

import { Props } from './types'

export const MapCard = ({ lat, lng, map }: Props) => {
  return (
    <Marker lat={lat} lng={lng} map={map} preventInteraction>
      {/* TODO: use palette */}
      <Card />
    </Marker>
  )
}
