---
title: 接口与 Base URL
icon: material-symbols:api-rounded
description: 汇总 puaai 的常用 Base URL、兼容协议入口和专用路径。
---

# 接口与 Base URL

这页把最常用的 Base URL、兼容协议入口和专用路径集中放在一起。接客户端、对照 SDK 参数或排查入口填错时，直接查这里即可。

## Base URL

当前站点对外文档和主站最终会挂在同一域名下：

```text
https://puaai.xyz
```

因此大多数客户端只需要填写这个 Base URL，再根据协议自动拼接路径。

::: important 不要把 `/docs` 当成 API Base URL
`https://puaai.xyz/docs` 是文档站地址，不是网关入口。实际发请求时，Base URL 应该填主域名，即 `https://puaai.xyz`。
:::

::: tip 先记一条经验规则
如果你不确定应该走哪个入口，先用当前 Key 查模型列表。能从哪个兼容层正常返回，再围绕那个兼容层继续接。
:::

## Claude / Anthropic 兼容

| 方法 | 路径 | 说明 |
| --- | --- | --- |
| `POST` | `/v1/messages` | Claude Messages 兼容入口 |
| `POST` | `/v1/messages/count_tokens` | 统计 token，部分平台 / 分组不支持 |
| `GET` | `/v1/models` | 查看当前 Key 可见模型 |
| `GET` | `/v1/usage` | 查看用量信息 |

## OpenAI 兼容

| 方法 | 路径 | 说明 |
| --- | --- | --- |
| `POST` | `/v1/chat/completions` | 旧版 OpenAI 兼容入口 |
| `POST` | `/v1/responses` | 新版 Responses API |
| `POST` | `/responses` | 不带 `/v1` 的别名入口 |
| `GET` | `/responses` | Responses WebSocket 入口 |
| `GET` | `/v1/models` | 复用通用模型列表 |

## Gemini 原生兼容

| 方法 | 路径 | 说明 |
| --- | --- | --- |
| `GET` | `/v1beta/models` | 列出当前 Key 可见模型 |
| `GET` | `/v1beta/models/:model` | 查看模型详情 |
| `POST` | `/v1beta/models/{model}:generateContent` | 生成内容 |
| `POST` | `/v1beta/models/{model}:streamGenerateContent` | 流式生成 |

## Antigravity 专用入口

Antigravity 既可能参与混合调度，也有专用路径。

| 方法 | 路径 | 说明 |
| --- | --- | --- |
| `GET` | `/antigravity/models` | Antigravity 可用模型 |
| `POST` | `/antigravity/v1/messages` | Claude 兼容专用入口 |
| `GET` | `/antigravity/v1/models` | Claude 兼容模型列表 |
| `GET` | `/antigravity/v1beta/models` | Gemini 兼容模型列表 |
| `POST` | `/antigravity/v1beta/models/*` | Gemini 兼容专用入口 |

## Sora 专用入口

| 方法 | 路径 | 说明 |
| --- | --- | --- |
| `POST` | `/sora/v1/chat/completions` | Sora 专用聊天 / 生成入口 |
| `GET` | `/sora/v1/models` | 当前 Key 下的 Sora 模型 |
| `GET` | `/sora/media/*filepath` | 媒体代理 |
| `GET` | `/sora/media-signed/*filepath` | 签名媒体代理 |

## 选入口的简单规则

### 你可以这样记

- Claude / Claude Code 一类客户端，优先看 `/v1/messages`
- OpenAI SDK / OpenAI 兼容客户端，优先看 `/v1/chat/completions` 或 `/v1/responses`
- Gemini SDK / CLI，优先看 `/v1beta/*`
- Sora 单独走 `/sora/*`

### 遇到不确定时

先用当前 Key 请求模型列表：

```bash
curl https://puaai.xyz/v1/models -H "Authorization: Bearer sk-your-api-key"
```

或：

```bash
curl https://puaai.xyz/v1beta/models -H "Authorization: Bearer sk-your-api-key"
```

能通哪个，就优先围绕那个兼容层继续接。
