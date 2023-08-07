import { PropsWithChildren } from 'react'

export type Props = PropsWithChildren<{
  position: google.maps.LatLng | google.maps.LatLngLiteral
  pane?: keyof google.maps.MapPanes
  map: google.maps.Map
  zIndex?: number
  preventInteraction?: boolean
}>
