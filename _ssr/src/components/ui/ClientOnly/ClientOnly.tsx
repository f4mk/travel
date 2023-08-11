import {
  ComponentType,
  type FC,
  lazy,
  Suspense,
  useEffect,
  useState
} from 'react'

import { Props } from './types'

export const ClientOnly: FC<Props> = ({ component, fallback }) => {
  const [Component, setComponent] = useState<ComponentType<any> | null>(null)

  useEffect(() => {
    setComponent(() => lazy(component))
  }, [component])

  return <Suspense fallback={fallback}>{Component && <Component />}</Suspense>
}
