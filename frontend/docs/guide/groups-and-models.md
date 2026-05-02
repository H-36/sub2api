---
title: 分组与模型
icon: material-symbols:hub-rounded
description: 用一页理解 API Key、分组、模型列表和协议入口之间的关系。
---

# 分组与模型

很多接入问题看起来像“模型不可用”，其实是 Key、分组、模型列表和协议入口没有对齐。这页把它们串起来。

::: important 先记住一句话
API Key 绑定的分组决定你能看到哪些模型、走哪类协议、按什么规则计费。
:::

## 四个概念

| 概念 | 作用 | 出错时常见表现 |
| --- | --- | --- |
| API Key | 请求身份和额度入口 | `401`、未授权、用量查不到 |
| 分组 | 决定平台、模型、倍率和额度限制 | 模型列表为空、请求被拒绝 |
| 模型列表 | 当前 Key 的实时可用模型 | 手填模型名失败 |
| 协议入口 | 请求路径和请求体格式 | `/v1/messages` 和 `/v1/chat/completions` 混用失败 |

## 推荐判断顺序

<div class="docs-step-list">
  <div class="docs-step-item">
    <span>1</span>
    <div>
      <strong>确认 Key 绑定了分组</strong>
      <p>进入 API Keys 页面，检查这把 Key 是否选择了状态正常的分组。</p>
    </div>
  </div>
  <div class="docs-step-item">
    <span>2</span>
    <div>
      <strong>查询当前模型列表</strong>
      <p>大多数场景查 `/v1/models`；特殊客户端按对应接入页单独配置。</p>
    </div>
  </div>
  <div class="docs-step-item">
    <span>3</span>
    <div>
      <strong>复制模型名到客户端</strong>
      <p>不要照搬旧截图或别人配置里的模型名，以当前 Key 的返回结果为准。</p>
    </div>
  </div>
  <div class="docs-step-item">
    <span>4</span>
    <div>
      <strong>选择匹配的协议入口</strong>
      <p>Claude 用 Messages，OpenAI 用 Chat / Responses；特殊客户端再看对应接入页。</p>
    </div>
  </div>
</div>

## 常见匹配关系

| 使用场景 | 优先入口 | 模型列表 |
| --- | --- | --- |
| Claude Code / Anthropic 兼容客户端 | `/v1/messages` | `/v1/models` |
| Codex / OpenAI Responses 兼容客户端 | `/v1/responses` | `/v1/models` |
| OpenAI Chat Completions 客户端 | `/v1/chat/completions` | `/v1/models` |
| 生图接口 | `/v1/images/generations` 或 `/v1/images/edits` | `/v1/models` |

::: warning 不要把模型列表入口当成正式请求入口
`/v1/models` 是模型列表接口，不代表 `/v1/messages`、`/v1/chat/completions` 和 `/v1/responses` 可以随意互换。OpenAI 分组通常走 Chat / Responses，Claude / Anthropic 分组通常走 Messages。
:::

## 最小自检

```bash
curl https://puaai.xyz/v1/models \
  -H "Authorization: Bearer sk-your-api-key"
```

继续阅读：

- [第一个请求](./first-request.md)
- [模型与自检](../reference/models.md)
- [分组说明](../reference/groups.md)
