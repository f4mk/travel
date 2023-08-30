import { OverlayView } from '#/components/ui/map/OverlayView'

import { MarkerContainer } from './MarkerContainer'
import { Props } from './types'

export const Marker = ({ lat, lng, preventInteraction, children }: Props) => {
  return (
    <OverlayView
      position={{
        lat: lat,
        lng: lng
      }}
      preventInteraction={preventInteraction}
    >
      <MarkerContainer>{children}</MarkerContainer>
    </OverlayView>
  )
}
