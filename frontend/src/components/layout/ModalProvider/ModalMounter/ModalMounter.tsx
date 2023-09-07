import { useEffect } from 'react'

import { Props } from './types'

export const ModalMounter = ({ onMount, children }: Props) => {
  useEffect(() => {
    onMount(true)
  }, [onMount])

  return <>{children}</>
}
