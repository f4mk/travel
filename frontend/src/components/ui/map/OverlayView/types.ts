import { PropsWithChildren } from 'react'

export type Props = PropsWithChildren<{
  position: google.maps.LatLng | google.maps.LatLngLiteral
  preventInteraction?: boolean
}>
