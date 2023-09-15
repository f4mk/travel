import { paths } from '#/types/external/user'

export type CreateUserRequest =
  paths['/users']['post']['requestBody']['content']['application/json']
export type CreateUserResponse =
  paths['/users']['post']['responses']['201']['content']['application/json']

export type CreateUserError =
  paths['/users']['post']['responses']['400']['content']['application/json']

export type DeleteUserRequest =
  paths['/users']['delete']['requestBody']['content']['application/json']
export type DeleteUserError =
  paths['/users']['delete']['responses']['400']['content']['application/json']
export type DeleteUserResponse =
  paths['/users']['delete']['responses']['200']['content']['application/json']

export type UpdateUserRequest =
  paths['/users']['put']['requestBody']['content']['application/json']
export type UpdateUserError =
  paths['/users']['put']['responses']['400']['content']['application/json']
export type UpdateUserResponse =
  paths['/users']['put']['responses']['200']['content']['application/json']

export type GetUserRequest =
  paths['/users/{id}']['get']['parameters']['path']['id']
export type GetUserError =
  paths['/users/{id}']['get']['responses']['401']['content']['application/json']
export type GetUserResponse =
  paths['/users/{id}']['get']['responses']['200']['content']['application/json']

export type GetMeRequest = paths['/users/me']['get']
export type GetMeError =
  paths['/users/me']['get']['responses']['401']['content']['application/json']
export type GetMeResponse =
  paths['/users/me']['get']['responses']['200']['content']['application/json']

export type VerifyUserRequest =
  paths['/users/verify']['post']['requestBody']['content']['application/json']
export type VerifyUserError =
  paths['/users/verify']['post']['responses']['403']['content']['application/json']
export type VerifyUserResponse =
  paths['/users/verify']['post']['responses']['201']['content']['application/json']
