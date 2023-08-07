import { useEffect, useMemo } from 'react'
import { createPortal } from 'react-dom'

import { createOverlay } from './Overlay'
import { Props } from './types'
export const OverlayView = ({
  position,
  pane = 'floatPane',
  map,
  zIndex,
  preventInteraction,
  children
}: Props) => {
  const container = useMemo(() => {
    const div = document.createElement('div')
    div.id = '___map-overlay'
    div.style.position = 'absolute'
    return div
  }, [])

  const overlay = useMemo(() => {
    return createOverlay(container, pane, position)
  }, [container, pane, position])

  useEffect(() => {
    overlay?.setMap(map)
    preventInteraction &&
      google.maps.OverlayView.preventMapHitsAndGesturesFrom(container)
    return () => overlay?.setMap(null)
  }, [map, overlay, container, preventInteraction])

  useEffect(() => {
    container.style.zIndex = `${zIndex}`
  }, [zIndex, container])

  // useEffect(() => {
  //   function checkOverlayClick(event: MouseEvent | TouchEvent) {
  //     const target = event.target as Element
  //     if (target.closest('#___map-overlay')) {
  //       event.stopPropagation()
  //       if ('touches' in event) {
  //         event.preventDefault()
  //       }
  //     }
  //   }

  //   document.addEventListener('mousedown', checkOverlayClick, {
  //     passive: false
  //   })
  //   return () => document.removeEventListener('mousedown', checkOverlayClick)
  // }, [])

  return createPortal(children, container)
}
