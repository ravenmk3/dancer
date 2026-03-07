export interface Zone {
  zone: string
  record_count: number
  created_at: number
  updated_at: number
}

export interface CreateZoneForm {
  zone: string
}
