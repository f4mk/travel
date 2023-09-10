import { paths } from '#/types/external/user'

export type CreateUserRequest =
  paths['/user']['post']['requestBody']['content']['application/json']
export type CreateUserResponse =
  paths['/user']['post']['responses']['201']['content']['application/json']

export type CreateUserError =
  paths['/user']['post']['responses']['400']['content']['application/json']

export type DeleteUserRequest =
  paths['/user']['delete']['requestBody']['content']['application/json']
export type DeleteUserError =
  paths['/user']['delete']['responses']['400']['content']['application/json']
export type DeleteUserResponse =
  paths['/user']['delete']['responses']['200']['content']['application/json']

export type UpdateUserRequest =
  paths['/user']['put']['requestBody']['content']['application/json']
export type UpdateUserError =
  paths['/user']['put']['responses']['400']['content']['application/json']
export type UpdateUserResponse =
  paths['/user']['put']['responses']['200']['content']['application/json']

export type GetUserRequest =
  paths['/user/{id}']['get']['parameters']['path']['id']
export type GetUserError =
  paths['/user/{id}']['get']['responses']['401']['content']['application/json']
export type GetUserResponse =
  paths['/user/{id}']['get']['responses']['200']['content']['application/json']

export type GetMeRequest = paths['/user/me']['get']
export type GetMeError =
  paths['/user/me']['get']['responses']['401']['content']['application/json']
export type GetMeResponse =
  paths['/user/me']['get']['responses']['200']['content']['application/json']
