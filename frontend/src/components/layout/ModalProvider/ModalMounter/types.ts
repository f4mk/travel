import { PropsWithChildren } from 'react'

export type Props = PropsWithChildren & { onMount: (a: boolean) => void }
