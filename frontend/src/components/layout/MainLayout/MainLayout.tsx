import { Outlet } from 'react-router-dom'

import { Meta } from './Meta'
import * as S from './styled'

export const MainLayout = () => {
  return (
    <S.Section>
      <Meta />
      <Outlet />
    </S.Section>
  )
}
