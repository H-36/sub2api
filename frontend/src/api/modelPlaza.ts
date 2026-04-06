import { apiClient } from './client'

export interface ModelPlazaSummary {
  platform_count: number
  group_count: number
  model_count: number
}

export interface ModelPlazaModel {
  name: string
  billing_mode: string
  input_price_1m: number | null
  output_price_1m: number | null
  cache_write_price_1m: number | null
  cache_read_price_1m: number | null
}

export interface ModelPlazaGroup {
  id: number
  name: string
  platform: string
  rate_multiplier: number
  model_count: number
  models: ModelPlazaModel[]
}

export interface ModelPlazaPlatform {
  platform: string
  label: string
  group_count: number
  groups: ModelPlazaGroup[]
}

export interface ModelPlazaResponse {
  summary: ModelPlazaSummary
  platforms: ModelPlazaPlatform[]
}

export async function getModelPlaza(): Promise<ModelPlazaResponse> {
  const { data } = await apiClient.get<ModelPlazaResponse>('/model-plaza')
  return data
}

export const modelPlazaAPI = {
  getModelPlaza
}

export default modelPlazaAPI
