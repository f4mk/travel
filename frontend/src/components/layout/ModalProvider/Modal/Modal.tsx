import { Suspense, useEffect, useRef, useState } from 'react'
import { Modal as ModalUI } from '@mantine/core'
import FocusTrap from 'focus-trap-react'

import { CenteredLoader } from '#/components/ui/CenteredLoader'

import { ModalMounter } from '../ModalMounter'

import { Props } from './types'

export const Modal = ({ opened, onClose, children }: Props) => {
  const ref = useRef<HTMLDivElement>(null)
  const [mounted, setMounted] = useState(false)
  useEffect(() => {
    const handleOutsideClick = (event: MouseEvent) => {
      const clickedElement = event.target as HTMLElement
      if (ref.current && !ref.current.contains(clickedElement)) {
        onClose()
      }
    }

    document.addEventListener('pointerdown', handleOutsideClick)

    return () => {
      document.removeEventListener('pointerdown', handleOutsideClick)
    }
  }, [onClose])

  return (
    <ModalUI
      opened={opened}
      onClose={onClose}
      withCloseButton={false}
      closeOnClickOutside
    >
      {/* NOTE: Focus trap by mantine doesnt work here */}
      <FocusTrap active={mounted}>
        <div ref={ref} tabIndex={-1}>
          <Suspense fallback={<CenteredLoader />}>
            <ModalMounter onMount={() => setMounted(true)}>
              {children}
            </ModalMounter>
          </Suspense>
        </div>
      </FocusTrap>
    </ModalUI>
  )
}
