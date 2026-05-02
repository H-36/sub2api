---
title: Open WebUI / LobeChat 接入指南
icon: material-symbols:chat-rounded
description: 在 Open WebUI、LobeChat 等 OpenAI 兼容聊天前端中配置 puaai。
---

# Open WebUI / LobeChat 接入指南

Open WebUI、LobeChat 这类聊天前端通常按 OpenAI 兼容接口接入。配置重点是 Base URL 要填到 `/v1`。

## OpenAI 兼容配置

| 字段 | 推荐值 |
| --- | --- |
| API Provider | OpenAI Compatible |
| Base URL / API URL | `https://puaai.xyz/v1` |
| API Key | `sk-your-api-key` |
| Model | 从 `/v1/models` 复制，或让前端自动拉取 |

## Open WebUI

常见环境变量示例：

```bash
OPENAI_API_BASE_URL=https://puaai.xyz/v1
OPENAI_API_KEY=sk-your-api-key
```

如果你在管理后台添加连接，填写：

```text
Base URL: https://puaai.xyz/v1
API Key: sk-your-api-key
```

## LobeChat

在模型供应商中添加 OpenAI 兼容服务：

```text
API Key: sk-your-api-key
Proxy URL / Base URL: https://puaai.xyz/v1
Model: 从 /v1/models 复制
```

<div class="docs-shot-callout">
  <strong>截图建议</strong>
  <p>建议分别补 Open WebUI 连接配置页、LobeChat Provider 设置页截图，标出 Base URL 和模型选择位置。</p>
</div>

## 验证

```bash
curl https://puaai.xyz/v1/chat/completions \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer sk-your-api-key" \
  -d '{
    "model": "gpt-5.4-mini",
    "messages": [
      { "role": "user", "content": "Reply with: puaai ok" }
    ]
  }'
```

## 常见问题

### 模型下拉为空

先确认 `/v1/models` 能返回结果。若 API 正常但前端不显示，可以手动添加模型名。

### 聊天前端支持图片吗？

这取决于前端是否支持 OpenAI 兼容图片接口。图片生成和编辑接口说明见 [生图接口](../integration/image-generation.md)。
