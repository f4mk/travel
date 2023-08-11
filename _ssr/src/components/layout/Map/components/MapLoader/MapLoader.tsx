import { LoadingOverlay } from '@mantine/core'

import { Props } from './types'

export const MapLoader = ({ isLoading }: Props) => {
  return <LoadingOverlay visible={isLoading} />
}
