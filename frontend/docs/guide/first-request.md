---
title: 第一个请求
icon: material-symbols:terminal-rounded
description: 用最小 curl 请求验证 puaai 的 Base URL、鉴权和协议入口都已经配置正确。
---

# 第一个请求

这一页提供最小可用的请求模板。建议顺序很简单：**先看模型列表，再发送正式请求。**

::: important 推荐顺序
1. 先用模型列表接口确认当前 Key 可见什么模型。
2. 再复制模型名发第一条正式请求。
3. curl 跑通之后，再去接客户端或 SDK。
:::

## 通用准备

```text
Base URL: https://puaai.xyz
Authorization: Bearer sk-your-api-key
```

::: tip 文档站地址不是 Base URL
`https://puaai.xyz/docs` 是文档站；绝大多数客户端填写的应该是主域名 `https://puaai.xyz`。
:::

## 1. 模型列表自检

### Claude / OpenAI 兼容分组

```bash
curl https://puaai.xyz/v1/models \
  -H "Authorization: Bearer sk-your-api-key"
```

### Gemini 分组

```bash
curl https://puaai.xyz/v1beta/models \
  -H "Authorization: Bearer sk-your-api-key"
```

::: tip 为什么这里要分开
`/v1/models` 和 `/v1beta/models` 分别对应不同兼容层。Gemini SDK / CLI 场景优先看 `v1beta`，其余大多数 Claude / OpenAI 兼容客户端优先看 `v1`。
:::

## 2. Claude / Anthropic 兼容请求

```bash
curl https://puaai.xyz/v1/messages \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer sk-your-api-key" \
  -x POST \
  -d '{
    "model": "claude-sonnet-4-5",
    "max_tokens": 256,
    "messages": [
      { "role": "user", "content": "Reply with: puaai ok" }
    ]
  }'
```

## 3. OpenAI Chat Completions 兼容请求

```bash
curl https://puaai.xyz/v1/chat/completions \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer sk-your-api-key" \
  -x POST \
  -d '{
    "model": "gpt-5.4-mini",
    "messages": [
      { "role": "user", "content": "Reply with: puaai ok" }
    ]
  }'
```

## 4. OpenAI Responses 兼容请求

```bash
curl https://puaai.xyz/v1/responses \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer sk-your-api-key" \
  -x POST \
  -d '{
    "model": "gpt-5.4-mini",
    "input": "Reply with: puaai ok"
  }'
```

## 5. Gemini 原生请求

```bash
curl "https://puaai.xyz/v1beta/models/gemini-2.5-flash:generateContent" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer sk-your-api-key" \
  -x POST \
  -d '{
    "contents": [
      {
        "parts": [
          { "text": "Reply with: puaai ok" }
        ]
      }
    ]
  }'
```

## 6. 成功与失败的判断

### 成功时你通常会看到

- 模型列表接口返回一个模型数组
- 消息接口返回正常结果或流式片段
- 账单 / 用量后台出现记录

### 失败时最先排查

1. Base URL 填错
2. Key 没绑分组
3. 模型名不在当前分组可见范围内
4. 入口协议和分组平台不匹配

::: warning 不要一开始就在复杂客户端里盲配
像 Cline、Open WebUI、LobeChat、Claude Code、Gemini CLI 这类客户端，往往会把协议细节和请求体包装起来。先跑一遍 curl，排错成本会低很多。
:::

## 7. 接入客户端前的建议

如果你接的是 Cline、Open WebUI、LobeChat、Claude Code、Gemini CLI 一类客户端，先用 `curl` 跑通一次，再去填客户端配置，会少掉大半排错时间。
