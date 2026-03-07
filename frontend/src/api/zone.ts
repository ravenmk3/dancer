import type { Zone, CreateZoneForm } from '@/types/zone'
import request from '@/utils/request'

export const getZoneListApi = () => {
  return request.post<{ zones: Zone[] }>('/dns/zones/list')
}

export const getZoneApi = (zone: string) => {
  return request.post<{ zone: Zone }>('/dns/zones/get', { zone })
}

export const createZoneApi = (data: CreateZoneForm) => {
  return request.post<{ zone: Zone }>('/dns/zones/create', data)
}

export const updateZoneApi = (zone: string) => {
  return request.post<{ zone: Zone }>('/dns/zones/update', { zone })
}

export const deleteZoneApi = (zone: string) => {
  return request.post('/dns/zones/delete', { zone })
}
