import { Outlet } from 'react-router-dom'

// import { ModalProvider } from '#/components/layout/ModalProvider'
import { Meta } from './Meta'
import * as S from './styled'
import { Props } from './types'

export const MainLayout = ({ ModalProvider }: Props) => {
  return (
    <S.Section>
      <Meta />
      <ModalProvider>
        <Outlet />
      </ModalProvider>
    </S.Section>
  )
}
