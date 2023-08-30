import { useCallback, useState } from 'react'
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
    googleMapsApiKey: import.meta.env.VITE_MAP_API_KEY
  })
  const [map, setMap] = useState<google.maps.Map | undefined>(undefined)
  const [zoom, setZoom] = useState(10)
  const center = useCurrentLocation()

  const onLoad = useCallback((map: google.maps.Map) => {
    setMap(map)
  }, [])

  const onUnmount = useCallback(() => {
    setMap(undefined)
  }, [])

  if (loadError) {
    return <MapError message={loadError.message} />
  }

  return isLoaded ? (
    <GoogleMap
      mapContainerStyle={mapContainerStyles}
      center={center}
      zoom={zoom}
      onLoad={onLoad}
      onUnmount={onUnmount}
      options={defaultMapOptions}
      clickableIcons={true}
      mapTypeId={'terrain'}
      onZoomChanged={() => {
        setZoom(map?.getZoom() || zoom)
      }}
    >
      ({map && children})
    </GoogleMap>
  ) : (
    <MapLoader isLoading={!isLoaded} />
  )
}
