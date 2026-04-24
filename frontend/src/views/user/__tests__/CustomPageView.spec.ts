import { describe, expect, it, vi, beforeEach } from 'vitest'
import { shallowMount } from '@vue/test-utils'
import CustomPageView from '../CustomPageView.vue'

const routeState = vi.hoisted(() => ({
  params: { id: 'migrated_purchase_subscription' },
}))

const routerReplace = vi.hoisted(() => vi.fn())

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
    useRouter: () => ({
      replace: routerReplace,
    }),
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
  beforeEach(() => {
    routerReplace.mockReset()
    appStoreState.publicSettingsLoaded = true
    appStoreState.fetchPublicSettings.mockReset()
  })

  it('redirects standalone checkout custom pages to the internal purchase view', () => {
    shallowMount(CustomPageView)

    expect(routerReplace).toHaveBeenCalledWith('/purchase')
  })
})
