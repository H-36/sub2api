---
title: Hermes Agent 接入指南
icon: material-symbols:automation-rounded
description: 在 Hermes Agent 或类似 Agent 工具中使用 puaai 的 OpenAI 兼容接口。
---

# Hermes Agent 接入指南

Hermes Agent 或类似 Agent 工具通常可以按 OpenAI 兼容服务配置。由于不同版本配置文件差异较大，这页给出通用字段和排错原则。

## 推荐配置

| 字段 | 推荐值 |
| --- | --- |
| Provider | OpenAI Compatible |
| Base URL | `https://puaai.xyz/v1` |
| API Key | `sk-your-api-key` |
| Model | 从 `/v1/models` 返回结果复制 |

## 通用 JSON 示例

如果工具支持 JSON 配置，可以参考：

```json
{
  "provider": "openai",
  "baseURL": "https://puaai.xyz/v1",
  "apiKey": "sk-your-api-key",
  "model": "gpt-5.4-mini"
}
```

如果工具把 Base URL 和路径分开配置，Base URL 填：

```text
https://puaai.xyz
```

路径使用：

```text
/v1/chat/completions
```

## Agent 工具额外注意

- 第一次测试先关闭复杂工具调用，只发普通文本消息。
- 如果需要函数调用、文件编辑、浏览器控制等能力，确认当前模型和客户端都支持。
- 长任务建议设置合理超时时间，避免客户端过早断开。

<div class="docs-shot-callout">
  <strong>截图建议</strong>
  <p>如果你后续确定 Hermes Agent 的固定配置界面，建议补 Provider、Base URL、API Key、Model 四个字段的截图。</p>
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
