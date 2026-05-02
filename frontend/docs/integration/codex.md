---
title: Codex 接入指南
icon: material-symbols:code-rounded
description: 使用当前站点提供的 OpenAI Responses 兼容入口接入 Codex CLI。
---

# Codex 接入指南

Codex 使用 OpenAI Responses 兼容方式接入 puaai。

## 准备信息

| 字段 | 推荐值 |
| --- | --- |
| Base URL | `https://puaai.xyz` |
| Wire API | `responses` |
| API Key | `sk-your-api-key` |
| 模型名 | 从 `/v1/models` 返回结果复制 |
| 配置文件 | `~/.codex/config.toml` 和 `~/.codex/auth.json` |

::: important Codex 的 Base URL 不要带 `/v1`
Codex Responses 兼容配置里 `base_url` 推荐填 `https://puaai.xyz`，由 Codex 自己拼接 Responses 路径。
:::

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

<div class="docs-shot-callout">
  <strong>截图建议</strong>
  <p>建议补 puaai 后台“使用密钥”里的 Codex 配置示例截图，以及本地 `config.toml` 编辑完成后的局部截图。真实 Key 只放在 `auth.json`，截图必须打码。</p>
</div>

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

启动后发一条简单消息：

```text
Reply with: puaai ok
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

## Windows / WSL 注意

Windows 原生 Codex 和 WSL 里的 Codex 配置路径不同：

| 环境 | 配置路径 |
| --- | --- |
| Windows | `%USERPROFILE%\\.codex\\config.toml` |
| WSL / Linux | `~/.codex/config.toml` |

你在哪个环境运行 `codex`，就改哪个环境里的文件。

## 常见问题

### `401` 或鉴权失败

确认 `~/.codex/auth.json` 是合法 JSON，并且 Key 没有多余空格或换行。

### 请求走错地址

确认 `base_url = "https://puaai.xyz"`，不要写成 `https://puaai.xyz/docs`。

### 模型不可用

用当前 Key 请求 `/v1/models`，把 `model` 和 `review_model` 都改成列表里真实返回的模型名。
