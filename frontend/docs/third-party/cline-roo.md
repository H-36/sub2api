---
title: Cline / Roo Code 接入指南
icon: material-symbols:developer-mode-rounded
description: 在 VS Code 的 Cline 或 Roo Code 插件中使用 puaai 兼容接口。
---

# Cline / Roo Code 接入指南

Cline 和 Roo Code 都常用于 VS Code 内的 AI 编程工作流。它们通常支持 OpenAI Compatible 或 Anthropic Compatible 供应商。

## OpenAI 兼容配置

| 字段 | 推荐值 |
| --- | --- |
| Provider | OpenAI Compatible |
| Base URL | `https://puaai.xyz/v1` |
| API Key | `sk-your-api-key` |
| Model | 从 `/v1/models` 复制 |

## Anthropic 兼容配置

| 字段 | 推荐值 |
| --- | --- |
| Provider | Anthropic Compatible |
| Base URL | `https://puaai.xyz` |
| API Key | `sk-your-api-key` |
| Model | 从 `/v1/models` 复制 |

如果插件自动拼接 `/v1/messages`，Base URL 填 `https://puaai.xyz`；如果插件要求完整 OpenAI 风格地址，Base URL 填 `https://puaai.xyz/v1`。

## 推荐配置流程

1. 先用 curl 查询模型列表。
2. 打开插件设置。
3. 选择 OpenAI Compatible 或 Anthropic Compatible。
4. 填入 Base URL、API Key、模型名。
5. 先用小任务测试，例如“解释当前文件第一段代码”。

<div class="docs-shot-callout">
  <strong>截图建议</strong>
  <p>建议补 Cline / Roo Code Provider 下拉框截图，以及 Base URL、Model 填写位置截图。真实 Key 必须打码。</p>
</div>

## 常见问题

### 插件显示模型不可用

确认模型名是从当前 Key 的 `/v1/models` 返回结果复制的，不要使用别的站点或旧教程里的模型名。

### curl 正常，但插件失败

检查插件是否把 `/v1` 拼了两次。常见错误形态是请求到了 `/v1/v1/chat/completions` 或 `/v1/v1/messages`。
