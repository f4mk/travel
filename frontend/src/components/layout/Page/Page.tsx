import { Outlet } from 'react-router-dom'

import { Content } from './components/Content'
import { Header } from './components/Header'

export const Page = () => {
  return (
    <div>
      <Header />
      <Content>
        <Outlet />
      </Content>
    </div>
  )
}
