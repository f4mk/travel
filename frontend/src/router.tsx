import { Suspense } from 'react'
import {
  createBrowserRouter,
  createRoutesFromElements,
  Route,
} from 'react-router-dom'

import { MainLayout } from '#/components/layout/MainLayout'
import { Page } from '#/components/layout/Page'
import { CenteredLoader } from '#/components/ui/CenteredLoader'
import { lazy } from '#/utils/lazy'

const { BlogPage } = lazy(() => import('#/pages/BlogPage'))
const { IndexPage } = lazy(() => import('#/pages/IndexPage'))
const { MapPage } = lazy(() => import('#/pages/MapPage'))
const { LandingPage } = lazy(() => import('#/pages/LandingPage'))

export const router = createBrowserRouter(
  createRoutesFromElements(
    <Route path="/" element={<MainLayout />}>
      <Route index element={<LandingPage />} />
      <Route path="/app" element={<Page />}>
        <Route
          index
          element={
            <Suspense fallback={<CenteredLoader />}>
              <IndexPage />
            </Suspense>
          }
        />
        <Route
          path="map"
          element={
            <Suspense fallback={<CenteredLoader />}>
              <MapPage />
            </Suspense>
          }
        />
        <Route
          path="blog"
          element={
            <Suspense fallback={<CenteredLoader />}>
              <BlogPage />
            </Suspense>
          }
        />
      </Route>
    </Route>
  )
)
