import { ComponentType, Suspense } from 'react'
import {
  createBrowserRouter,
  createRoutesFromElements,
  Route,
} from 'react-router-dom'

import { MainLayout } from '#/components/layout/MainLayout'
import { Page } from '#/components/layout/Page'
import { CenteredLoader } from '#/components/ui/CenteredLoader'
import { lazy } from '#/utils'

import { ROUTES } from './constants/routes'

const { BlogPage } = lazy(() => import('#/pages/BlogPage'))
const { IndexPage } = lazy(() => import('#/pages/IndexPage'))
const { MapPage } = lazy(() => import('#/pages/MapPage'))
const { LandingPage } = lazy(() => import('#/pages/LandingPage'))
const { VerifyPage } = lazy(() => import('#/pages/VerifyPage'))
const { ResetPasswordPage } = lazy(() => import('#/pages/ResetPasswordPage'))
const { ConfirmCreatePage } = lazy(() => import('#/pages/ConfirmCreatePage'))
const { NotFoundPage } = lazy(() => import('#/pages/NotFoundPage'))

export const createRouter = (ModalProvider: ComponentType) =>
  createBrowserRouter(
    createRoutesFromElements(
      <Route
        path={ROUTES.ROOT}
        element={<MainLayout ModalProvider={ModalProvider} />}
      >
        <Route
          index
          element={
            <Suspense fallback={<CenteredLoader />}>
              <LandingPage />
            </Suspense>
          }
        />
        <Route
          path={ROUTES.USER_VERIFY}
          element={
            <Suspense fallback={<CenteredLoader />}>
              <VerifyPage />
            </Suspense>
          }
        />
        <Route
          path={ROUTES.USER_CREATE}
          element={
            <Suspense fallback={<CenteredLoader />}>
              <ConfirmCreatePage />
            </Suspense>
          }
        />
        <Route
          path={ROUTES.PASSWORD_RESET}
          element={
            <Suspense fallback={<CenteredLoader />}>
              <ResetPasswordPage />
            </Suspense>
          }
        />
        <Route path={ROUTES.APP.ROOT} element={<Page />}>
          <Route
            index
            element={
              <Suspense fallback={<CenteredLoader />}>
                <IndexPage />
              </Suspense>
            }
          />
          <Route
            path={ROUTES.APP.MAP}
            element={
              <Suspense fallback={<CenteredLoader />}>
                <MapPage />
              </Suspense>
            }
          />
          <Route
            path={ROUTES.APP.BLOG}
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
