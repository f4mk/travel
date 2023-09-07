export const useIsMobile = () => {
  if (window.innerWidth <= 768) {
    return true
  }

  return false
}
