#!/bin/sh
set -eu

SCRIPT_DIR="$(CDPATH= cd "$(dirname "$0")" && pwd -P)"
REPO_ROOT="$(cd "${SCRIPT_DIR}/../.." && pwd)"
APP_DIR="${REPO_ROOT}/third_party/gpt_image_playground"
FRONTEND_DIR="${REPO_ROOT}/frontend"
OUTPUT_DIR="${FRONTEND_DIR}/public/image-playground-app"

if [ ! -f "${APP_DIR}/package.json" ]; then
  echo "gpt_image_playground source not found at ${APP_DIR}" >&2
  exit 1
fi

if ! command -v npm >/dev/null 2>&1; then
  echo "npm is required to build gpt_image_playground" >&2
  exit 1
fi

cd "${APP_DIR}"

if [ ! -d node_modules ]; then
  npm ci
fi

PATCH_BACKUP_DIR="$(mktemp -d)"
restore_sub2api_patch() {
  if [ -d "${PATCH_BACKUP_DIR}/files" ]; then
    (
      cd "${PATCH_BACKUP_DIR}/files"
      find . -type f | while IFS= read -r file; do
        mkdir -p "${APP_DIR}/$(dirname "${file}")"
        cp "${file}" "${APP_DIR}/${file}"
      done
    )
  fi
  rm -rf "${PATCH_BACKUP_DIR}"
}
trap restore_sub2api_patch EXIT INT TERM

SUB2API_PATCH_BACKUP_DIR="${PATCH_BACKUP_DIR}" node <<'NODE'
const fs = require('node:fs')
const path = require('node:path')

const appDir = process.cwd()
const backupDir = process.env.SUB2API_PATCH_BACKUP_DIR

if (!backupDir) {
  console.error('SUB2API_PATCH_BACKUP_DIR is required')
  process.exit(1)
}

function backupFile(relativePath) {
  const source = path.join(appDir, relativePath)
  const target = path.join(backupDir, 'files', relativePath)

  if (fs.existsSync(target)) return
  fs.mkdirSync(path.dirname(target), { recursive: true })
  fs.copyFileSync(source, target)
}

function patchFile(relativePath, replacements) {
  const file = path.join(appDir, relativePath)
  let source = fs.readFileSync(file, 'utf8')
  const original = source

  for (const [needle, replacement] of replacements) {
    if (!source.includes(needle)) {
      if (replacement !== '' && source.includes(replacement)) continue
      console.error(`Sub2API image playground patch failed: ${relativePath} did not contain expected snippet`)
      process.exit(1)
    }
    source = source.replace(needle, replacement)
  }

  if (source !== original) {
    backupFile(relativePath)
    fs.writeFileSync(file, source)
  }
}

function patchFilePattern(relativePath, pattern, replacement) {
  const file = path.join(appDir, relativePath)
  let source = fs.readFileSync(file, 'utf8')
  const original = source

  if (!pattern.test(source)) {
    if (replacement !== '' && source.includes(replacement)) return
    console.error(`Sub2API image playground patch failed: ${relativePath} did not match expected pattern`)
    process.exit(1)
  }

  pattern.lastIndex = 0
  source = source.replace(pattern, replacement)

  if (source !== original) {
    backupFile(relativePath)
    fs.writeFileSync(file, source)
  }
}

patchFile('src/App.tsx', [
  [
`      const provider: ApiProvider | null = providerParam === 'fal'
        ? 'fal'
        : ['openai', 'openai-compatible'].includes(providerParam)
          ? 'openai'
          : null`,
`      const provider: ApiProvider | null = ['openai', 'openai-compatible'].includes(providerParam)
        ? 'openai'
        : null`,
  ],
])

patchFile('src/lib/api.ts', [
  [
`import { callFalAiImageApi } from './falAiImageApi'
`,
``,
  ],
  [
`  if (profile.provider === 'fal') return callFalAiImageApi(opts, profile)

  return callOpenAICompatibleImageApi(opts, profile)`,
`  return callOpenAICompatibleImageApi(opts, profile)`,
  ],
])

patchFile('src/lib/apiProfiles.ts', [
  [
`export const DEFAULT_FAL_BASE_URL = 'https://fal.run'
export const DEFAULT_FAL_MODEL = 'openai/gpt-image-2'
`,
``,
  ],
  [
`export function createDefaultFalProfile(overrides: Partial<ApiProfile> = {}): ApiProfile {
  return {
    id: \`fal-\${Date.now().toString(36)}-\${Math.random().toString(36).slice(2, 6)}\`,
    name: '新配置',
    provider: 'fal',
    baseUrl: DEFAULT_FAL_BASE_URL,
    apiKey: '',
    model: DEFAULT_FAL_MODEL,
    timeout: DEFAULT_API_TIMEOUT,
    apiMode: 'images',
    codexCli: false,
    apiProxy: false,
    ...overrides,
  }
}

`,
``,
  ],
  [
`export function switchApiProfileProvider(profile: ApiProfile, provider: ApiProvider): ApiProfile {
  if (provider === 'fal') {
    return {
      ...profile,
      provider,
      baseUrl: DEFAULT_FAL_BASE_URL,
      model: DEFAULT_FAL_MODEL,
      apiMode: 'images',
      codexCli: false,
      apiProxy: false,
    }
  }

  return {
    ...profile,
    provider,
    baseUrl: DEFAULT_BASE_URL,
    model: DEFAULT_IMAGES_MODEL,
  }
}`,
`export function switchApiProfileProvider(profile: ApiProfile, _provider: ApiProvider): ApiProfile {
  return {
    ...profile,
    provider: 'openai',
    baseUrl: DEFAULT_BASE_URL,
    model: DEFAULT_IMAGES_MODEL,
  }
}`,
  ],
  [
`  const provider: ApiProvider = record.provider === 'fal' ? 'fal' : 'openai'
  const defaults = provider === 'fal' ? createDefaultFalProfile(fallback) : createDefaultOpenAIProfile(fallback)
  const apiMode: ApiMode = record.apiMode === 'responses' ? 'responses' : 'images'`,
`  const provider: ApiProvider = 'openai'
  const isOpenAIProvider = !('provider' in record) || record.provider === 'openai'
  const defaults = createDefaultOpenAIProfile(fallback)
  const apiMode: ApiMode = isOpenAIProvider && record.apiMode === 'responses' ? 'responses' : 'images'`,
  ],
  [
`    baseUrl: typeof record.baseUrl === 'string' ? record.baseUrl : defaults.baseUrl,
    apiKey: typeof record.apiKey === 'string' ? record.apiKey : defaults.apiKey,
    model: typeof record.model === 'string' && record.model.trim() ? record.model : defaults.model,`,
`    baseUrl: isOpenAIProvider && typeof record.baseUrl === 'string' ? record.baseUrl : defaults.baseUrl,
    apiKey: isOpenAIProvider && typeof record.apiKey === 'string' ? record.apiKey : defaults.apiKey,
    model: isOpenAIProvider && typeof record.model === 'string' && record.model.trim() ? record.model : defaults.model,`,
  ],
])

patchFile('src/lib/paramCompatibility.ts', [
  [
`import { getActiveApiProfile } from './apiProfiles'
`,
`import { getActiveApiProfile } from './apiProfiles'
`,
  ],
  [
`export const DEFAULT_FAL_IMAGE_SIZE = '1360x1024'
export const MAX_FAL_OUTPUT_IMAGES = 4
export const MAX_OPENAI_OUTPUT_IMAGES = 10`,
`export const MAX_OPENAI_OUTPUT_IMAGES = 10`,
  ],
  [
`export function getOutputImageLimitForSettings(settings: AppSettings) {
  return getActiveApiProfile(settings).provider === 'fal' ? MAX_FAL_OUTPUT_IMAGES : MAX_OPENAI_OUTPUT_IMAGES
}`,
`export function getOutputImageLimitForSettings(_settings: AppSettings) {
  return MAX_OPENAI_OUTPUT_IMAGES
}`,
  ],
  [
`  if (activeProfile.provider === 'openai' && activeProfile.codexCli) {
    nextParams.quality = DEFAULT_PARAMS.quality
  }

  if (activeProfile.provider === 'fal') {
    if (nextParams.size === 'auto') nextParams.size = DEFAULT_FAL_IMAGE_SIZE
    if (nextParams.quality === 'auto') nextParams.quality = 'high'
    nextParams.moderation = DEFAULT_PARAMS.moderation
    nextParams.output_compression = DEFAULT_PARAMS.output_compression
  }`,
`  if (activeProfile.codexCli) {
    nextParams.quality = DEFAULT_PARAMS.quality
  }`,
  ],
])

patchFilePattern('src/store.ts', /import \{ getFalErrorMessage, getFalQueuedImageResult, getFalQueueStatus \} from '\.\/lib\/falAiImageApi'\n/, '')
patchFilePattern('src/store.ts', /const FAL_RECOVERY_POLL_MS = 10_000\nconst falRecoveryTimers = new Map<string, ReturnType<typeof setTimeout>>\(\)\n/, '')
patchFilePattern('src/store.ts', /\nfunction getFalRecoveryProfile[\s\S]*?\n}\n\n\/\*\* 初始化：从 IndexedDB 加载任务和图片缓存，清理孤立图片 \*\//, `
/** 初始化：从 IndexedDB 加载任务和图片缓存，清理孤立图片 */`)
patchFile('src/store.ts', [
  [
`  for (const task of tasks) {
    if (
      task.apiProvider === 'fal' &&
      task.falRequestId &&
      task.falEndpoint &&
      (task.status === 'running' || task.falRecoverable)
    ) {
      scheduleFalRecovery(task.id, 0)
    }
  }

`,
``,
  ],
  [
`  const taskProvider = task.apiProvider ?? activeProfile.provider
  let falRequestInfo: { requestId: string; endpoint: string } | null = task.falRequestId && task.falEndpoint
    ? { requestId: task.falRequestId, endpoint: task.falEndpoint }
    : null

  if (taskProvider === 'openai') {
    scheduleOpenAIWatchdog(taskId, activeProfile.timeout)
  }`,
`  const taskProvider = 'openai'

  scheduleOpenAIWatchdog(taskId, activeProfile.timeout)`,
  ],
  [
`      onFalRequestEnqueued: (request) => {
        falRequestInfo = request
        updateTaskInStore(taskId, {
          falRequestId: request.requestId,
          falEndpoint: request.endpoint,
          falRecoverable: false,
        })
      },
`,
``,
  ],
  [
`    const shouldStoreApiResponseMetadata = taskProvider !== 'fal'`,
`    const shouldStoreApiResponseMetadata = true`,
  ],
])
patchFilePattern('src/store.ts', /    const latestFalRequestInfo = falRequestInfo \?\? \(latestTask\.falRequestId && latestTask\.falEndpoint[\s\S]*?      useStore\.getState\(\)\.setDetailTaskId\(taskId\)\n    \}/, `    updateTaskInStore(taskId, {
      status: 'error',
      error: err instanceof Error ? err.message : String(err),
      falRecoverable: false,
      finishedAt: Date.now(),
      elapsed: Date.now() - task.createdAt,
    })
    useStore.getState().setDetailTaskId(taskId)`)

patchFile('src/components/InputBar.tsx', [
  [
`import { DEFAULT_FAL_IMAGE_SIZE, getChangedParams, getOutputImageLimitForSettings, normalizeParamsForSettings } from '../lib/paramCompatibility'`,
`import { getChangedParams, getOutputImageLimitForSettings, normalizeParamsForSettings } from '../lib/paramCompatibility'`,
  ],
  [
`  const activeProfile = getActiveApiProfile(settings)
  const activeProvider = activeProfile.provider
  const isFalProvider = activeProvider === 'fal'
  const moderationDisabled = settings.apiMode === 'responses' || isFalProvider
  const compressionDisabled = params.output_format === 'png' || isFalProvider
  const outputImageLimit = getOutputImageLimitForSettings(settings)
  const nLimitHintText = isFalProvider
    ? \`fal.ai 最大请求数量为 \${outputImageLimit}\`
    : \`OpenAI 最大请求数量为 \${outputImageLimit}\`
  const displaySize = isFalProvider && params.size === 'auto'
    ? DEFAULT_FAL_IMAGE_SIZE
    : normalizeImageSize(params.size) || DEFAULT_PARAMS.size
  const qualityOptions = isFalProvider
    ? [
        { label: 'low', value: 'low' },
        { label: 'medium', value: 'medium' },
        { label: 'high', value: 'high' },
      ]
    : [
        { label: 'auto', value: 'auto' },
        { label: 'low', value: 'low' },
        { label: 'medium', value: 'medium' },
        { label: 'high', value: 'high' },
      ]`,
`  const activeProfile = getActiveApiProfile(settings)
  const moderationDisabled = activeProfile.apiMode === 'responses'
  const compressionDisabled = params.output_format === 'png'
  const outputImageLimit = getOutputImageLimitForSettings(settings)
  const nLimitHintText = \`OpenAI 最大请求数量为 \${outputImageLimit}\`
  const displaySize = normalizeImageSize(params.size) || DEFAULT_PARAMS.size
  const qualityOptions = [
    { label: 'auto', value: 'auto' },
    { label: 'low', value: 'low' },
    { label: 'medium', value: 'medium' },
    { label: 'high', value: 'high' },
  ]`,
  ],
  [
`  const showQualityHint = () => {
    if (settings.codexCli || isFalProvider) setQualityHintVisible(true)
  }

  const showSizeHint = () => {
    if (isFalProvider) setSizeHintVisible(true)
  }`,
`  const showQualityHint = () => {
    if (settings.codexCli) setQualityHintVisible(true)
  }

  const showSizeHint = () => {
    setSizeHintVisible(false)
  }`,
  ],
  [
`  const startSizeHintTouch = () => {
    if (!isFalProvider) return
    sizeHintTimerRef.current = window.setTimeout(() => {
      setSizeHintVisible(true)
      sizeHintTimerRef.current = null
    }, 450)
  }`,
`  const startSizeHintTouch = () => {
    setSizeHintVisible(false)
  }`,
  ],
  [
`  const startQualityHintTouch = () => {
    if (!settings.codexCli && !isFalProvider) return
    qualityHintTimerRef.current = window.setTimeout(() => {
      setQualityHintVisible(true)
      qualityHintTimerRef.current = null
    }, 450)
  }`,
`  const startQualityHintTouch = () => {
    if (!settings.codexCli) return
    qualityHintTimerRef.current = window.setTimeout(() => {
      setQualityHintVisible(true)
      qualityHintTimerRef.current = null
    }, 450)
  }`,
  ],
  [
`        <ButtonTooltip
          visible={isFalProvider && sizeHintVisible}
          text={<>fal.ai 不支持 <code className="rounded bg-white/10 px-1 py-0.5 font-mono">auto</code> 参数</>}
        />`,
`        <ButtonTooltip
          visible={false && sizeHintVisible}
          text=""
        />`,
  ],
  [
`          value={settings.codexCli ? 'auto' : isFalProvider && params.quality === 'auto' ? 'high' : params.quality}`,
`          value={settings.codexCli ? 'auto' : params.quality}`,
  ],
  [
`          visible={(settings.codexCli || isFalProvider) && qualityHintVisible}
          text={isFalProvider ? <>fal.ai 不支持 <code className="rounded bg-white/10 px-1 py-0.5 font-mono">auto</code> 参数</> : 'Codex CLI 不支持质量参数'}`,
`          visible={settings.codexCli && qualityHintVisible}
          text="Codex CLI 不支持质量参数"`,
  ],
  [
`          text={isFalProvider ? 'fal.ai 不支持压缩率参数' : '仅 JPEG 和 WebP 支持压缩率'}`,
`          text="仅 JPEG 和 WebP 支持压缩率"`,
  ],
  [
`          text={isFalProvider ? 'fal.ai 不支持审核参数' : 'Responses API 不支持审核参数'}`,
`          text="Responses API 不支持审核参数"`,
  ],
  [
`          currentSize={isFalProvider && params.size === 'auto' ? DEFAULT_FAL_IMAGE_SIZE : params.size}
          onSelect={(size) => setParams({ size })}
          onClose={() => setShowSizePicker(false)}
          allowAuto={!isFalProvider}`,
`          currentSize={params.size}
          onSelect={(size) => setParams({ size })}
          onClose={() => setShowSizePicker(false)}
          allowAuto`,
  ],
])

patchFile('src/components/SettingsModal.tsx', [
  [
`  DEFAULT_FAL_BASE_URL,
  DEFAULT_FAL_MODEL,
`,
``,
  ],
  [
`function providerLabel(provider: string) {
  return provider === 'fal' ? 'fal.ai' : 'OpenAI'
}`,
`function providerLabel(_provider: string) {
  return 'OpenAI'
}`,
  ],
  [
`      const normalizedBaseUrl = profile.provider === 'fal'
        ? profile.baseUrl.trim().replace(/\\/+$/, '') || DEFAULT_FAL_BASE_URL
        : normalizeBaseUrl(profile.baseUrl.trim() || DEFAULT_SETTINGS.baseUrl)
      const defaultModel = profile.provider === 'fal' ? DEFAULT_FAL_MODEL : getDefaultModelForMode(profile.apiMode)`,
`      const normalizedBaseUrl = normalizeBaseUrl(profile.baseUrl.trim() || DEFAULT_SETTINGS.baseUrl)
      const defaultModel = getDefaultModelForMode(profile.apiMode)`,
  ],
  [
`                  options={[{ label: 'OpenAI 兼容接口', value: 'openai' }, { label: 'fal.ai', value: 'fal' }]}`,
`                  options={[{ label: 'OpenAI 兼容接口', value: 'openai' }]}`,
  ],
  [
`                    placeholder={activeProfile.provider === 'fal' ? 'FAL_KEY' : 'sk-...'}`,
`                    placeholder="sk-..."`,
  ],
  [
`                  placeholder={activeProfile.provider === 'fal' ? DEFAULT_FAL_MODEL : getDefaultModelForMode(activeProfile.apiMode ?? DEFAULT_SETTINGS.apiMode)}`,
`                  placeholder={getDefaultModelForMode(activeProfile.apiMode ?? DEFAULT_SETTINGS.apiMode)}`,
  ],
  [
`                  {activeProfile.provider === 'fal' ? (
                    <>当前适配 <code className="rounded bg-gray-100 px-1 py-0.5 dark:bg-white/[0.06]">{DEFAULT_FAL_MODEL}</code>。</>
                  ) : (activeProfile.apiMode ?? DEFAULT_SETTINGS.apiMode) === 'responses' ? (
                    <>Responses API 需要使用支持 <code className="rounded bg-gray-100 px-1 py-0.5 dark:bg-white/[0.06]">image_generation</code> 工具的文本模型，例如 <code className="rounded bg-gray-100 px-1 py-0.5 dark:bg-white/[0.06]">{DEFAULT_RESPONSES_MODEL}</code>。</>
                  ) : (
                    <>Images API 需要使用 GPT Image 模型，例如 <code className="rounded bg-gray-100 px-1 py-0.5 dark:bg-white/[0.06]">{DEFAULT_IMAGES_MODEL}</code>。</>
                  )}`,
`                  {(activeProfile.apiMode ?? DEFAULT_SETTINGS.apiMode) === 'responses' ? (
                    <>Responses API 需要使用支持 <code className="rounded bg-gray-100 px-1 py-0.5 dark:bg-white/[0.06]">image_generation</code> 工具的文本模型，例如 <code className="rounded bg-gray-100 px-1 py-0.5 dark:bg-white/[0.06]">{DEFAULT_RESPONSES_MODEL}</code>。</>
                  ) : (
                    <>Images API 需要使用 GPT Image 模型，例如 <code className="rounded bg-gray-100 px-1 py-0.5 dark:bg-white/[0.06]">{DEFAULT_IMAGES_MODEL}</code>。</>
                  )}`,
  ],
])

patchFile('src/components/DetailModal.tsx', [
  [
`  const taskProviderName = taskProvider === 'fal' ? 'fal.ai' : taskProvider ? 'OpenAI' : '未知'`,
`  const taskProviderName = taskProvider ? 'OpenAI' : '未知'`,
  ],
])

patchFile('src/lib/apiProfiles.test.ts', [
  [
`  DEFAULT_FAL_BASE_URL,
  DEFAULT_FAL_MODEL,
`,
``,
  ],
  [
`      activeProfileId: 'imported-fal',
    })

    expect(merged.profiles.map((profile) => profile.id)).toEqual(['imported-openai', 'imported-fal'])
    expect(merged.activeProfileId).toBe('imported-fal')`,
`      activeProfileId: 'imported-openai',
    })

    expect(merged.profiles.map((profile) => profile.id)).toEqual(['imported-openai'])
    expect(merged.activeProfileId).toBe('imported-openai')`,
  ],
  [
`    expect(merged.profiles).toHaveLength(3)
    expect(merged.activeProfileId).toBe(DEFAULT_OPENAI_PROFILE_ID)
    expect(merged.profiles[0]).toMatchObject({ apiKey: 'current-key', model: 'current-model' })
    expect(merged.profiles[1]).toMatchObject({ name: 'Imported OpenAI', provider: 'openai', apiKey: 'imported-key' })
    expect(merged.profiles[2]).toMatchObject({ name: 'Imported fal', provider: 'fal', apiKey: 'fal-key' })
    expect(new Set(merged.profiles.map((profile) => profile.id)).size).toBe(3)`,
`    expect(merged.profiles).toHaveLength(2)
    expect(merged.activeProfileId).toBe(DEFAULT_OPENAI_PROFILE_ID)
    expect(merged.profiles[0]).toMatchObject({ apiKey: 'current-key', model: 'current-model' })
    expect(merged.profiles[1]).toMatchObject({ name: 'Imported OpenAI', provider: 'openai', apiKey: 'imported-key' })
    expect(new Set(merged.profiles.map((profile) => profile.id)).size).toBe(2)`,
  ],
  [
`    expect(merged.profiles).toHaveLength(2)
    expect(merged.profiles[0]).toMatchObject({ apiKey: 'current-key', model: 'current-model' })
    expect(merged.profiles[1]).toMatchObject({ provider: 'fal', apiKey: 'fal-key', model: DEFAULT_FAL_MODEL })`,
`    expect(merged.profiles).toHaveLength(1)
    expect(merged.profiles[0]).toMatchObject({ apiKey: 'current-key', model: 'current-model' })`,
  ],
])

patchFilePattern('src/lib/apiProfiles.test.ts', /\n        \{\n          id: '[^']*fal',\n          name: '[^']*fal',\n          provider: 'fal',\n          baseUrl: DEFAULT_FAL_BASE_URL,\n          apiKey: 'fal-key',\n          model: DEFAULT_FAL_MODEL,\n          timeout: 300,\n          apiMode: 'images',\n          codexCli: false,\n          apiProxy: false,\n        \},/g, '')
patchFilePattern('src/lib/paramCompatibility.test.ts', /import \{ createDefaultFalProfile, DEFAULT_SETTINGS, normalizeSettings \} from '\.\/apiProfiles'\n/, `import { DEFAULT_SETTINGS, normalizeSettings } from './apiProfiles'
`)
patchFilePattern('src/lib/paramCompatibility.test.ts', /\n\n  it\('limits fal\.ai output count to 4'[\s\S]*?\n  }\)/, '')
NODE

npm run build
restore_sub2api_patch
trap - EXIT INT TERM

rm -rf "${OUTPUT_DIR}"
mkdir -p "${OUTPUT_DIR}"
cp -a "${APP_DIR}/dist/." "${OUTPUT_DIR}/"

# The embedded Sub2API version is an online tool, not a standalone PWA.
# Disable service worker/PWA integration as a post-build step so upstream
# sources stay easy to sync.
rm -f "${OUTPUT_DIR}/sw.js" "${OUTPUT_DIR}/manifest.webmanifest"

if [ -f "${OUTPUT_DIR}/index.html" ]; then
  IMAGE_PLAYGROUND_INDEX_HTML="${OUTPUT_DIR}/index.html" node <<'NODE'
const fs = require('node:fs')

const file = process.env.IMAGE_PLAYGROUND_INDEX_HTML
if (!file || !fs.existsSync(file)) {
  process.exit(0)
}

const marker = 'data-sub2api-image-playground-theme'
let html = fs.readFileSync(file, 'utf8')
const originalHtml = html

html = html
  .replace(/[ \t]*<link[^>]+rel=["']manifest["'][^>]*>\n?/gi, '')
  .replace(/[ \t]*<link[^>]+rel=["']apple-touch-icon["'][^>]*>\n?/gi, '')

if (!html.includes(marker)) {
  const bridge = `    <style ${marker}>
      html,
      body,
      #root {
        background-color: var(--sub2api-image-playground-bg, #f9fafb) !important;
      }

      html.dark {
        color-scheme: dark;
      }

      html:not(.dark) {
        color-scheme: light;
      }
    </style>
    <script ${marker} nonce="__CSP_NONCE_VALUE__">
      (() => {
        const MESSAGE_TYPE = 'sub2api:image-playground-theme'
        const DEFAULT_LIGHT_BG = '#f9fafb'
        const DEFAULT_DARK_BG = '#020617'

        const normalizeTheme = (value) => value === 'dark' ? 'dark' : 'light'
        const isCssColor = (value) => {
          if (!value || typeof value !== 'string' || value.length > 80) return false
          if (typeof CSS !== 'undefined' && typeof CSS.supports === 'function') {
            return CSS.supports('color', value)
          }
          const probe = document.createElement('span')
          probe.style.color = value
          return probe.style.color !== ''
        }
        const normalizeBackground = (value, theme) => (
          isCssColor(value) ? value : theme === 'dark' ? DEFAULT_DARK_BG : DEFAULT_LIGHT_BG
        )
        const applyTheme = (payload = {}) => {
          const theme = normalizeTheme(payload.theme)
          const background = normalizeBackground(payload.background, theme)

          document.documentElement.classList.toggle('dark', theme === 'dark')
          document.documentElement.style.setProperty('--sub2api-image-playground-bg', background)

          const themeColor = document.querySelector('meta[name="theme-color"]')
          if (themeColor) themeColor.setAttribute('content', background)
          if (document.body) {
            document.body.style.setProperty('background-color', background, 'important')
          }

          const root = document.getElementById('root')
          if (root) {
            root.style.setProperty('background-color', background, 'important')
          }
        }

        const params = new URLSearchParams(window.location.search)
        let latestTheme = {
          theme: normalizeTheme(params.get('theme')),
          background: params.get('sub2apiBg') || ''
        }

        applyTheme(latestTheme)
        document.addEventListener('DOMContentLoaded', () => applyTheme(latestTheme), { once: true })

        window.addEventListener('message', (event) => {
          if (event.origin !== window.location.origin) return
          const data = event.data
          if (!data || data.type !== MESSAGE_TYPE) return
          latestTheme = {
            theme: data.theme,
            background: data.background
          }
          applyTheme(latestTheme)
        })
      })()
    </script>`

  html = html.replace(/\s*<\/head>/, `\n${bridge}\n  </head>`)
}

if (html !== originalHtml) {
  fs.writeFileSync(file, html)
}
NODE
fi

# Upstream uses media-query dark mode. Convert the generated CSS to class-based
# dark mode so the iframe follows Sub2API's current theme instead of the OS.
IMAGE_PLAYGROUND_OUTPUT_DIR="${OUTPUT_DIR}" node <<'NODE'
const fs = require('node:fs')
const path = require('node:path')

const outputDir = process.env.IMAGE_PLAYGROUND_OUTPUT_DIR
const assetDir = outputDir ? path.join(outputDir, 'assets') : ''

if (!assetDir || !fs.existsSync(assetDir)) {
  process.exit(0)
}

function findMatchingBrace(source, openBraceIndex) {
  let depth = 0
  for (let i = openBraceIndex; i < source.length; i += 1) {
    if (source[i] === '{') depth += 1
    if (source[i] === '}') {
      depth -= 1
      if (depth === 0) return i
    }
  }
  return -1
}

function splitSelectorList(selector) {
  const parts = []
  let current = ''
  let escaped = false
  let bracketDepth = 0
  let parenDepth = 0

  for (const char of selector) {
    if (escaped) {
      current += char
      escaped = false
      continue
    }
    if (char === '\\') {
      current += char
      escaped = true
      continue
    }
    if (char === '[') bracketDepth += 1
    if (char === ']') bracketDepth = Math.max(0, bracketDepth - 1)
    if (char === '(') parenDepth += 1
    if (char === ')') parenDepth = Math.max(0, parenDepth - 1)
    if (char === ',' && bracketDepth === 0 && parenDepth === 0) {
      parts.push(current)
      current = ''
      continue
    }
    current += char
  }

  parts.push(current)
  return parts
}

function prefixSelectorList(selector) {
  return splitSelectorList(selector)
    .map((part) => {
      const trimmed = part.trim()
      return trimmed ? `.dark ${trimmed}` : trimmed
    })
    .join(',')
}

function prefixRules(block) {
  let output = ''
  let index = 0

  while (index < block.length) {
    const openBraceIndex = block.indexOf('{', index)
    if (openBraceIndex === -1) {
      output += block.slice(index)
      break
    }

    const selector = block.slice(index, openBraceIndex).trim()
    const closeBraceIndex = findMatchingBrace(block, openBraceIndex)
    if (closeBraceIndex === -1) {
      output += block.slice(index)
      break
    }

    const body = block.slice(openBraceIndex + 1, closeBraceIndex)
    const nextSelector = selector.startsWith('@') ? selector : prefixSelectorList(selector)
    output += `${nextSelector}{${body}}`
    index = closeBraceIndex + 1
  }

  return output
}

function convertDarkMedia(css) {
  const pattern = /@media\s*\(\s*prefers-color-scheme\s*:\s*dark\s*\)\s*\{/g
  let output = ''
  let cursor = 0
  let match

  while ((match = pattern.exec(css)) !== null) {
    const openBraceIndex = pattern.lastIndex - 1
    const closeBraceIndex = findMatchingBrace(css, openBraceIndex)
    if (closeBraceIndex === -1) {
      break
    }

    output += css.slice(cursor, match.index)
    output += prefixRules(css.slice(openBraceIndex + 1, closeBraceIndex))
    cursor = closeBraceIndex + 1
    pattern.lastIndex = cursor
  }

  output += css.slice(cursor)
  return output
}

for (const entry of fs.readdirSync(assetDir)) {
  if (!entry.endsWith('.css')) continue
  const file = path.join(assetDir, entry)
  const css = fs.readFileSync(file, 'utf8')
  fs.writeFileSync(file, convertDarkMedia(css))
}
NODE

IMAGE_PLAYGROUND_OUTPUT_DIR="${OUTPUT_DIR}" node <<'NODE'
const fs = require('node:fs')
const path = require('node:path')

const outputDir = process.env.IMAGE_PLAYGROUND_OUTPUT_DIR
if (!outputDir) {
  console.error('Sub2API image playground theme bridge was not applied')
  process.exit(1)
}

function* walk(dir) {
  if (!fs.existsSync(dir)) return
  for (const entry of fs.readdirSync(dir, { withFileTypes: true })) {
    const file = path.join(dir, entry.name)
    if (entry.isDirectory()) {
      yield* walk(file)
    } else if (entry.isFile()) {
      yield file
    }
  }
}

let failed = false
const indexFile = path.join(outputDir, 'index.html')
if (fs.existsSync(indexFile)) {
  const html = fs.readFileSync(indexFile, 'utf8')
  if (!html.includes('data-sub2api-image-playground-theme')) failed = true
}

const assetDir = path.join(outputDir, 'assets')
for (const file of walk(assetDir)) {
  if (!file.endsWith('.css')) continue
  const css = fs.readFileSync(file, 'utf8')
  if (/prefers-color-scheme\s*:\s*dark/.test(css)) failed = true
}

if (failed) {
  console.error('Sub2API image playground theme bridge was not applied')
  process.exit(1)
}
NODE

IMAGE_PLAYGROUND_OUTPUT_DIR="${OUTPUT_DIR}" node <<'NODE'
const fs = require('node:fs')
const path = require('node:path')

const outputDir = process.env.IMAGE_PLAYGROUND_OUTPUT_DIR
const assetDir = outputDir ? path.join(outputDir, 'assets') : ''

if (!assetDir || !fs.existsSync(assetDir)) {
  process.exit(0)
}

function* walk(dir) {
  for (const entry of fs.readdirSync(dir, { withFileTypes: true })) {
    const file = path.join(dir, entry.name)
    if (entry.isDirectory()) {
      yield* walk(file)
    } else if (entry.isFile()) {
      yield file
    }
  }
}

let failed = false
for (const file of walk(assetDir)) {
  if (!file.endsWith('.js')) continue
  const source = fs.readFileSync(file, 'utf8')
  const patched = source
    .replace(/if\s*\(\s*["']serviceWorker["']\s*in\s*navigator\s*\)\s*\{/g, 'if (false && "serviceWorker" in navigator) {')
    .replace(/(["']serviceWorker["']\s*in\s*navigator)\s*&&/g, 'false && $1 &&')

  if (patched !== source) {
    fs.writeFileSync(file, patched)
  }

  let match
  const serviceWorkerAnd = /(["']serviceWorker["']\s*in\s*navigator\s*&&)/g
  while ((match = serviceWorkerAnd.exec(patched)) !== null) {
    const prefix = patched.slice(0, match.index)
    if (!/false\s*&&\s*$/.test(prefix)) failed = true
  }

  const serviceWorkerIf = /if\s*\(\s*(["']serviceWorker["']\s*in\s*navigator)/g
  while ((match = serviceWorkerIf.exec(patched)) !== null) {
    const snippet = patched.slice(match.index, match.index + 80)
    if (!/^if\s*\(\s*false\s*&&/.test(snippet)) failed = true
  }
}

if (failed) {
  console.error('service worker registration guard was not patched')
  process.exit(1)
}
NODE

echo "Image Playground built into ${OUTPUT_DIR}"
