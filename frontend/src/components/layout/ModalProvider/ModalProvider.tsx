import {
  ComponentType,
  ReactElement,
  useCallback,
  useMemo,
  useState
} from 'react'
import { useDisclosure } from '@mantine/hooks'

import { Modal } from './Modal'
import { ModalContext } from './ModalContext'
import { Props } from './types'

export const ModalProvider = ({ children }: Props) => {
  const [opened, { open, close }] = useDisclosure(false)
  const [ModalContent, setModalContent] =
    useState<Nullable<ComponentType>>(null)

  const showModal = useCallback(
    (Content: ReactElement) => {
      setModalContent(() => () => Content)
      open()
    },
    [open]
  )
  const hideModal = useCallback(
    (cb?: () => void) => {
      if (cb && typeof cb === 'function') {
        cb()
        console.error('Modal onClose callback must be a function')
      }
      close()
    },
    [close]
  )
  const value = useMemo(() => {
    return {
      showModal,
      hideModal,
      isOpened: opened
    }
  }, [hideModal, showModal, opened])
  return (
    <ModalContext.Provider value={value}>
      {children}
      <Modal opened={opened} onClose={hideModal}>
        {ModalContent && <ModalContent />}
      </Modal>
    </ModalContext.Provider>
  )
}
