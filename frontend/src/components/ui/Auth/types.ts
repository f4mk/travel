export type Props = {
  activeTab: EFormView
  onClose: () => void
}
export enum EFormView {
  SIGN_IN = 'signIn',
  SIGN_UP = 'signUp',
}
