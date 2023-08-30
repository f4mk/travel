import { useEffect, useRef } from 'react'
import {
  OverlayView as DefaultOverlayView,
  OverlayViewF
} from '@react-google-maps/api'

import type { Props } from './types'

export const OverlayView = ({
  position,
  children,
  preventInteraction
}: Props) => {
  const containerRef = useRef(null)
  useEffect(() => {
    if (containerRef.current && preventInteraction) {
      google.maps.OverlayView.preventMapHitsAndGesturesFrom(
        containerRef.current
      )
    }
  }, [containerRef, preventInteraction])

  return (
    <OverlayViewF
      position={position}
      mapPaneName={DefaultOverlayView.OVERLAY_MOUSE_TARGET}
    >
      <div ref={containerRef}>{children}</div>
    </OverlayViewF>
  )
}
