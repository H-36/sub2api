/**
 * Model Plaza API（公开端点，可匿名访问）
 *
 * 上游新视图以分组为中心：分组信息 + 模型渠道定价 + LiteLLM 官方参考价。
 * 本地历史 /models 视图仍使用 summary/platforms 与扁平 1M 单价字段，后端会同时返回，
 * 因此这里保留兼容字段，避免旧入口在合并后空白。
 */

import { apiClient } from './client'
import type { UserSupportedModelPricing } from './channels'

export interface ModelPlazaSummary {
  platform_count: number
  group_count: number
  model_count: number
}

/** LiteLLM 官方参考价（USD per token，字段缺失 = 官方数据未覆盖）。 */
export interface PlazaOfficialPricing {
  input_price: number | null
  output_price: number | null
  /** 5m 缓存写入（= LiteLLM cache_creation）。 */
  cache_write_price: number | null
  /** 1h 缓存写入（LiteLLM cache_creation_above_1hr），多数模型缺失。 */
  cache_write_1h_price?: number | null
  cache_read_price: number | null
}

export interface ModelPlazaModel {
  name: string
  platform: string
  pricing: UserSupportedModelPricing | null
  official_pricing: PlazaOfficialPricing | null

  /** 兼容本地历史 /models 视图：计费模式与预换算的 $/1M token。 */
  billing_mode: string
  input_price_1m: number | null
  output_price_1m: number | null
  cache_write_price_1m: number | null
  cache_read_price_1m: number | null
}

export type PlazaModel = ModelPlazaModel

export interface ModelPlazaGroup {
  id: number
  name: string
  description: string
  platform: string
  /** 'standard' | 'subscription' */
  subscription_type: string
  rate_multiplier: number
  /** 登录且管理员为该用户配了专属倍率时返回；生效倍率 = user_rate ?? rate_multiplier。 */
  user_rate_multiplier?: number
  peak_rate_enabled: boolean
  peak_start: string
  peak_end: string
  peak_rate_multiplier: number
  is_exclusive: boolean
  /** 兼容本地历史 /models 视图。 */
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
  /** 管理员配置的全局价格说明（Markdown）。 */
  description: string
  groups: ModelPlazaGroup[]

  /** 兼容本地历史 /models 视图。 */
  summary: ModelPlazaSummary
  platforms: ModelPlazaPlatform[]
}

/** 获取模型广场数据。开关未启用时后端返回 404。 */
export async function getModelPlaza(options?: { signal?: AbortSignal }): Promise<ModelPlazaResponse> {
  const { data } = await apiClient.get<ModelPlazaResponse>('/model-plaza', {
    signal: options?.signal
  })
  return data
}

export const modelPlazaAPI = { getModelPlaza }

export default modelPlazaAPI
