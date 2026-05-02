---
title: Gemini CLI 接入指南
icon: material-symbols:diamond-rounded
description: 使用 puaai 的 Gemini 原生兼容入口配置 Gemini CLI。
---

# Gemini CLI 接入指南

Gemini CLI 建议走 puaai 的 Gemini 原生兼容入口。配置前先用当前 Key 查询 `/v1beta/models`，确认可用模型。

## 前置检查

```bash
curl https://puaai.xyz/v1beta/models \
  -H "Authorization: Bearer sk-your-api-key"
```

从返回结果里复制模型名，例如：

```text
gemini-2.5-flash
```

## 环境变量方式

如果你的 Gemini CLI 支持从环境变量读取配置，可以在当前 shell 写入：

```bash
export GEMINI_API_KEY="sk-your-api-key"
export GOOGLE_GEMINI_BASE_URL="https://puaai.xyz/v1beta"
export GEMINI_MODEL="gemini-2.5-flash"
```

然后运行：

```bash
gemini
```

## `~/.gemini/.env`

长期使用时，可以创建或编辑：

```text
~/.gemini/.env
```

写入：

```bash
GEMINI_API_KEY=sk-your-api-key
GOOGLE_GEMINI_BASE_URL=https://puaai.xyz/v1beta
GEMINI_MODEL=gemini-2.5-flash
```

如果你的 Gemini CLI 版本使用不同变量名，请优先以当前 CLI 的官方说明或后台“使用密钥”生成结果为准。

## 验证请求

先跑最小 API：

```bash
curl "https://puaai.xyz/v1beta/models/gemini-2.5-flash:generateContent" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer sk-your-api-key" \
  -d '{
    "contents": [
      {
        "parts": [
          { "text": "Reply with: puaai ok" }
        ]
      }
    ]
  }'
```

再启动 CLI：

```bash
gemini
```

## 截图教程建议

<div class="docs-shot-callout">
  <strong>截图建议</strong>
  <p>建议补三张 puaai 自有截图：API Keys 页面复制 Key、模型列表自检返回、Gemini CLI 成功回复。不要使用第三方站点截图。</p>
</div>

## 常见问题

### 模型列表能查到，但 CLI 失败

优先检查 CLI 写入的模型名是不是从 `/v1beta/models` 复制的，以及 Base URL 是否误填成了 `https://puaai.xyz/docs`。

### 需要用 OpenAI 兼容而不是 Gemini 原生

如果客户端只支持 OpenAI 兼容接口，就不要走这页，改用 `/v1` 入口和 OpenAI 兼容模型。
