export type RequestArgs<T> = {
  url: string
  method?: 'GET' | 'POST' | 'PUT' | 'DELETE'
  body?: T
}
