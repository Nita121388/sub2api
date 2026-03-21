import { apiClient } from './client'
import type {
  FetchOptions,
  PaginatedResponse,
  RequestLog,
  RequestLogQueryParams
} from '@/types'

export async function query(
  params: RequestLogQueryParams,
  options: FetchOptions = {}
): Promise<PaginatedResponse<RequestLog>> {
  const { data } = await apiClient.get<PaginatedResponse<RequestLog>>('/request-logs', {
    ...options,
    params
  })
  return data
}

export async function getById(id: number, options: FetchOptions = {}): Promise<RequestLog> {
  const { data } = await apiClient.get<RequestLog>(`/request-logs/${id}`, options)
  return data
}

export const requestLogsAPI = {
  query,
  getById
}
