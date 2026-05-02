import { defineClientConfig } from 'vuepress/client'

import CurlBuilder from './components/CurlBuilder.vue'

export default defineClientConfig({
  enhance({ app }) {
    app.component('CurlBuilder', CurlBuilder)
  }
})
