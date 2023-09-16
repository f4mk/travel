import { Suspense } from 'react'
import {
  createBrowserRouter,
  createRoutesFromElements,
  Route,
} from 'react-router-dom'

import { MainLayout } from '#/components/layout/MainLayout'
import { Page } from '#/components/layout/Page'
import { CenteredLoader } from '#/components/ui/CenteredLoader'
import { lazy } from '#/utils'

import { ERoutes } from './constants/routes'

const { BlogPage } = lazy(() => import('#/pages/BlogPage'))
const { IndexPage } = lazy(() => import('#/pages/IndexPage'))
const { MapPage } = lazy(() => import('#/pages/MapPage'))
const { LandingPage } = lazy(() => import('#/pages/LandingPage'))
const { VerifyPage } = lazy(() => import('#/pages/VerifyPage'))
const { ResetPasswordPage } = lazy(() => import('#/pages/ResetPasswordPage'))
const { ConfirmCreatePage } = lazy(() => import('#/pages/ConfirmCreatePage'))
const { NotFoundPage } = lazy(() => import('#/pages/NotFoundPage'))

export const router = createBrowserRouter(
  createRoutesFromElements(
    <Route path={ERoutes.ROOT} element={<MainLayout />}>
      <Route
        index
        element={
          <Suspense fallback={<CenteredLoader />}>
            <LandingPage />
          </Suspense>
        }
      />
      <Route
        path={ERoutes.USER_VERIFY}
        element={
          <Suspense fallback={<CenteredLoader />}>
            <VerifyPage />
          </Suspense>
        }
      />
      <Route
        path={ERoutes.USER_CREATE}
        element={
          <Suspense fallback={<CenteredLoader />}>
            <ConfirmCreatePage />
          </Suspense>
        }
      />
      <Route
        path={ERoutes.PASSWORD_RESET}
        element={
          <Suspense fallback={<CenteredLoader />}>
            <ResetPasswordPage />
          </Suspense>
        }
      />
      <Route path={ERoutes.APP} element={<Page />}>
        <Route
          index
          element={
            <Suspense fallback={<CenteredLoader />}>
              <IndexPage />
            </Suspense>
          }
        />
        <Route
          path={ERoutes.MAP}
          element={
            <Suspense fallback={<CenteredLoader />}>
              <MapPage />
            </Suspense>
          }
        />
        <Route
          path={ERoutes.BLOG}
          element={
            <Suspense fallback={<CenteredLoader />}>
              <BlogPage />
            </Suspense>
          }
        />
      </Route>
      <Route
        path="*"
        element={
          <Suspense fallback={<CenteredLoader />}>
            <NotFoundPage />
          </Suspense>
        }
      />
    </Route>
  )
)
