<template>
  <section class="docs-curl-builder" aria-label="Curl request builder">
    <div class="docs-curl-builder__grid">
      <label>
        <span>协议</span>
        <select v-model="protocol">
          <option value="chat">OpenAI Chat</option>
          <option value="responses">OpenAI Responses</option>
          <option value="anthropic">Anthropic Messages</option>
          <option value="gemini">Gemini 原生</option>
          <option value="models">模型列表</option>
        </select>
      </label>

      <label>
        <span>API Key</span>
        <input v-model="apiKey" type="text" placeholder="sk-your-api-key" />
      </label>

      <label>
        <span>模型</span>
        <input v-model="model" type="text" placeholder="gpt-5.4-mini" />
      </label>

      <label>
        <span>提示词</span>
        <input v-model="prompt" type="text" placeholder="Reply with: puaai ok" />
      </label>
    </div>

    <pre class="docs-curl-builder__output"><code>{{ command }}</code></pre>
  </section>
</template>

<script setup lang="ts">
import { computed, ref, watch } from 'vue'

type Protocol = 'chat' | 'responses' | 'anthropic' | 'gemini' | 'models'

const apiBaseUrl = 'https://puaai.xyz'

const protocol = ref<Protocol>('chat')
const apiKey = ref('sk-your-api-key')
const model = ref('gpt-5.4-mini')
const prompt = ref('Reply with: puaai ok')

watch(protocol, (next) => {
  if (next === 'anthropic') model.value = 'claude-sonnet-4-5'
  if (next === 'gemini') model.value = 'gemini-2.5-flash'
  if (next === 'chat' || next === 'responses') model.value = 'gpt-5.4-mini'
})

const command = computed(() => {
  const key = apiKey.value.trim() || 'sk-your-api-key'
  const modelName = model.value.trim() || 'gpt-5.4-mini'
  const text = prompt.value.trim() || 'Reply with: puaai ok'

  if (protocol.value === 'models') {
    return [
      `curl ${apiBaseUrl}/v1/models \\`,
      `  -H "Authorization: Bearer ${key}"`
    ].join('\n')
  }

  if (protocol.value === 'responses') {
    return [
      `curl ${apiBaseUrl}/v1/responses \\`,
      '  -H "Content-Type: application/json" \\',
      `  -H "Authorization: Bearer ${key}" \\`,
      `  -d '${JSON.stringify({ model: modelName, input: text }, null, 2)}'`
    ].join('\n')
  }

  if (protocol.value === 'anthropic') {
    return [
      `curl ${apiBaseUrl}/v1/messages \\`,
      '  -H "Content-Type: application/json" \\',
      `  -H "Authorization: Bearer ${key}" \\`,
      `  -d '${JSON.stringify({
        model: modelName,
        max_tokens: 256,
        messages: [{ role: 'user', content: text }]
      }, null, 2)}'`
    ].join('\n')
  }

  if (protocol.value === 'gemini') {
    return [
      `curl "${apiBaseUrl}/v1beta/models/${modelName}:generateContent" \\`,
      '  -H "Content-Type: application/json" \\',
      `  -H "Authorization: Bearer ${key}" \\`,
      `  -d '${JSON.stringify({
        contents: [{ parts: [{ text }] }]
      }, null, 2)}'`
    ].join('\n')
  }

  return [
    `curl ${apiBaseUrl}/v1/chat/completions \\`,
    '  -H "Content-Type: application/json" \\',
    `  -H "Authorization: Bearer ${key}" \\`,
    `  -d '${JSON.stringify({
      model: modelName,
      messages: [{ role: 'user', content: text }]
    }, null, 2)}'`
  ].join('\n')
})
</script>
