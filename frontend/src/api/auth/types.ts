import { paths } from '#/types/external/auth'
export type LoginRequest =
  paths['/auth/login']['post']['requestBody']['content']['application/json']
export type LoginResponse =
  paths['/auth/login']['post']['responses']['201']['content']['application/json']
export type LoginError =
  paths['/auth/login']['post']['responses']['400']['content']['application/json']
