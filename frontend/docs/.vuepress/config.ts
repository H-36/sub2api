import { viteBundler } from '@vuepress/bundler-vite'
import { defineUserConfig } from 'vuepress'
import { hopeTheme } from 'vuepress-theme-hope'
import path from 'node:path'

import { docsMeta } from './docsMeta'
import { navbar } from './navbar'
import { sidebar } from './sidebar'

export default defineUserConfig({
  lang: 'zh-CN',
  title: docsMeta.siteName,
  description: 'puaai 对外接入文档，包含快速开始、CLI 配置、第三方调用、API 地址、计费与价格和常见问题。',
  base: '/docs/',
  dest: path.resolve(__dirname, '../../../backend/internal/web/dist/docs'),
  clientConfigFile: path.resolve(__dirname, './client.ts'),
  head: [
    ['meta', { name: 'theme-color', content: '#0f766e' }],
    ['meta', { property: 'og:title', content: docsMeta.siteName }],
    ['meta', { property: 'og:type', content: 'website' }],
    ['meta', { property: 'og:url', content: docsMeta.docsUrl }],
    ['meta', { property: 'og:description', content: 'puaai 接入、CLI 配置与第三方工具使用文档。' }]
  ],
  theme: hopeTheme({
    hostname: docsMeta.siteUrl,
    favicon: 'logo.png',
    logo: 'logo.png',
    navbar,
    sidebar,
    breadcrumb: false,
    breadcrumbIcon: false,
    contributors: false,
    editLink: false,
    lastUpdated: true,
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
