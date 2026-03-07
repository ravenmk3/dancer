import type { Domain, CreateDomainForm } from '@/types/domain'
import request from '@/utils/request'

export const getDomainListApi = (zone: string) => {
  return request.post<{ domains: Domain[] }>('/dns/domains/list', { zone })
}

export const getDomainApi = (zone: string, domain: string) => {
  return request.post<{ domain: Domain }>('/dns/domains/get', { zone, domain })
}

export const createDomainApi = (data: CreateDomainForm) => {
  return request.post<{ domain: Domain }>('/dns/domains/create', data)
}

export const updateDomainApi = (data: CreateDomainForm) => {
  return request.post<{ domain: Domain }>('/dns/domains/update', data)
}

export const deleteDomainApi = (zone: string, domain: string) => {
  return request.post('/dns/domains/delete', { zone, domain })
}
