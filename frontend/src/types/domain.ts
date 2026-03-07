export interface Domain {
  zone: string
  domain: string
  name: string
  ips: string[]
  ttl: number
  record_count: number
  created_at: number
  updated_at: number
}

export interface CreateDomainForm {
  zone: string
  domain: string
  ips: string[]
  ttl: number
}
