import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest'
import {
  buildCustomPageProxyUrl,
  buildEmbeddedUrl,
  detectTheme,
  isStandaloneCheckoutUrl,
} from '../embedded-url'

describe('embedded-url', () => {
  const originalLocation = window.location

  beforeEach(() => {
    Object.defineProperty(window, 'location', {
      value: {
        origin: 'https://app.example.com',
        href: 'https://app.example.com/user/purchase',
      },
      writable: true,
      configurable: true,
    })
  })

  afterEach(() => {
    Object.defineProperty(window, 'location', {
      value: originalLocation,
      writable: true,
      configurable: true,
    })
    document.documentElement.classList.remove('dark')
    vi.restoreAllMocks()
  })

  it('adds embedded query parameters including locale and source context', () => {
    const result = buildEmbeddedUrl(
      'https://pay.example.com/checkout?plan=pro',
      42,
      'token-123',
      'dark',
      'zh-CN',
    )

    const url = new URL(result)
    expect(url.searchParams.get('plan')).toBe('pro')
    expect(url.searchParams.get('user_id')).toBe('42')
    expect(url.searchParams.get('token')).toBe('token-123')
    expect(url.searchParams.get('theme')).toBe('dark')
    expect(url.searchParams.get('lang')).toBe('zh-CN')
    expect(url.searchParams.get('ui_mode')).toBe('embedded')
    expect(url.searchParams.get('src_host')).toBe('https://app.example.com')
    expect(url.searchParams.get('src_url')).toBe('https://app.example.com/user/purchase')
  })

  it('omits optional params when they are empty', () => {
    const result = buildEmbeddedUrl('https://pay.example.com/checkout', undefined, '', 'light')

    const url = new URL(result)
    expect(url.searchParams.get('theme')).toBe('light')
    expect(url.searchParams.get('ui_mode')).toBe('embedded')
    expect(url.searchParams.has('user_id')).toBe(false)
    expect(url.searchParams.has('token')).toBe(false)
    expect(url.searchParams.has('lang')).toBe(false)
  })

  it('returns original string for invalid url input', () => {
    expect(buildEmbeddedUrl('not a url', 1, 'token')).toBe('not a url')
  })

  it('returns the original checkout url when context appending is disabled', () => {
    const checkoutUrl = 'https://pay.ldxp.cn/shop/BSEJH4PV/4j9om0'

    expect(buildEmbeddedUrl(checkoutUrl, 42, 'token-123', 'dark', 'zh-CN', {
      appendContext: false,
    })).toBe(checkoutUrl)
  })

  it('detects standalone checkout urls', () => {
    expect(isStandaloneCheckoutUrl('https://pay.ldxp.cn/shop/BSEJH4PV/4j9om0')).toBe(true)
    expect(isStandaloneCheckoutUrl('https://example.com/checkout/session')).toBe(true)
    expect(isStandaloneCheckoutUrl('https://example.com/docs')).toBe(false)
  })

  it('builds custom page proxy urls', () => {
    expect(buildCustomPageProxyUrl('migrated_purchase_subscription')).toBe(
      '/api/v1/custom-pages/migrated_purchase_subscription/proxy',
    )
    expect(buildCustomPageProxyUrl('pay/item 1')).toBe('/api/v1/custom-pages/pay%2Fitem%201/proxy')
  })

  it('detects dark mode from document root class', () => {
    document.documentElement.classList.add('dark')
    expect(detectTheme()).toBe('dark')
  })
})
