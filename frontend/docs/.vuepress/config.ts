import { viteBundler } from '@vuepress/bundler-vite'
import { defineUserConfig } from 'vuepress'
import { hopeTheme } from 'vuepress-theme-hope'
import path from 'node:path'

import { navbar } from './navbar'
import { sidebar } from './sidebar'

export default defineUserConfig({
  lang: 'zh-CN',
  title: 'puaai 文档',
  description: 'puaai 对外接入文档，包含 API 地址、计费与价格、API Keys 和 Codex / Claude Code / OpenCode 接入说明。',
  base: '/docs/',
  dest: path.resolve(__dirname, '../../../backend/internal/web/dist/docs'),
  head: [
    ['meta', { name: 'theme-color', content: '#0f766e' }],
    ['meta', { property: 'og:title', content: 'puaai 文档' }],
    ['meta', { property: 'og:type', content: 'website' }],
    ['meta', { property: 'og:url', content: 'https://puaai.xyz/docs' }],
    ['meta', { property: 'og:description', content: '简洁的 puaai 接入文档。' }]
  ],
  theme: hopeTheme({
    hostname: 'https://puaai.xyz',
    favicon: 'logo.png',
    logo: 'logo.png',
    navbar,
    sidebar,
    breadcrumb: false,
    breadcrumbIcon: false,
    contributors: false,
    editLink: false,
    lastUpdated: false,
    pageInfo: [],
    titleIcon: true,
    toc: {
      levels: [2, 3]
    },
    displayFooter: true,
    footer: 'puaai Docs',
    copyright: 'Copyright © 2026 puaai',
    plugins: {
      icon: {
        assets: 'iconify'
      },
      slimsearch: {
        indexContent: true
      }
    }
  }),
  bundler: viteBundler()
})
