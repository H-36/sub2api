---
title: 获取 API Key
icon: mdi:key-variant
shortTitle: 获取 Key
description: 首次使用 puaai 时，如何尽快创建一把可用的 API Key。
---

# 获取 API Key

如果你刚注册或刚登录 puaai，这一页会带你用最短路径创建第一把可用的 API Key，并完成首次可用性检查。

::: important 先记住最短路径
1. 登录主站。
2. 进入 `API Keys` 页面创建 Key。
3. 给 Key 绑定一个有效分组。
4. 用这把 Key 先请求一次模型列表。
:::

## 1. 注册并登录

访问 `https://puaai.xyz`，完成注册或登录。

如果站点开启了以下能力，你可能还会遇到：

- 邮箱验证码
- 邀请码
- 优惠码
- 第三方登录

这些都是站点级设置，不影响网关协议本身。

## 2. 进入 API Keys 页面

登录后进入用户后台的 **API Keys** 页面，点击“创建 API Key”。

建议第一次创建时这样填：

| 字段 | 建议 |
| --- | --- |
| 名称 | 写明用途，例如 `openwebui-test` |
| 分组 | 先选一个当前可用、自己知道用途的分组 |
| 自定义 Key | 不填，先使用系统生成 |
| 额度限制 | 测试期可以先设一个小额上限 |
| 过期时间 | 测试期建议设置，方便控风险 |
| IP 白名单 | 服务端固定出口时再启用 |

::: tip 分组就是 Key 的“能力边界”
分组不仅影响价格倍率，也影响平台协议、可见模型、是否支持某些特定入口，以及是否允许订阅模式或日 / 周 / 月限额。
:::

## 3. 复制并保存 Key

创建成功后会看到一串类似下面的值：

```text
sk-xxxxxxxxxxxxxxxx
```

后续所有请求都把它放在请求头里：

```http
Authorization: Bearer sk-xxxxxxxxxxxxxxxx
```

::: warning 只有拿到 `sk-...` 还不够
如果这把 Key 没有绑定有效分组，网关仍然会直接拦截请求。首次接入最常见的误区，就是“以为创建成功就等于可用”。
:::

## 4. 先做最小自检

拿到 Key 之后，不要马上在复杂客户端里调半天，先做一个最小请求：

```bash
curl https://puaai.xyz/v1/models \
  -H "Authorization: Bearer sk-your-api-key"
```

如果返回正常，基本说明你这把 Key 已经可以开始接入。

::: tip 为什么先查 `/v1/models`
这一步能同时验证 Key 是否有效、分组是否绑定成功、Base URL 是否正确，以及当前分组下是否至少存在可调度模型。
:::

## 5. 看不懂分组时怎么选

第一次最容易卡住的其实不是“拿 Key”，而是“拿到 Key 后选哪个分组”。

建议按这个顺序判断：

1. 先看分组的平台类型，是 `anthropic`、`openai`、`gemini`、`antigravity` 还是 `sora`
2. 再看这个分组主要给哪类客户端用
3. 最后看倍率、限额、是否专属、是否需要订阅

更具体说明见 [分组说明](../reference/groups.md)。

## 6. 常见卡点

### Key 创建成功，但请求直接报错

优先排查：

- 这把 Key 是否绑定了分组
- 分组是否处于 `active`
- 你请求的入口是否和分组平台匹配

### Key 能创建，但在客户端里不好用

优先排查：

- Base URL 是否填成了 `https://puaai.xyz`
- 模型名是不是直接手填错了
- 客户端走的是 `messages`、`responses`、`chat/completions` 还是 `v1beta`

### 想查这把 Key 的额度和用量

可以直接用站点提供的公开查询页：

`https://puaai.xyz/key-usage`
