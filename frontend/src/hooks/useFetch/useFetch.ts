const fetcher = (...args: Parameters<typeof fetch>) =>
  fetch(...args).then((res) => res.json())

export const useFetch = (url: string) => {
  return fetcher(url)
}
