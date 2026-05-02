---
title: 第三方调用总览
icon: material-symbols:extension-rounded
description: 按客户端类型选择 puaai 的 Base URL、协议入口和模型自检方式。
---

# 第三方调用总览

第三方工具接入 puaai 时，核心就是三件事：Base URL、API Key、模型名。不同工具只是配置入口不同。

## 先选客户端类型

| 客户端 | 推荐页面 | 典型 Base URL |
| --- | --- | --- |
| Curl / 自写请求 | [Curl / API 测试](./curl.md) | `https://puaai.xyz` |
| OpenCode | [OpenCode 接入指南](../integration/opencode.md) | `https://puaai.xyz/v1` |
| Kilo Code | [Kilo Code 接入指南](./kilo-code.md) | `https://puaai.xyz/v1` |
| Zed | [Zed 接入指南](./zed.md) | `https://puaai.xyz/v1` |
| Cline / Roo Code | [Cline / Roo Code 接入指南](./cline-roo.md) | `https://puaai.xyz/v1` |
| Open WebUI / LobeChat | [Open WebUI / LobeChat](./webui-lobechat.md) | `https://puaai.xyz/v1` |
| Hermes Agent | [Hermes Agent](./hermes.md) | `https://puaai.xyz/v1` |

## 接入前固定动作

1. 创建 API Key，并绑定可用分组。
2. 用当前 Key 查询模型列表。
3. 从模型列表复制模型名。
4. 在客户端里填写 Base URL、API Key、模型名。
5. 发一条最小消息测试。

## 协议选择

| 分组或工具类型 | 推荐协议 |
| --- | --- |
| OpenAI 兼容工具 | OpenAI Chat 或 Responses |
| Claude / Anthropic 兼容工具 | Anthropic Messages |
| Gemini 工具 | Gemini 原生 `/v1beta` |
| 不确定工具能力 | 先用 OpenAI Chat 兼容配置测试 |

::: tip 最小化变量
第一次接入时先关闭复杂功能，例如工具调用、图片、长上下文、代理插件和自定义 header。先跑通最小文本请求，再逐步打开高级能力。
:::
