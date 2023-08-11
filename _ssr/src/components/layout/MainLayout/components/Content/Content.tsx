import { PropsWithChildren } from 'react'

import * as S from './styled'

export type ContentProps = PropsWithChildren
export const Content = ({ children }: ContentProps) => {
  return <S.Main>{children}</S.Main>
}
