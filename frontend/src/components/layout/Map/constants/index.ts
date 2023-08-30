export const defaultMapOptions: google.maps.MapOptions = {
  disableDefaultUI: true,
  minZoom: 3,
  zoomControl: true,
  restriction: {
    strictBounds: true,
    latLngBounds: {
      north: 85,
      south: -85,
      west: -180,
      east: 180
    }
  }
}

export const mapContainerStyles = { width: '100%', height: '100%' }
