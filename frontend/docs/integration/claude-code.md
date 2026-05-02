---
title: Claude Code 接入指南
icon: material-symbols:smart-toy-rounded
description: 使用当前站点提供的 Anthropic 兼容入口接入 Claude Code。
---

# Claude Code 接入指南

Claude Code 使用 Anthropic 兼容方式接入 puaai。

::: tip 先拿当前 Key 自检
配置前建议先请求 `/v1/models`，确认这把 Key 当前能看到 Claude / Anthropic 兼容模型。
:::

## 准备信息

| 字段 | 推荐值 |
| --- | --- |
| Base URL | `https://puaai.xyz` |
| Auth Token | `sk-your-api-key` |
| 模型名 | 从 `/v1/models` 返回结果复制 |
| 配置文件 | `~/.claude/settings.json` |

## 临时环境变量

```bash
export ANTHROPIC_BASE_URL="https://puaai.xyz"
export ANTHROPIC_AUTH_TOKEN="sk-your-api-key"
export CLAUDE_CODE_DISABLE_NONESSENTIAL_TRAFFIC=1

claude
```

这个方式只对当前终端窗口生效，适合临时测试。

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

写完后重新打开 Claude Code，让配置重新加载。

<div class="docs-shot-callout">
  <strong>截图建议</strong>
  <p>建议补两张 puaai 自有截图：后台“使用密钥”生成 Claude Code 配置、终端里 Claude Code 成功启动并回复。截图前请打码 API Key。</p>
</div>

## 验证

```bash
curl https://puaai.xyz/v1/models \
  -H "Authorization: Bearer sk-your-api-key"
```

```bash
claude
```

进入 Claude Code 后发一条简单消息：

```text
Reply with: puaai ok
```

## Antigravity 分组

如果当前 Key 走的是 Antigravity 分组，Base URL 改成：

```bash
export ANTHROPIC_BASE_URL="https://puaai.xyz/antigravity"
export ANTHROPIC_AUTH_TOKEN="sk-your-api-key"
export CLAUDE_CODE_DISABLE_NONESSENTIAL_TRAFFIC=1
```

如果后台 **API Keys -> 使用密钥** 生成的 Claude Code 配置和文档示例不同，以后台当前生成结果为准。

## 常见问题

### curl 可以，但 Claude Code 不走 puaai

检查 `~/.claude/settings.json` 是否写在当前系统用户目录下。Windows 和 WSL 是两套目录，在哪个环境运行 Claude Code，就改哪个环境的配置。

### 请求仍然访问官方地址

确认 `ANTHROPIC_BASE_URL` 已写入 Claude Code 能读取的环境中。临时 `export` 只对当前终端有效。

### 模型不可用

先查 `/v1/models`，从返回结果复制模型名。不要照搬其他站点或旧截图里的模型名。
