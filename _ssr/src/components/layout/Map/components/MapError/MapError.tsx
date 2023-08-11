import * as S from './styled'

export type MapErrorProps = {
  message: string
}

export const MapError = ({ message }: MapErrorProps) => {
  return <S.Div>{message}</S.Div>
}
