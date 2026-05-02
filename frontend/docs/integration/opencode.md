---
title: OpenCode 接入指南
icon: material-symbols:deployed-code-rounded
description: 使用当前站点的兼容入口接入 OpenCode，并根据分组类型选择合适的 provider。
---

# OpenCode 接入指南

OpenCode 需要根据当前分组类型选择不同的 provider 和 Base URL。

## 准备信息

| 字段 | 推荐值 |
| --- | --- |
| 配置文件 | `~/.config/opencode/opencode.json` |
| OpenAI Base URL | `https://puaai.xyz/v1` |
| API Key | `sk-your-api-key` |
| 模型名 | 从 `/v1/models` 或 `/v1beta/models` 复制 |

## 配置文件

```text
~/.config/opencode/opencode.json
```

## 常见配置示例

OpenAI 兼容分组可以先用下面这份配置：

```json
{
  "provider": {
    "openai": {
      "options": {
        "baseURL": "https://puaai.xyz/v1",
        "apiKey": "sk-your-api-key"
      },
      "models": {
        "gpt-5.4": {
          "name": "GPT-5.4",
          "limit": {
            "context": 1050000,
            "output": 128000
          },
          "options": {
            "store": false
          },
          "variants": {
            "low": {},
            "medium": {},
            "high": {},
            "xhigh": {}
          }
        },
        "gpt-5.3-codex": {
          "name": "GPT-5.3 Codex",
          "limit": {
            "context": 400000,
            "output": 128000
          },
          "options": {
            "store": false
          },
          "variants": {
            "low": {},
            "medium": {},
            "high": {},
            "xhigh": {}
          }
        }
      }
    }
  },
  "agent": {
    "build": {
      "options": {
        "store": false
      }
    },
    "plan": {
      "options": {
        "store": false
      }
    }
  },
  "$schema": "https://opencode.ai/config.json"
}
```

<div class="docs-shot-callout">
  <strong>截图建议</strong>
  <p>建议补 OpenCode 配置文件位置截图、模型配置片段截图和终端成功运行截图。API Key 建议放入本机安全配置或环境变量，截图时必须打码。</p>
</div>

## 分组与 provider 对应关系

| 分组类型 | 推荐 provider | Base URL |
| --- | --- | --- |
| OpenAI 兼容分组 | `openai` | `https://puaai.xyz/v1` |
| Claude / Anthropic 兼容分组 | `anthropic` | `https://puaai.xyz/v1` |
| Gemini 分组 | `gemini` | `https://puaai.xyz/v1beta` |
| Antigravity Claude | `antigravity-claude` | `https://puaai.xyz/antigravity/v1` |
| Antigravity Gemini | `antigravity-gemini` | `https://puaai.xyz/antigravity/v1beta` |

## 验证

OpenAI / Claude 兼容分组：

```bash
curl https://puaai.xyz/v1/models \
  -H "Authorization: Bearer sk-your-api-key"
```

Gemini 分组：

```bash
curl https://puaai.xyz/v1beta/models \
  -H "Authorization: Bearer sk-your-api-key"
```

启动 OpenCode 后先发小任务测试：

```text
Reply with: puaai ok
```

如果后台 **API Keys -> 使用密钥** 生成的 OpenCode 配置和文档示例不同，以后台当前生成结果为准。

## 常见问题

### Base URL 应该带 `/v1` 吗？

OpenCode 的 OpenAI provider 通常填 `https://puaai.xyz/v1`。如果你用的 provider 会自动拼 `/v1`，出现 `/v1/v1` 错误时改成 `https://puaai.xyz`。

### 多个模型怎么配？

在 `models` 里为每个模型添加一个条目，key 使用真实模型名。模型名建议从当前 Key 的模型列表复制。

### Gemini 分组怎么配？

Gemini 分组优先使用 `https://puaai.xyz/v1beta`，并用 `/v1beta/models` 查询模型。
