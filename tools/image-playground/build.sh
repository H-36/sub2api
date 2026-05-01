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

npm run build

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
    <script ${marker}>
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
