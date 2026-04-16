---
title: API Keys 说明
icon: mdi:key-variant
shortTitle: API Keys
description: 创建、管理和安全使用 puaai API Key。
---

# API Keys 说明

这页说明如何创建、保存和验证 puaai API Key。

## 创建 API Key

登录后进入用户后台的 **API Keys** 页面，点击“创建 API Key”。

建议第一次创建时这样填：

| 字段 | 建议 |
| --- | --- |
| 名称 | 写明用途，例如 `codex-test` |
| 分组 | 选择一个可用分组 |
| 自定义 Key | 留空时系统生成 |
| 额度限制 | 测试期建议设置 |
| 过期时间 | 测试期建议设置 |

分组会影响协议入口、可见模型和价格口径，所以创建后建议立即做一次模型列表自检。

## 保存 Key

创建成功后会看到一串类似下面的值：

```text
sk-xxxxxxxxxxxxxxxx
```

后续所有请求都把它放在请求头里：

```http
Authorization: Bearer sk-xxxxxxxxxxxxxxxx
```

## 最小验证

拿到 Key 之后，不要马上在复杂客户端里调半天，先做一个最小请求：

```bash
curl https://puaai.xyz/v1/models \
  -H "Authorization: Bearer sk-your-api-key"
```

Gemini 分组使用：

```bash
curl https://puaai.xyz/v1beta/models \
  -H "Authorization: Bearer sk-your-api-key"
```

## 工具配置

如果你要接 `Codex`、`Claude Code` 或 `OpenCode`，可以直接在 **API Keys** 页面点“使用密钥”查看当前这把 Key 的配置示例。

如果后台“使用密钥”生成的配置和文档示例不同，以后台当前生成结果为准。

## 安全建议

- 不要把 Key 暴露到公开前端
- 测试期建议设置额度限制
- 长期不用的 Key 及时删除或设置过期时间
