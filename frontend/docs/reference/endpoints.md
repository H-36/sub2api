---
title: API 地址
icon: material-symbols:api-rounded
description: 汇总 puaai 的 Base URL、常用兼容入口和最小请求示例。
---

# API 地址

puaai 提供 OpenAI、Anthropic 和 Gemini 兼容接口。手写 HTTP 请求时，使用主域名加具体路径；客户端配置 Base URL 时，要看它会不会自动追加 `/v1`。

## 地址该怎么填

| 场景 | 地址 |
| --- | --- |
| API 根域名 | `https://puaai.xyz` |
| 本页 curl 示例 | `https://puaai.xyz` + 下方路径 |
| OpenAI 兼容 Base URL | 常见填 `https://puaai.xyz/v1` |
| Claude Code / Codex 配置 | 按对应接入页填 `https://puaai.xyz` |
| 文档站 | `https://puaai.xyz/docs` |

## 常用接口

| 类型 | 路径 | 说明 |
| --- | --- | --- |
| OpenAI Responses | `/v1/responses` | OpenAI Responses 兼容入口 |
| OpenAI Chat | `/v1/chat/completions` | OpenAI Chat Completions 兼容入口 |
| OpenAI Images | `/v1/images/generations` | OpenAI 兼容图片生成入口 |
| OpenAI Image Edits | `/v1/images/edits` | OpenAI 兼容图片编辑入口 |
| Claude / Anthropic | `/v1/messages` | Claude Messages 兼容入口 |
| 模型列表 | `/v1/models` | 大多数接入场景常用 |
| Gemini 原生 | `/v1beta/models` | 仅明确接 Gemini SDK / CLI 时使用 |

## 请求示例

如果你不确定当前可用模型，先用模型列表接口自检。

### 模型列表

```bash
curl https://puaai.xyz/v1/models \
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

`your-model-name` 请替换成当前 Key 在对应模型列表里实际返回的模型名。

::: warning 正式请求入口要和分组平台匹配
`/v1/models` 能返回模型列表，不代表正式请求入口可以混用。OpenAI 分组通常走 Chat / Responses；Claude / Anthropic 分组通常走 Messages。
:::
