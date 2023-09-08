export type JwtPayload = {
  sub: string
  exp: number
  iat: number
  jti: string
  roles: string[]
  token_version: number
}
