export const isEqualArrays = (a: any[], b: any[]) => {
  if (a.length !== b.length) {
    return false
  }
  return a.every((item, idx) => Object.is(item, b[idx]))
}
