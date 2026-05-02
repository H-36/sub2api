---
title: Kilo Code 接入指南
icon: material-symbols:code-blocks-rounded
description: 在 Kilo Code 中使用 puaai 的 OpenAI 或 Anthropic 兼容接口。
---

# Kilo Code 接入指南

Kilo Code 一般可以按 OpenAI 兼容或 Anthropic 兼容方式接入。第一次配置建议先用 OpenAI 兼容，确认能对话后再按需要切换 Claude / Anthropic 兼容。

## 准备信息

| 字段 | 推荐值 |
| --- | --- |
| API Provider | OpenAI Compatible |
| Base URL | `https://puaai.xyz/v1` |
| API Key | `sk-your-api-key` |
| Model | 从 `/v1/models` 返回结果复制 |

## 配置步骤

1. 打开 Kilo Code 设置。
2. 找到模型供应商或 API Provider。
3. 选择 OpenAI Compatible。
4. 填写 Base URL：`https://puaai.xyz/v1`。
5. 填写 API Key。
6. 填写从 `/v1/models` 复制的模型名。
7. 保存后发送一条最小消息。

<div class="docs-shot-callout">
  <strong>截图建议</strong>
  <p>建议补一张 Kilo Code Provider 设置截图，标出 Base URL、API Key、Model 三个字段。截图中不要展示真实 Key。</p>
</div>

## Anthropic 兼容方式

如果你要走 Claude / Anthropic 兼容，可以尝试：

| 字段 | 推荐值 |
| --- | --- |
| API Provider | Anthropic Compatible |
| Base URL | `https://puaai.xyz` 或 `https://puaai.xyz/v1` |
| API Key | `sk-your-api-key` |
| Model | 从 `/v1/models` 返回结果复制 |

不同版本的 Kilo Code 对 Anthropic Compatible 的 Base URL 处理可能不同。如果请求路径重复出现 `/v1/v1/messages`，就把 Base URL 从 `https://puaai.xyz/v1` 改成 `https://puaai.xyz`。

## 验证

```bash
curl https://puaai.xyz/v1/models \
  -H "Authorization: Bearer sk-your-api-key"
```

如果 curl 能通但 Kilo Code 不通，优先检查 Base URL 末尾是否多写了 `/docs` 或重复 `/v1`。
