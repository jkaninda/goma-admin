export interface Pageable {
  current_page: number
  size: number
  total_pages: number
  total_elements: number
  empty: boolean
}

export interface PaginatedResponse<T> {
  success: boolean
  data: T[]
  pageable: Pageable
}
