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

const { Blog } = lazy(() => import('#/pages/Blog'))
const { Index } = lazy(() => import('#/pages/Index'))
const { Map } = lazy(() => import('#/pages/Map'))
const { Landing } = lazy(() => import('#/pages/Landing'))
const { VerifyAccount } = lazy(() => import('#/pages/VerifyAccount'))
const { ResetPassword } = lazy(() => import('#/pages/ResetPassword'))
const { ConfirmCreateUser } = lazy(() => import('#/pages/ConfirmCreateUser'))
const { NotFound } = lazy(() => import('#/pages/NotFound'))

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
              <Landing />
            </Suspense>
          }
        />
        <Route
          path={ROUTES.USER_VERIFY}
          element={
            <Suspense fallback={<CenteredLoader />}>
              <VerifyAccount />
            </Suspense>
          }
        />
        <Route
          path={ROUTES.USER_CREATE}
          element={
            <Suspense fallback={<CenteredLoader />}>
              <ConfirmCreateUser />
            </Suspense>
          }
        />
        <Route
          path={ROUTES.PASSWORD_RESET}
          element={
            <Suspense fallback={<CenteredLoader />}>
              <ResetPassword />
            </Suspense>
          }
        />
        <Route path={ROUTES.APP.ROOT} element={<Page />}>
          <Route
            index
            element={
              <Suspense fallback={<CenteredLoader />}>
                <Index />
              </Suspense>
            }
          />
          <Route
            path={ROUTES.APP.MAP}
            element={
              <Suspense fallback={<CenteredLoader />}>
                <Map />
              </Suspense>
            }
          />
          <Route
            path={ROUTES.APP.BLOG}
            element={
              <Suspense fallback={<CenteredLoader />}>
                <Blog />
              </Suspense>
            }
          />
        </Route>
        <Route
          path="*"
          element={
            <Suspense fallback={<CenteredLoader />}>
              <NotFound />
            </Suspense>
          }
        />
      </Route>
    )
  )
