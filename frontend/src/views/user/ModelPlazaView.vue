<template>
  <AppLayout>
    <div class="mx-auto w-full max-w-[1400px] space-y-6">
      <div v-if="loading" class="flex justify-center py-12">
        <div
          class="h-8 w-8 animate-spin rounded-full border-2 border-primary-500 border-t-transparent"
        ></div>
      </div>

      <template v-else-if="data">
        <section class="grid gap-4 md:grid-cols-3">
          <article
            v-for="card in summaryCards"
            :key="card.key"
            class="relative overflow-hidden rounded-[28px] border border-white/70 bg-white/90 p-5 shadow-[0_18px_50px_-30px_rgba(15,23,42,0.45)] backdrop-blur dark:border-white/10 dark:bg-dark-900/85"
          >
            <div :class="['absolute inset-x-0 top-0 h-1.5', card.accent]"></div>
            <div class="flex items-start justify-between gap-4">
              <div>
                <p class="text-sm font-medium text-gray-500 dark:text-dark-400">
                  {{ card.label }}
                </p>
                <p class="mt-3 text-3xl font-semibold tracking-tight text-gray-900 dark:text-white">
                  {{ card.value }}
                </p>
              </div>
              <div :class="['flex h-11 w-11 items-center justify-center rounded-2xl', card.iconBg]">
                <Icon :name="card.icon" size="md" :class="card.iconColor" />
              </div>
            </div>
          </article>
        </section>

        <section
          class="overflow-hidden rounded-[32px] border border-white/70 bg-white/90 shadow-[0_22px_60px_-35px_rgba(15,23,42,0.45)] backdrop-blur dark:border-white/10 dark:bg-dark-900/85"
        >
          <div class="border-b border-gray-100/80 px-5 py-4 dark:border-dark-700/80">
            <div class="flex flex-wrap items-center justify-between gap-3">
              <div>
                <h2 class="text-lg font-semibold text-gray-900 dark:text-white">
                  {{ t('modelPlaza.platformsTitle') }}
                </h2>
                <p class="text-sm text-gray-500 dark:text-dark-400">
                  {{ t('modelPlaza.platformsHint') }}
                </p>
              </div>
              <button
                type="button"
                class="rounded-full border border-gray-200 px-3 py-1.5 text-sm font-medium text-gray-600 transition hover:border-primary-200 hover:text-primary-600 dark:border-dark-700 dark:text-dark-300 dark:hover:border-primary-500/40 dark:hover:text-primary-300"
                @click="togglePlatformFilter(null)"
              >
                {{ t('modelPlaza.allPlatforms') }}
              </button>
            </div>
          </div>

          <div class="grid gap-3 p-5 md:grid-cols-2 xl:grid-cols-4">
            <button
              v-for="platform in platformCards"
              :key="platform.platform"
              type="button"
              :class="[
                'group rounded-[28px] border p-4 text-left transition-all duration-200',
                selectedPlatform === platform.platform
                  ? 'border-primary-300 bg-primary-50 shadow-[0_16px_36px_-28px_rgba(14,165,233,0.9)] dark:border-primary-500/50 dark:bg-primary-500/10'
                  : 'border-gray-200/80 bg-gradient-to-br from-white to-gray-50 hover:-translate-y-0.5 hover:border-gray-300 hover:shadow-[0_16px_36px_-30px_rgba(15,23,42,0.45)] dark:border-dark-700 dark:from-dark-900 dark:to-dark-800 dark:hover:border-dark-500'
              ]"
              @click="togglePlatformFilter(platform.platform)"
            >
              <div class="flex items-start justify-between gap-4">
                <div
                  :class="[
                    'flex h-12 w-12 items-center justify-center rounded-2xl',
                    platformBadgeClass(platform.platform)
                  ]"
                >
                  <PlatformIcon :platform="normalizePlatform(platform.platform)" size="lg" />
                </div>
                <span
                  :class="[
                    'rounded-full px-2.5 py-1 text-xs font-medium',
                    selectedPlatform === platform.platform
                      ? 'bg-white text-primary-700 dark:bg-primary-500/15 dark:text-primary-200'
                      : 'bg-gray-100 text-gray-600 dark:bg-dark-700 dark:text-dark-200'
                  ]"
                >
                  {{ platform.groupCount }} {{ t('modelPlaza.groupUnit') }}
                </span>
              </div>

              <div class="mt-5 space-y-2">
                <h3 class="text-base font-semibold text-gray-900 dark:text-white">
                  {{ platform.label }}
                </h3>
                <p class="text-sm text-gray-500 dark:text-dark-400">
                  {{ t('modelPlaza.platformCardModels', { count: platform.modelCount }) }}
                </p>
              </div>
            </button>
          </div>
        </section>

        <section class="space-y-4">
          <div class="flex flex-wrap items-center justify-between gap-3">
            <div>
              <h2 class="text-xl font-semibold text-gray-900 dark:text-white">
                {{ t('modelPlaza.groupsTitle') }}
              </h2>
              <p class="text-sm text-gray-500 dark:text-dark-400">
                {{ activeFilterDescription }}
              </p>
            </div>

            <div
              v-if="selectedPlatformLabel"
              class="inline-flex items-center gap-2 rounded-full border border-primary-200 bg-primary-50 px-3 py-1.5 text-sm font-medium text-primary-700 dark:border-primary-500/30 dark:bg-primary-500/10 dark:text-primary-200"
            >
              <span>{{ selectedPlatformLabel }}</span>
              <button
                type="button"
                class="rounded-full p-0.5 transition hover:bg-primary-100 dark:hover:bg-primary-500/20"
                @click="togglePlatformFilter(null)"
              >
                <Icon name="x" size="xs" />
              </button>
            </div>
          </div>

          <div
            v-if="filteredGroups.length === 0"
            class="rounded-[28px] border border-dashed border-gray-300 bg-white/70 px-6 py-12 text-center dark:border-dark-600 dark:bg-dark-900/70"
          >
            <div
              class="mx-auto flex h-14 w-14 items-center justify-center rounded-2xl bg-gray-100 dark:bg-dark-800"
            >
              <Icon name="cube" size="lg" class="text-gray-400 dark:text-dark-400" />
            </div>
            <h3 class="mt-4 text-lg font-semibold text-gray-900 dark:text-white">
              {{ t('modelPlaza.emptyTitle') }}
            </h3>
            <p class="mt-2 text-sm text-gray-500 dark:text-dark-400">
              {{ t('modelPlaza.emptyDescription') }}
            </p>
          </div>

          <div v-else class="grid gap-4 md:grid-cols-2 2xl:grid-cols-3">
            <button
              v-for="group in filteredGroups"
              :key="group.id"
              type="button"
              class="group overflow-hidden rounded-[30px] border border-white/70 bg-white/90 p-5 text-left shadow-[0_18px_50px_-35px_rgba(15,23,42,0.55)] transition-all duration-200 hover:-translate-y-1 hover:shadow-[0_26px_60px_-35px_rgba(15,23,42,0.65)] dark:border-white/10 dark:bg-dark-900/90"
              @click="openGroup(group)"
            >
              <div class="flex items-start justify-between gap-4">
                <div class="min-w-0">
                  <h3
                    class="truncate text-xl font-semibold tracking-tight text-gray-900 dark:text-white"
                  >
                    {{ group.name }}
                  </h3>
                </div>

                <div
                  class="rounded-2xl bg-gray-100/90 p-2 text-gray-400 transition group-hover:bg-primary-50 group-hover:text-primary-600 dark:bg-dark-800 dark:text-dark-400 dark:group-hover:bg-primary-500/10 dark:group-hover:text-primary-300"
                >
                  <Icon name="arrowRight" size="sm" />
                </div>
              </div>

              <div class="mt-6 grid grid-cols-2 gap-3">
                <div class="rounded-2xl bg-emerald-50 px-4 py-3 dark:bg-emerald-500/10">
                  <p
                    class="text-xs font-medium uppercase tracking-[0.18em] text-emerald-600 dark:text-emerald-300"
                  >
                    {{ t('modelPlaza.rateLabel') }}
                  </p>
                  <p class="mt-2 text-xl font-semibold text-emerald-700 dark:text-emerald-200">
                    {{ formatMultiplier(group.rate_multiplier) }}
                  </p>
                </div>
                <div class="rounded-2xl bg-sky-50 px-4 py-3 dark:bg-sky-500/10">
                  <p
                    class="text-xs font-medium uppercase tracking-[0.18em] text-sky-600 dark:text-sky-300"
                  >
                    {{ t('modelPlaza.modelsLabel') }}
                  </p>
                  <p class="mt-2 text-xl font-semibold text-sky-700 dark:text-sky-200">
                    {{ group.model_count }}
                  </p>
                </div>
              </div>
            </button>
          </div>
        </section>
      </template>

      <div
        v-else
        class="rounded-[28px] border border-dashed border-gray-300 bg-white/80 px-6 py-12 text-center dark:border-dark-600 dark:bg-dark-900/70"
      >
        <div
          class="mx-auto flex h-14 w-14 items-center justify-center rounded-2xl bg-red-50 dark:bg-red-500/10"
        >
          <Icon name="exclamationTriangle" size="lg" class="text-red-500 dark:text-red-300" />
        </div>
        <h3 class="mt-4 text-lg font-semibold text-gray-900 dark:text-white">
          {{ t('modelPlaza.loadFailedTitle') }}
        </h3>
        <p class="mt-2 text-sm text-gray-500 dark:text-dark-400">
          {{ t('modelPlaza.loadFailedDescription') }}
        </p>
        <button type="button" class="btn btn-primary mt-5" @click="loadModelPlaza">
          {{ t('common.retry') }}
        </button>
      </div>
    </div>

    <Teleport to="body">
      <Transition name="fade">
        <div
          v-if="selectedGroup"
          class="fixed inset-0 z-50 bg-slate-950/35 backdrop-blur-[2px]"
          @click="closeDrawer"
        >
          <Transition name="slide-panel">
            <aside
              v-if="selectedGroup"
              class="absolute inset-y-0 right-0 flex w-full max-w-2xl flex-col overflow-hidden border-l border-white/60 bg-white/95 shadow-[0_30px_80px_-30px_rgba(15,23,42,0.6)] backdrop-blur dark:border-white/10 dark:bg-dark-950/95"
              @click.stop
            >
              <div class="border-b border-gray-100 px-5 py-4 dark:border-dark-700">
                <div class="flex items-start justify-between gap-4">
                  <div class="min-w-0">
                    <div class="flex flex-wrap items-center gap-2">
                      <span
                        :class="[
                          'inline-flex h-8 w-8 items-center justify-center rounded-2xl',
                          platformBadgeClass(selectedGroup.platform)
                        ]"
                      >
                        <PlatformIcon
                          :platform="normalizePlatform(selectedGroup.platform)"
                          size="sm"
                        />
                      </span>
                      <span
                        class="rounded-full bg-gray-100 px-2.5 py-1 text-xs font-medium text-gray-600 dark:bg-dark-800 dark:text-dark-200"
                      >
                        {{ resolvePlatformLabel(selectedGroup.platform) }}
                      </span>
                      <span
                        class="rounded-full bg-emerald-50 px-2.5 py-1 text-xs font-medium text-emerald-700 dark:bg-emerald-500/10 dark:text-emerald-200"
                      >
                        {{ t('modelPlaza.rateShort', { rate: formatMultiplier(selectedGroup.rate_multiplier) }) }}
                      </span>
                    </div>
                    <h2
                      class="mt-3 truncate text-2xl font-semibold tracking-tight text-gray-900 dark:text-white"
                    >
                      {{ selectedGroup.name }}
                    </h2>
                    <p class="mt-2 text-sm text-gray-500 dark:text-dark-400">
                      {{ t('modelPlaza.groupModelCount', { count: selectedGroup.model_count }) }}
                    </p>
                  </div>

                  <button
                    type="button"
                    class="rounded-2xl border border-gray-200 p-2 text-gray-500 transition hover:border-gray-300 hover:bg-gray-100 hover:text-gray-700 dark:border-dark-700 dark:text-dark-300 dark:hover:border-dark-500 dark:hover:bg-dark-800 dark:hover:text-white"
                    @click="closeDrawer"
                  >
                    <Icon name="x" size="md" />
                  </button>
                </div>
              </div>

              <div class="flex-1 overflow-y-auto p-5">
                <div
                  v-if="selectedGroup.models.length === 0"
                  class="rounded-[28px] border border-dashed border-gray-300 px-5 py-12 text-center dark:border-dark-600"
                >
                  <div
                    class="mx-auto flex h-12 w-12 items-center justify-center rounded-2xl bg-gray-100 dark:bg-dark-800"
                  >
                    <Icon name="sparkles" size="md" class="text-gray-400 dark:text-dark-400" />
                  </div>
                  <h3 class="mt-4 text-lg font-semibold text-gray-900 dark:text-white">
                    {{ t('modelPlaza.noModelsTitle') }}
                  </h3>
                  <p class="mt-2 text-sm text-gray-500 dark:text-dark-400">
                    {{ t('modelPlaza.noModelsDescription') }}
                  </p>
                </div>

                <div v-else class="grid gap-3 md:grid-cols-2">
                  <article
                    v-for="model in selectedGroup.models"
                    :key="model.name"
                    class="rounded-[26px] border border-gray-200/80 bg-gradient-to-br from-white via-white to-gray-50 p-4 shadow-[0_18px_42px_-34px_rgba(15,23,42,0.65)] dark:border-dark-700 dark:from-dark-900 dark:via-dark-900 dark:to-dark-800"
                  >
                    <div class="flex items-start justify-between gap-3">
                      <h3
                        class="min-w-0 text-sm font-semibold leading-6 text-gray-900 dark:text-white"
                      >
                        {{ model.name }}
                      </h3>
                      <span
                        class="shrink-0 rounded-full bg-gray-100 px-2.5 py-1 text-[11px] font-medium uppercase tracking-[0.16em] text-gray-600 dark:bg-dark-700 dark:text-dark-200"
                      >
                        {{ billingModeLabel(model.billing_mode) }}
                      </span>
                    </div>

                    <div class="mt-4 space-y-2.5">
                      <div
                        class="flex items-center justify-between rounded-2xl bg-gray-50 px-3 py-2 dark:bg-dark-800/80"
                      >
                        <span class="text-xs font-medium text-gray-500 dark:text-dark-400">
                          {{ t('modelPlaza.inputPrice') }}
                        </span>
                        <span class="text-sm font-semibold text-gray-900 dark:text-white">
                          {{ formatPrice(model.input_price_1m) }}
                        </span>
                      </div>
                      <div
                        class="flex items-center justify-between rounded-2xl bg-gray-50 px-3 py-2 dark:bg-dark-800/80"
                      >
                        <span class="text-xs font-medium text-gray-500 dark:text-dark-400">
                          {{ t('modelPlaza.outputPrice') }}
                        </span>
                        <span class="text-sm font-semibold text-gray-900 dark:text-white">
                          {{ formatPrice(model.output_price_1m) }}
                        </span>
                      </div>
                      <div
                        class="flex items-center justify-between rounded-2xl bg-gray-50 px-3 py-2 dark:bg-dark-800/80"
                      >
                        <span class="text-xs font-medium text-gray-500 dark:text-dark-400">
                          {{ t('modelPlaza.cacheWritePrice') }}
                        </span>
                        <span class="text-sm font-semibold text-gray-900 dark:text-white">
                          {{ formatPrice(model.cache_write_price_1m) }}
                        </span>
                      </div>
                      <div
                        class="flex items-center justify-between rounded-2xl bg-gray-50 px-3 py-2 dark:bg-dark-800/80"
                      >
                        <span class="text-xs font-medium text-gray-500 dark:text-dark-400">
                          {{ t('modelPlaza.cacheReadPrice') }}
                        </span>
                        <span class="text-sm font-semibold text-gray-900 dark:text-white">
                          {{ formatPrice(model.cache_read_price_1m) }}
                        </span>
                      </div>
                    </div>
                  </article>
                </div>
              </div>
            </aside>
          </Transition>
        </div>
      </Transition>
    </Teleport>
  </AppLayout>
</template>

<script setup lang="ts">
import { computed, onMounted, ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { modelPlazaAPI, type ModelPlazaGroup, type ModelPlazaResponse } from '@/api/modelPlaza'
import { useAppStore } from '@/stores/app'
import type { GroupPlatform } from '@/types'
import PlatformIcon from '@/components/common/PlatformIcon.vue'
import AppLayout from '@/components/layout/AppLayout.vue'
import Icon from '@/components/icons/Icon.vue'

interface PlatformCard {
  platform: string
  label: string
  groupCount: number
  modelCount: number
}

const { t } = useI18n()
const appStore = useAppStore()

const loading = ref(false)
const data = ref<ModelPlazaResponse | null>(null)
const selectedPlatform = ref<string | null>(null)
const selectedGroupId = ref<number | null>(null)

const priceFormatter = new Intl.NumberFormat(undefined, {
  minimumFractionDigits: 0,
  maximumFractionDigits: 4
})

const summaryCards = computed(() => {
  const summary = data.value?.summary
  return [
    {
      key: 'platforms',
      label: t('modelPlaza.summaryPlatforms'),
      value: summary?.platform_count ?? 0,
      icon: 'grid' as const,
      accent: 'bg-gradient-to-r from-sky-400 to-cyan-400',
      iconBg: 'bg-sky-50 dark:bg-sky-500/10',
      iconColor: 'text-sky-600 dark:text-sky-300'
    },
    {
      key: 'groups',
      label: t('modelPlaza.summaryGroups'),
      value: summary?.group_count ?? 0,
      icon: 'database' as const,
      accent: 'bg-gradient-to-r from-emerald-400 to-teal-400',
      iconBg: 'bg-emerald-50 dark:bg-emerald-500/10',
      iconColor: 'text-emerald-600 dark:text-emerald-300'
    },
    {
      key: 'models',
      label: t('modelPlaza.summaryModels'),
      value: summary?.model_count ?? 0,
      icon: 'cube' as const,
      accent: 'bg-gradient-to-r from-amber-400 to-orange-400',
      iconBg: 'bg-amber-50 dark:bg-amber-500/10',
      iconColor: 'text-amber-600 dark:text-amber-300'
    }
  ]
})

const platformCards = computed<PlatformCard[]>(() =>
  (data.value?.platforms ?? []).map((platform) => ({
    platform: platform.platform,
    label: platform.label,
    groupCount: platform.group_count,
    modelCount: platform.groups.reduce((total, group) => total + group.model_count, 0)
  }))
)

const allGroups = computed<ModelPlazaGroup[]>(() =>
  (data.value?.platforms ?? []).flatMap((platform) => platform.groups)
)

const filteredGroups = computed<ModelPlazaGroup[]>(() => {
  if (!selectedPlatform.value) return allGroups.value
  return allGroups.value.filter((group) => group.platform === selectedPlatform.value)
})

const selectedGroup = computed<ModelPlazaGroup | null>(
  () => filteredGroups.value.find((group) => group.id === selectedGroupId.value) ?? null
)

const selectedPlatformLabel = computed(() => {
  if (!selectedPlatform.value) return ''
  return (
    platformCards.value.find((platform) => platform.platform === selectedPlatform.value)?.label ??
    selectedPlatform.value
  )
})

const activeFilterDescription = computed(() => {
  if (!selectedPlatformLabel.value) {
    return t('modelPlaza.groupsHint', { count: filteredGroups.value.length })
  }
  return t('modelPlaza.filteredGroupsHint', {
    platform: selectedPlatformLabel.value,
    count: filteredGroups.value.length
  })
})

watch(selectedPlatform, () => {
  if (!selectedGroup.value) {
    selectedGroupId.value = null
  }
})

async function loadModelPlaza() {
  loading.value = true
  try {
    data.value = await modelPlazaAPI.getModelPlaza()
  } catch (error) {
    console.error('Failed to load model plaza:', error)
    data.value = null
    appStore.showError(t('modelPlaza.failedToLoad'))
  } finally {
    loading.value = false
  }
}

function togglePlatformFilter(platform: string | null) {
  selectedPlatform.value = selectedPlatform.value === platform ? null : platform
}

function openGroup(group: ModelPlazaGroup) {
  selectedGroupId.value = group.id
}

function closeDrawer() {
  selectedGroupId.value = null
}

function resolvePlatformLabel(platform: string): string {
  return (
    platformCards.value.find((item) => item.platform === platform)?.label ??
    data.value?.platforms.find((item) => item.platform === platform)?.label ??
    platform
  )
}

function normalizePlatform(platform: string): GroupPlatform | undefined {
  switch (platform) {
    case 'openai':
    case 'anthropic':
    case 'gemini':
    case 'antigravity':
      return platform
    default:
      return undefined
  }
}

function platformBadgeClass(platform: string): string {
  switch (platform) {
    case 'openai':
      return 'bg-emerald-50 text-emerald-700 ring-1 ring-emerald-200 dark:bg-emerald-500/10 dark:text-emerald-200 dark:ring-emerald-500/20'
    case 'anthropic':
      return 'bg-amber-50 text-amber-700 ring-1 ring-amber-200 dark:bg-amber-500/10 dark:text-amber-200 dark:ring-amber-500/20'
    case 'gemini':
      return 'bg-sky-50 text-sky-700 ring-1 ring-sky-200 dark:bg-sky-500/10 dark:text-sky-200 dark:ring-sky-500/20'
    case 'antigravity':
      return 'bg-rose-50 text-rose-700 ring-1 ring-rose-200 dark:bg-rose-500/10 dark:text-rose-200 dark:ring-rose-500/20'
    default:
      return 'bg-gray-100 text-gray-700 ring-1 ring-gray-200 dark:bg-dark-700 dark:text-dark-100 dark:ring-dark-600'
  }
}

function formatMultiplier(rate: number): string {
  const fixed = Number.isInteger(rate) ? rate.toFixed(0) : rate.toFixed(rate >= 10 ? 1 : 2)
  return `${fixed.replace(/\.0+$/, '').replace(/(\.\d*[1-9])0+$/, '$1')}x`
}

function formatPrice(price: number | null): string {
  if (price == null) return '-'
  return `$${priceFormatter.format(price)}/1M`
}

function billingModeLabel(mode: string): string {
  switch (mode) {
    case 'token':
      return t('modelPlaza.billingModeToken')
    case 'per_request':
      return t('modelPlaza.billingModeRequest')
    case 'image':
      return t('modelPlaza.billingModeImage')
    default:
      return mode
  }
}

onMounted(() => {
  loadModelPlaza()
})
</script>

<style scoped>
.slide-panel-enter-active,
.slide-panel-leave-active {
  transition:
    transform 0.24s ease,
    opacity 0.24s ease;
}

.slide-panel-enter-from,
.slide-panel-leave-to {
  opacity: 0;
  transform: translateX(24px);
}
</style>
