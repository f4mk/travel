import { useCallback, useRef, useState } from 'react'
import { GoogleMap, useLoadScript } from '@react-google-maps/api'

import { useGetLocale } from '#/hooks'

import { MapError } from './components/MapError'
import { MapLoader } from './components/MapLoader'
import { defaultMapOptions, mapContainerStyles } from './constants'
import { useCurrentLocation } from './hooks'
import { Props } from './types'

export const Map = ({ children }: Props) => {
  const { isLoaded, loadError } = useLoadScript({
    id: 'google-map-script',
    language: useGetLocale(),
    googleMapsApiKey: import.meta.env.VITE_MAP_API_KEY,
  })
  const mapRef = useRef<Nullable<google.maps.Map>>(null)
  const [zoom, setZoom] = useState(10)
  const center = useCurrentLocation()

  const handleLoad = useCallback((map: google.maps.Map) => {
    mapRef.current = map
  }, [])

  const handleUnmount = useCallback(() => {
    mapRef.current = null
  }, [])

  const handleZoomChanged = useCallback(() => {
    setZoom((prev) => mapRef.current?.getZoom() || prev)
  }, [])

  if (loadError) {
    return <MapError message={loadError.message} />
  }

  return isLoaded ? (
    <GoogleMap
      mapContainerStyle={mapContainerStyles}
      center={center}
      zoom={zoom}
      onLoad={handleLoad}
      onUnmount={handleUnmount}
      options={defaultMapOptions}
      clickableIcons={true}
      mapTypeId={'terrain'}
      onZoomChanged={handleZoomChanged}
      onClick={(e) => console.log(e)}
    >
      {children}
    </GoogleMap>
  ) : (
    <MapLoader isLoading={!isLoaded} />
  )
}
