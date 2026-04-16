---
title: Codex 接入指南
icon: material-symbols:code-rounded
description: 使用当前站点提供的 OpenAI Responses 兼容入口接入 Codex CLI。
---

# Codex 接入指南

Codex 使用 OpenAI Responses 兼容方式接入 puaai。

## `~/.codex/config.toml`

```toml
model_provider = "OpenAI"
model = "gpt-5.4"
review_model = "gpt-5.4"
model_reasoning_effort = "xhigh"
disable_response_storage = true
network_access = "enabled"
windows_wsl_setup_acknowledged = true
model_context_window = 1000000
model_auto_compact_token_limit = 900000

[model_providers.OpenAI]
name = "OpenAI"
base_url = "https://puaai.xyz"
wire_api = "responses"
requires_openai_auth = true
```

如果当前分组看不到 `gpt-5.4`，请把 `model` 和 `review_model` 改成 `/v1/models` 实际返回的模型名。

## `~/.codex/auth.json`

```json
{
  "OPENAI_API_KEY": "sk-your-api-key"
}
```

## 验证

```bash
curl https://puaai.xyz/v1/models \
  -H "Authorization: Bearer sk-your-api-key"
```

```bash
codex
```

## 可选：开启 WebSocket v2

如果你需要 WebSocket v2，可以把 `config.toml` 改成：

```toml
model_provider = "OpenAI"
model = "gpt-5.4"
review_model = "gpt-5.4"
model_reasoning_effort = "xhigh"
disable_response_storage = true
network_access = "enabled"
windows_wsl_setup_acknowledged = true
model_context_window = 1000000
model_auto_compact_token_limit = 900000

[model_providers.OpenAI]
name = "OpenAI"
base_url = "https://puaai.xyz"
wire_api = "responses"
supports_websockets = true
requires_openai_auth = true

[features]
responses_websockets_v2 = true
```

如果后台 **API Keys -> 使用密钥** 生成的 Codex 配置和文档示例不同，以后台当前生成结果为准。
