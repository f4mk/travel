import {
  createBrowserRouter,
  createRoutesFromElements,
  Route
} from 'react-router-dom'

import { MainLayout } from '#/components/layout/MainLayout'
import { BlogPage, IndexPage, MapPage } from '#/pages'

export const router = createBrowserRouter(
  createRoutesFromElements(
    <Route path="/" element={<MainLayout />}>
      <Route index element={<IndexPage />} />
      <Route path="map" element={<MapPage />} />
      <Route path="blog" element={<BlogPage />} />
    </Route>
  )
)
