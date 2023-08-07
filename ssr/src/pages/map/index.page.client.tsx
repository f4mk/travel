import { Map } from '#/components/layout/Map'
import { MapCard } from '#/components/ui/map/MapCard'
import { Pin } from '#/components/ui/map/Pin'

import * as S from './styled'

const pins = [
  {
    id: '1',
    lat: 40,
    lng: 41
  },
  {
    id: '2',
    lat: 40.02,
    lng: 41.02
  }
]

const cards = [
  {
    id: '3',
    lat: 39.025,
    lng: 42.015
  },
  {
    id: '4',
    lat: 44.12,
    lng: 45.12
  }
]

export const Page = () => {
  return (
    <S.Div>
      {/* <ClientOnly
        component={() =>
          import('#/components/layout/Map').then((module) => ({
            default: module.Map
          }))
        } */}
      <Map
        markers={(props) => (
          <>
            {pins.map(({ id, ...rest }) => (
              <Pin key={id} {...rest} {...props} />
            ))}
            {cards.map(({ id, ...rest }) => (
              <MapCard key={id} {...rest} {...props} />
            ))}
          </>
        )}
      />
    </S.Div>
  )
}
