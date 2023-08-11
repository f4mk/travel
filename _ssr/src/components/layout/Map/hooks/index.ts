const DEFAULT_CENTER = {
  lat: -3.745,
  lng: -38.523
}

export const useCurrentLocation = () => {
  let pos
  if (navigator.geolocation) {
    navigator.geolocation.getCurrentPosition(
      (position: GeolocationPosition) => {
        pos = {
          lat: position.coords.latitude,
          lng: position.coords.longitude
        }
      }
    )
  }
  return pos || DEFAULT_CENTER
}

export const useClientLanguage = () => {
  return window?.navigator.language || 'en-GB'
}
