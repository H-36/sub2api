---
title: Windows / WSL 使用说明
icon: material-symbols:desktop-windows-rounded
description: 在 Windows、PowerShell 和 WSL 中配置 puaai CLI 工具时最常见的路径和环境变量问题。
---

# Windows / WSL 使用说明

Windows 用户常见问题不在 puaai API 本身，而在配置文件路径、环境变量作用域和 WSL 网络环境。这里按场景整理。

## 选择运行环境

| 场景 | 建议 |
| --- | --- |
| Codex / Claude Code 做开发任务 | 优先 WSL 或类 Unix shell |
| 只做简单 API 测试 | PowerShell、CMD、Git Bash 都可以 |
| 编辑器插件 | 按插件实际运行环境填写配置 |

## PowerShell 临时变量

PowerShell 当前窗口临时生效：

```powershell
$env:OPENAI_API_KEY="sk-your-api-key"
$env:ANTHROPIC_AUTH_TOKEN="sk-your-api-key"
```

验证：

```powershell
curl.exe https://puaai.xyz/v1/models `
  -H "Authorization: Bearer sk-your-api-key"
```

## WSL 临时变量

WSL / Linux shell：

```bash
export OPENAI_API_KEY="sk-your-api-key"
export ANTHROPIC_AUTH_TOKEN="sk-your-api-key"
```

验证：

```bash
curl https://puaai.xyz/v1/models \
  -H "Authorization: Bearer sk-your-api-key"
```

## 配置文件路径

| 工具 | Windows 常见路径 | WSL 常见路径 |
| --- | --- | --- |
| Claude Code | `%USERPROFILE%\\.claude\\settings.json` | `~/.claude/settings.json` |
| Codex | `%USERPROFILE%\\.codex\\config.toml` | `~/.codex/config.toml` |
| Gemini CLI | `%USERPROFILE%\\.gemini\\.env` | `~/.gemini/.env` |

::: warning Windows 和 WSL 配置不共享
Windows 下的 `%USERPROFILE%\\.codex` 和 WSL 里的 `~/.codex` 是两套目录。你在哪个环境运行 CLI，就改哪个环境里的配置。
:::

## 常见排错

### PowerShell 里 `$` 被吞掉

如果命令或 JSON 里有 `$`，PowerShell 可能把它当变量展开。复杂请求建议放进文件，或使用 WSL / Git Bash 执行。

### `curl` 行尾续行失败

PowerShell 续行符是反引号：

```powershell
curl.exe https://puaai.xyz/v1/models `
  -H "Authorization: Bearer sk-your-api-key"
```

Linux / WSL 续行符是反斜杠：

```bash
curl https://puaai.xyz/v1/models \
  -H "Authorization: Bearer sk-your-api-key"
```

### 编辑器插件和终端配置不一致

很多编辑器插件不会读取你终端里的临时 `export`。这种情况请在插件设置里直接填写 Base URL、API Key 和模型名。
