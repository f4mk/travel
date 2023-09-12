import { useEffect } from 'react'

import { Props } from './types'

export const ModalMounter = ({ onMount, children }: Props) => {
  useEffect(() => {
    onMount(true)
    return () => onMount(false)
  }, [onMount])

  return <>{children}</>
}
