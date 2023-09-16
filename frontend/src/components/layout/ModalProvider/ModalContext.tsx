import { createContext, ReactElement } from 'react'

export const ModalContext = createContext({
  // eslint-disable-next-line
  showModal: (_: ReactElement, __?: Record<string,string>) => {},
  // eslint-disable-next-line
  hideModal: (_?:() => void) => {},
})
