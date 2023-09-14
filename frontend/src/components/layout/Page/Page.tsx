import { Suspense } from 'react'
import { Outlet } from 'react-router-dom'

import { CenteredLoader } from '#/components/ui/CenteredLoader'

import { Errors } from '../Errors'

import { Content } from './components/Content'
import { Footer } from './components/Footer'
import { Header } from './components/Header'

export const Page = () => {
  return (
    <div>
      <Errors>
        <Suspense fallback={<CenteredLoader />}>
          <Header />
        </Suspense>
      </Errors>
      <Content>
        <Outlet />
      </Content>
      <Footer />
    </div>
  )
}
