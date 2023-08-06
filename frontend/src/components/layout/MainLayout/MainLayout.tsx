import { Content } from './components/Content'
import { Footer } from './components/Footer'
import { Header } from './components/Header/Header'
import { Meta } from './components/Meta/Meta'
import * as S from './styled'
import { Props } from './types'

export const MainLayout = ({ children }: Props) => {
  return (
    <S.Section>
      <Meta />
      <Header />
      <Content>{children}</Content>
      <Footer />
    </S.Section>
  )
}
