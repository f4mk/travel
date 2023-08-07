import { Outlet } from 'react-router-dom'

import { Content } from './components/Content'
import { Footer } from './components/Footer'
import { Header } from './components/Header/Header'
import { Meta } from './components/Meta/Meta'
import * as S from './styled'

export const MainLayout = () => {
  return (
    <S.Section>
      <Meta />
      <Header />
      <Content>
        <Outlet />
      </Content>
      <Footer />
    </S.Section>
  )
}
