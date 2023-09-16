import { Component, ErrorInfo } from 'react'

import { isEqualArrays } from '#/utils'

import { Props, State } from './types'

export class ErrorBoundary extends Component<Props, State> {
  constructor(props: Props) {
    super(props)
    this.state = {
      error: undefined,
      resetList: props.resetList,
    }
  }
  static getDerivedStateFromError(error: Error, prevState: State) {
    return { ...prevState, error }
  }

  static getDerivedStateFromProps(props: Props, prevState: State) {
    const state = { ...prevState }
    if (!isEqualArrays(props.resetList, prevState.resetList)) {
      state.error = undefined
      state.resetList = props.resetList
    }
    return state
  }

  shouldComponentUpdate(nextProps: Props, nextState: State) {
    if (
      this.state.error !== nextState.error ||
      this.props.children !== nextProps.children
    ) {
      return true
    }
    return (
      Boolean(nextState.error) &&
      !isEqualArrays(this.props.resetList, nextProps.resetList)
    )
  }

  componentDidCatch(error: Error, info: ErrorInfo) {
    // TODO: add logging
    console.error(error.message, info.componentStack)
  }

  reset() {
    this.props.onReset?.()
    this.setState({ resetList: this.props.resetList, error: undefined })
  }

  render() {
    if (this.state.error) {
      const Fallback = this.props.fallback
      return <Fallback error={this.state.error} reset={() => this.reset()} />
    }
    return this.props.children
  }
}
