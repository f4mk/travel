import { DEFAULT_CENTER } from './consts'

export const useCurrentLocation = () => {
  let pos
  if (navigator.geolocation) {
    navigator.geolocation.getCurrentPosition(
      (position: GeolocationPosition) => {
        pos = {
          lat: position.coords.latitude,
          lng: position.coords.longitude,
        }
      }
    )
  }
  return pos || DEFAULT_CENTER
}
