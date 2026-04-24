---
title: API 地址
icon: material-symbols:api-rounded
description: 汇总 puaai 的 Base URL、常用兼容入口和最小请求示例。
---

# API 地址

puaai 提供 OpenAI、Anthropic 和 Gemini 兼容接口。大多数官方 SDK 或客户端只需要把 Base URL 指向 puaai，再填入 API Key 即可接入。

## Base URL

| 场景 | 地址 |
| --- | --- |
| OpenAI / Anthropic 兼容客户端 | `https://puaai.xyz` |
| Gemini 兼容客户端 | `https://puaai.xyz` |
| 文档站 | `https://puaai.xyz/docs` |

::: important 不要把 `/docs` 当成 API Base URL
`https://puaai.xyz/docs` 是文档站地址，不是网关入口。实际发请求时，Base URL 应该填主域名，即 `https://puaai.xyz`。
:::

## 常用接口

| 类型 | 路径 | 说明 |
| --- | --- | --- |
| OpenAI Responses | `/v1/responses` | OpenAI Responses 兼容入口 |
| OpenAI Chat | `/v1/chat/completions` | OpenAI Chat Completions 兼容入口 |
| OpenAI Images | `/v1/images/generations` | OpenAI 兼容图片生成入口 |
| OpenAI Image Edits | `/v1/images/edits` | OpenAI 兼容图片编辑入口 |
| Claude / Anthropic | `/v1/messages` | Claude Messages 兼容入口 |
| 模型列表 | `/v1/models` | OpenAI / Anthropic 兼容分组常用 |
| Gemini 模型列表 | `/v1beta/models` | Gemini 分组常用 |
| Gemini 请求 | `/v1beta/models/{model}:generateContent` | Gemini 原生兼容入口 |

## 请求示例

如果你不确定当前可用模型，先用模型列表接口自检。

### 模型列表

```bash
curl https://puaai.xyz/v1/models \
  -H "Authorization: Bearer sk-your-api-key"
```

Gemini 分组使用：

```bash
curl https://puaai.xyz/v1beta/models \
  -H "Authorization: Bearer sk-your-api-key"
```

### OpenAI Responses

```bash
curl https://puaai.xyz/v1/responses \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer sk-your-api-key" \
  -d '{
    "model": "your-model-name",
    "input": "Reply with: puaai ok"
  }'
```

### Anthropic Messages

```bash
curl https://puaai.xyz/v1/messages \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer sk-your-api-key" \
  -d '{
    "model": "your-model-name",
    "max_tokens": 128,
    "messages": [
      { "role": "user", "content": "Reply with: puaai ok" }
    ]
  }'
```

`your-model-name` 请替换成当前 Key 在 `/v1/models` 或 `/v1beta/models` 里实际返回的模型名。
