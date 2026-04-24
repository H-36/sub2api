import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest'
import { shallowMount } from '@vue/test-utils'
import CustomPageView from '../CustomPageView.vue'

const routeState = vi.hoisted(() => ({
  params: { id: 'migrated_purchase_subscription' },
}))

const appStoreState = vi.hoisted(() => ({
  publicSettingsLoaded: true,
  cachedPublicSettings: {
    payment_enabled: true,
    custom_menu_items: [
      {
        id: 'migrated_purchase_subscription',
        label: 'Recharge',
        icon_svg: '',
        url: 'https://pay.ldxp.cn/shop/BSEJH4PV/4j9om0',
        visibility: 'user',
        sort_order: 0,
      },
    ],
  },
  fetchPublicSettings: vi.fn(),
}))

vi.mock('vue-router', async () => {
  const actual = await vi.importActual<typeof import('vue-router')>('vue-router')
  return {
    ...actual,
    useRoute: () => routeState,
  }
})

vi.mock('vue-i18n', async () => {
  const actual = await vi.importActual<typeof import('vue-i18n')>('vue-i18n')
  return {
    ...actual,
    useI18n: () => ({
      t: (key: string) => key,
      locale: { value: 'zh-CN' },
    }),
  }
})

vi.mock('@/stores', () => ({
  useAppStore: () => appStoreState,
}))

vi.mock('@/stores/auth', () => ({
  useAuthStore: () => ({
    user: { id: 42 },
    token: 'token-123',
    isAdmin: false,
  }),
}))

vi.mock('@/stores/adminSettings', () => ({
  useAdminSettingsStore: () => ({
    customMenuItems: [],
  }),
}))

describe('CustomPageView', () => {
  const originalLocation = window.location
  const assign = vi.fn()

  beforeEach(() => {
    appStoreState.publicSettingsLoaded = true
    appStoreState.fetchPublicSettings.mockReset()
    assign.mockReset()
    Object.defineProperty(window, 'location', {
      value: {
        assign,
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
  })

  it('opens standalone checkout custom pages in the current window', () => {
    shallowMount(CustomPageView)

    expect(assign).toHaveBeenCalledWith('https://pay.ldxp.cn/shop/BSEJH4PV/4j9om0')
  })
})
