import { ReactNode } from 'react'

import type { MapMarker } from '#/components/ui/map/types'

export type Props = MapMarker & {
  preventInteraction?: boolean
  children: ReactNode
}
