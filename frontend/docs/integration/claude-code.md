---
title: Claude Code 接入指南
icon: material-symbols:smart-toy-rounded
description: 使用当前站点提供的 Anthropic 兼容入口接入 Claude Code。
---

# Claude Code 接入指南

Claude Code 使用 Anthropic 兼容方式接入 puaai。

## 临时环境变量

```bash
export ANTHROPIC_BASE_URL="https://puaai.xyz"
export ANTHROPIC_AUTH_TOKEN="sk-your-api-key"
export CLAUDE_CODE_DISABLE_NONESSENTIAL_TRAFFIC=1

claude
```

## `settings.json`

长期使用时，可以写入：

- macOS / Linux: `~/.claude/settings.json`
- Windows: `%userprofile%\\.claude\\settings.json`

```json
{
  "env": {
    "ANTHROPIC_BASE_URL": "https://puaai.xyz",
    "ANTHROPIC_AUTH_TOKEN": "sk-your-api-key",
    "CLAUDE_CODE_DISABLE_NONESSENTIAL_TRAFFIC": "1",
    "CLAUDE_CODE_ATTRIBUTION_HEADER": "0"
  }
}
```

## 验证

```bash
curl https://puaai.xyz/v1/models \
  -H "Authorization: Bearer sk-your-api-key"
```

```bash
claude
```

## Antigravity 分组

如果当前 Key 走的是 Antigravity 分组，Base URL 改成：

```bash
export ANTHROPIC_BASE_URL="https://puaai.xyz/antigravity"
export ANTHROPIC_AUTH_TOKEN="sk-your-api-key"
export CLAUDE_CODE_DISABLE_NONESSENTIAL_TRAFFIC=1
```

如果后台 **API Keys -> 使用密钥** 生成的 Claude Code 配置和文档示例不同，以后台当前生成结果为准。
