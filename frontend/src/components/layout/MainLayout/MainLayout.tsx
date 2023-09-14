import { Outlet } from 'react-router-dom'

import { ModalProvider } from '#/components/layout/ModalProvider'

import { Meta } from './Meta'
import * as S from './styled'

export const MainLayout = () => {
  return (
    <S.Section>
      <Meta />
      <ModalProvider>
        <Outlet />
      </ModalProvider>
    </S.Section>
  )
}
