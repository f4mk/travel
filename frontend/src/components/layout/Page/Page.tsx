import { Outlet } from 'react-router-dom'

import { Content } from './components/Content'
import { Footer } from './components/Footer'
import { Header } from './components/Header'

export const Page = () => {
  return (
    <div>
      <Header />
      <Content>
        <Outlet />
      </Content>
      <Footer />
    </div>
  )
}
