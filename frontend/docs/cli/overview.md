---
title: 通用前置步骤
icon: material-symbols:checklist-rounded
description: 配置 Claude Code、Codex、Gemini CLI 前先完成的环境、Key 和模型自检。
---

# 通用前置步骤

CLI 工具的配置差异很大，但前置检查基本一致。先完成这页，再进入具体工具页面，排错会轻很多。

## 你需要准备什么

| 项目 | 用途 | 建议 |
| --- | --- | --- |
| puaai 账号 | 登录主站和管理 Key | 使用 `https://puaai.xyz` |
| API Key | 写入 CLI 配置 | 先创建测试 Key |
| 可用分组 | 决定模型和协议入口 | 选择平台类型明确、状态正常的分组 |
| 模型名 | 写入 CLI 配置 | 从模型列表复制，不要手填 |
| 终端环境 | 运行 CLI | macOS / Linux / Windows / WSL 都可以 |

::: tip 后台示例优先
如果 API Keys 页面里的“使用密钥”生成了配置示例，并且和文档示例不同，以后台当前生成结果为准。
:::

## 1. 创建 API Key

进入主站：

```text
https://puaai.xyz
```

打开 **API Keys** 页面，创建一把测试 Key。第一次建议设置名称、分组、额度限制和过期时间。

<div class="docs-shot-callout">
  <strong>截图建议</strong>
  <p>这里适合放一张 puaai 后台 API Keys 列表截图，重点标出“创建 API Key”和“使用密钥”按钮。截图前请打码 Key、邮箱、余额和订单信息。</p>
</div>

## 2. 查询模型列表

OpenAI / Anthropic 兼容分组：

```bash
curl https://puaai.xyz/v1/models \
  -H "Authorization: Bearer sk-your-api-key"
```

Gemini 分组：

```bash
curl https://puaai.xyz/v1beta/models \
  -H "Authorization: Bearer sk-your-api-key"
```

## 3. 判断该选哪页

| 你要配置 | 继续阅读 |
| --- | --- |
| Claude Code | [Claude Code 接入指南](../integration/claude-code.md) |
| Codex CLI | [Codex 接入指南](../integration/codex.md) |
| Gemini CLI | [Gemini CLI 接入指南](./gemini.md) |
| Windows / WSL | [Windows / WSL 使用说明](./wsl.md) |
| OpenCode、Kilo Code、Zed、Cline 等 | [第三方调用总览](../third-party/overview.md) |

## 4. 常见环境检查

### Node.js

部分 CLI 或编辑器插件需要 Node.js。先检查：

```bash
node -v
npm -v
```

如果命令不存在，先安装 Node.js LTS，再继续配置。

### 配置目录

| 工具 | 常见配置路径 |
| --- | --- |
| Claude Code | `~/.claude/settings.json` |
| Codex | `~/.codex/config.toml` 和 `~/.codex/auth.json` |
| Gemini CLI | `~/.gemini/.env` |
| OpenCode | `~/.config/opencode/opencode.json` |

### Base URL

绝大多数 CLI 的 Base URL 填：

```text
https://puaai.xyz
```

OpenCode 一类工具如果要求填到协议前缀，可能使用：

```text
https://puaai.xyz/v1
```

Gemini 原生兼容路径通常是：

```text
https://puaai.xyz/v1beta
```

::: warning 不要把文档地址当 API 地址
`https://puaai.xyz/docs` 是文档站，不是 API Base URL。
:::
