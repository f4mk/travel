import { Suspense, useEffect, useRef, useState } from 'react'
import FocusTrap from 'focus-trap-react'

import { CenteredLoader } from '#/components/ui/CenteredLoader'

import { ModalMounter } from '../ModalMounter'

import * as S from './styled'
import { Props } from './types'

export const Modal = ({ opened, onClose, children, size = 'md' }: Props) => {
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
    <S.Modal
      opened={opened}
      onClose={onClose}
      withCloseButton={false}
      closeOnClickOutside
      size={size}
    >
      {/* NOTE: Focus trap by mantine doesnt work here */}
      <FocusTrap active={mounted}>
        <div ref={ref}>
          <Suspense fallback={<CenteredLoader />}>
            <ModalMounter onMount={setMounted}>{children}</ModalMounter>
          </Suspense>
        </div>
      </FocusTrap>
    </S.Modal>
  )
}
