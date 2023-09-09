import { paths } from '#/types/external/auth'
export type LoginRequest =
  paths['/auth/login']['post']['requestBody']['content']['application/json']
export type LoginResponse =
  paths['/auth/login']['post']['responses']['201']['content']['application/json']
export type LoginError =
  paths['/auth/login']['post']['responses']['400']['content']['application/json']

export type LogoutRequest =
  paths['/auth/logout']['post']['requestBody']['content']['application/json']
export type LogoutResponse =
  paths['/auth/logout']['post']['responses']['201']['content']['application/json']
export type LogoutError =
  paths['/auth/logout']['post']['responses']['400']['content']['application/json']

export type LogoutAllRequest =
  paths['/auth/logout/all']['post']['requestBody']['content']['application/json']
export type LogoutAllResponse =
  paths['/auth/logout/all']['post']['responses']['201']['content']['application/json']
export type LogoutAllError =
  paths['/auth/logout/all']['post']['responses']['400']['content']['application/json']

export type PasswordResetRequest =
  paths['/auth/password/reset']['post']['requestBody']['content']['application/json']
export type PasswordResetResponse =
  paths['/auth/password/reset']['post']['responses']['201']['content']['application/json']
export type PasswordResetError =
  paths['/auth/password/reset']['post']['responses']['400']['content']['application/json']
export type PasswordResetSubmitRequest =
  paths['/auth/password/reset/submit']['post']['requestBody']['content']['application/json']
export type PasswordResetSubmitResponse =
  paths['/auth/password/reset/submit']['post']['responses']['201']['content']['application/json']
export type PasswordResetSubmitError =
  paths['/auth/password/reset/submit']['post']['responses']['400']['content']['application/json']

export type PasswordChangeRequest =
  paths['/auth/password/change']['post']['requestBody']['content']['application/json']
export type PasswordChangeResponse =
  paths['/auth/password/change']['post']['responses']['201']['content']['application/json']
export type PasswordChangeError =
  paths['/auth/password/change']['post']['responses']['400']['content']['application/json']
