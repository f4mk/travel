import { Suspense } from 'react'
import { Outlet } from 'react-router-dom'

import { CenteredLoader } from '#/components/ui/CenteredLoader'

import { Errors } from '../Errors'

import { Content } from './components/Content'
import { Footer } from './components/Footer'
import { Header } from './components/Header'
import * as S from './styled'
export const Page = () => {
  return (
    <S.Div>
      <Errors>
        <Suspense fallback={<CenteredLoader />}>
          <Header />
        </Suspense>
      </Errors>
      <Content>
        <Outlet />
      </Content>
      <Footer />
    </S.Div>
  )
}
