import { EFormView } from '#/components/ui/ProfileMenu'

export type Props = {
  opened: boolean
  activeTab: EFormView
  setActiveTab: (tab: EFormView) => void
  onClose: () => void
}
