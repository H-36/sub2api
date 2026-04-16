---
title: OpenCode 接入指南
icon: material-symbols:deployed-code-rounded
description: 使用当前站点的兼容入口接入 OpenCode，并根据分组类型选择合适的 provider。
---

# OpenCode 接入指南

OpenCode 需要根据当前分组类型选择不同的 provider 和 Base URL。

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

如果后台 **API Keys -> 使用密钥** 生成的 OpenCode 配置和文档示例不同，以后台当前生成结果为准。
