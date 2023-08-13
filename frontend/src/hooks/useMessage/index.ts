import { useIntl } from 'react-intl'

export const useMessage = () => {
  const intl = useIntl()

  return intl.formatMessage
}
