import * as S from './styled'
import { Props } from './types'
export const MarkerContainer = ({ children }: Props) => {
  return (
    <S.Container>
      <S.Pinner>{children}</S.Pinner>
    </S.Container>
  )
}
