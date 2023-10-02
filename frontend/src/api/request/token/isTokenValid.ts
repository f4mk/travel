import { Temporal } from '@js-temporal/polyfill'

import { JwtPayload } from './types'

const isExpired = (exp: number) => {
  const expTime = Temporal.Instant.fromEpochSeconds(exp)
  if (Temporal.Instant.compare(expTime, Temporal.Now.instant()) === 1) {
    return false
  }
  return true
}

const decodeJwt = (token: string) => {
  try {
    const [_, payloadEnc, __] = token.split('.')
    const payloadDec = atob(payloadEnc)
    const payload = JSON.parse(payloadDec) as JwtPayload
    return payload
  } catch (e) {
    console.error('Failed to decode token', e)
    return null
  }
}

export const isTokenValid = (token: string) => {
  const payload = decodeJwt(token)
  if (!payload) {
    return false
  }
  try {
    const { exp } = payload
    return !isExpired(exp)
  } catch {
    return false
  }
}
