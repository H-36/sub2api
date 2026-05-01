#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(cd "${SCRIPT_DIR}/../.." && pwd)"
APP_DIR="${REPO_ROOT}/third_party/gpt_image_playground"
FRONTEND_DIR="${REPO_ROOT}/frontend"
OUTPUT_DIR="${FRONTEND_DIR}/public/image-playground-app"

if [[ ! -f "${APP_DIR}/package.json" ]]; then
  echo "gpt_image_playground source not found at ${APP_DIR}" >&2
  exit 1
fi

if ! command -v npm >/dev/null 2>&1; then
  echo "npm is required to build gpt_image_playground" >&2
  exit 1
fi

cd "${APP_DIR}"

if [[ ! -d node_modules ]]; then
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

if [[ -f "${OUTPUT_DIR}/index.html" ]]; then
  perl -0pi -e 's#[ \t]*<link[^>]+rel="manifest"[^>]*>\n?##g; s#[ \t]*<link[^>]+rel="apple-touch-icon"[^>]*>\n?##g' "${OUTPUT_DIR}/index.html"

  IMAGE_PLAYGROUND_INDEX_HTML="${OUTPUT_DIR}/index.html" node <<'NODE'
const fs = require('node:fs')

const file = process.env.IMAGE_PLAYGROUND_INDEX_HTML
if (!file || !fs.existsSync(file)) {
  process.exit(0)
}

const marker = 'data-sub2api-image-playground-theme'
let html = fs.readFileSync(file, 'utf8')

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

theme_check_failed=0
if [[ -f "${OUTPUT_DIR}/index.html" ]]; then
  perl -0ne 'exit(/data-sub2api-image-playground-theme/ ? 0 : 1)' "${OUTPUT_DIR}/index.html" || theme_check_failed=1
fi
while IFS= read -r -d '' file; do
  perl -0ne 'exit(/prefers-color-scheme\s*:\s*dark/ ? 1 : 0)' "$file" || theme_check_failed=1
done < <(find "${OUTPUT_DIR}/assets" -type f -name '*.css' -print0 2>/dev/null || true)
if [[ "${theme_check_failed}" -ne 0 ]]; then
  echo "Sub2API image playground theme bridge was not applied" >&2
  exit 1
fi

while IFS= read -r -d '' file; do
  perl -0pi -e 's#if\s*\(\s*["'\'']serviceWorker["'\'']\s*in\s*navigator\s*\)\s*\{#if (false && "serviceWorker" in navigator) {#g; s#(["'\'']serviceWorker["'\'']\s*in\s*navigator)\s*&&#false && $1 &&#g' "$file"
done < <(find "${OUTPUT_DIR}/assets" -type f -name '*.js' -print0 2>/dev/null || true)

registration_check_failed=0
while IFS= read -r -d '' file; do
  perl -0ne '
    while (/(["\047]serviceWorker["\047]\s*in\s*navigator\s*&&)/g) {
      my $prefix = substr($_, 0, $-[0]);
      exit 1 unless $prefix =~ /false\s*&&\s*$/;
    }
    while (/if\s*\(\s*(["\047]serviceWorker["\047]\s*in\s*navigator)/g) {
      my $snippet = substr($_, $-[0], 80);
      exit 1 unless $snippet =~ /^if\s*\(\s*false\s*&&/;
    }
  ' "$file" || registration_check_failed=1
done < <(find "${OUTPUT_DIR}/assets" -type f -name '*.js' -print0 2>/dev/null || true)

if [[ "${registration_check_failed}" -ne 0 ]]; then
  echo "service worker registration guard was not patched" >&2
  exit 1
fi

echo "Image Playground built into ${OUTPUT_DIR}"
