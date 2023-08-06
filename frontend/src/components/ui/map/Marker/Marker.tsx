import { OverlayView } from '#/components/ui/map/OverlayView'

// import { OverlayView } from '@react-google-maps/api'
import { MarkerContainer } from './components/MarkerContainer'
import { Props } from './types'

export const Marker = ({
  lat,
  lng,
  map,
  preventInteraction,
  children
}: Props) => {
  return (
    <OverlayView
      // mapPaneName={OverlayView.OVERLAY_MOUSE_TARGET}
      preventInteraction={preventInteraction}
      position={{
        lat: lat,
        lng: lng
      }}
      map={map}
    >
      <MarkerContainer>{children}</MarkerContainer>
    </OverlayView>
  )
}
