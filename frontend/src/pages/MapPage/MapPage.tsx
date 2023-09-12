import { Map } from '#/components/layout/Map'
import { MapCard } from '#/components/ui/map/MapCard'
import { Pin } from '#/components/ui/map/Pin'

import * as S from './styled'

const pins = [
  {
    id: '1',
    lat: 40,
    lng: 41,
  },
  {
    id: '2',
    lat: 40.02,
    lng: 41.02,
  },
]

const cards = [
  {
    id: '3',
    lat: 39.025,
    lng: 42.015,
  },
  {
    id: '4',
    lat: 44.12,
    lng: 45.12,
  },
]

export const MapPage = () => {
  return (
    <S.Div>
      <Map>
        <>
          {pins.map((pin) => (
            <Pin key={pin.id} {...pin} />
          ))}
          {cards.map((card) => (
            <MapCard key={card.id} {...card} />
          ))}
        </>
      </Map>
    </S.Div>
  )
}
