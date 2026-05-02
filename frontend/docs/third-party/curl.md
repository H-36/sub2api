---
title: Curl / API 测试
icon: material-symbols:http-rounded
description: 用最小 curl 请求验证 puaai 的 Base URL、Key、模型和协议入口。
---

# Curl / API 测试

配置任何客户端前，建议先用 curl 跑通。这样能把问题收敛到 Base URL、Key、模型名和协议入口。

## 一键生成请求

<CurlBuilder />

## 模型列表

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

## OpenAI Chat Completions

```bash
curl https://puaai.xyz/v1/chat/completions \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer sk-your-api-key" \
  -d '{
    "model": "gpt-5.4-mini",
    "messages": [
      { "role": "user", "content": "Reply with: puaai ok" }
    ]
  }'
```

## OpenAI Responses

```bash
curl https://puaai.xyz/v1/responses \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer sk-your-api-key" \
  -d '{
    "model": "gpt-5.4-mini",
    "input": "Reply with: puaai ok"
  }'
```

## Anthropic Messages

```bash
curl https://puaai.xyz/v1/messages \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer sk-your-api-key" \
  -d '{
    "model": "claude-sonnet-4-5",
    "max_tokens": 256,
    "messages": [
      { "role": "user", "content": "Reply with: puaai ok" }
    ]
  }'
```

## Gemini 原生

```bash
curl "https://puaai.xyz/v1beta/models/gemini-2.5-flash:generateContent" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer sk-your-api-key" \
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

## 成功判断

- 模型列表返回模型数组
- 消息接口返回文本结果
- 后台用量记录出现新的请求

## 失败时先看

| 现象 | 优先排查 |
| --- | --- |
| `401` | Key 是否复制完整，是否带 `Bearer` |
| 模型不存在 | 是否从当前 Key 的模型列表复制模型名 |
| 路径不存在 | Base URL 是否误填 `/docs` |
| 请求体错误 | 当前入口协议是否匹配 |
