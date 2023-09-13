import { Outlet } from 'react-router-dom'

import { Footer } from './components/Footer'
import { Meta } from './components/Meta/Meta'
import * as S from './styled'

export const MainLayout = () => {
  return (
    <S.Section>
      <Meta />
      <Outlet />
      <Footer />
    </S.Section>
  )
}
