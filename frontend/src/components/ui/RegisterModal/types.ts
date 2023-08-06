export type Props = {
  opened: boolean
  onClose: () => void
  onSwitch: () => void
}
export type FormValues = {
  username: string
  name: string
  lastname: string
  email: string
  password: string
  passwordRepeat: string
}
