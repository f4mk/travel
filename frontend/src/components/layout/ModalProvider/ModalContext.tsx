import { createContext, ReactElement } from 'react'

export const ModalContext = createContext({
  // eslint-disable-next-line
  showModal: (_: ReactElement) => {},
  // eslint-disable-next-line
  hideModal: (_?:() => void) => {},
  isOpened: false
})
