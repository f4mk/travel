export const getTokenFromHeader = (authHeader?: Nullable<string>) => {
  if (!authHeader) {
    throw new Error('Authorization failed: header is missing')
  }
  const [_, newToken] = authHeader.split(' ')
  if (!newToken) {
    throw new Error('Authorization failed: token is missing')
  }
  return newToken
}
