---
title: Zed 接入指南
icon: material-symbols:edit-square-rounded
description: 在 Zed 中通过 OpenAI 兼容配置使用 puaai。
---

# Zed 接入指南

Zed 可以通过 OpenAI 兼容配置接入 puaai。不同版本的 Zed 配置项可能会调整，下面给出通用思路和示例。

## 准备信息

```text
Base URL: https://puaai.xyz/v1
API Key: sk-your-api-key
Model: 从 /v1/models 返回结果复制
```

## settings 示例

打开 Zed 设置文件，添加或调整 OpenAI 兼容配置：

```json
{
  "language_models": {
    "openai": {
      "api_url": "https://puaai.xyz/v1",
      "available_models": [
        {
          "name": "gpt-5.4-mini",
          "display_name": "puaai GPT-5.4 Mini",
          "max_tokens": 128000
        }
      ]
    }
  }
}
```

API Key 通常建议放到 Zed 支持的密钥管理位置或环境变量中，不要提交到项目仓库。

::: tip 以 Zed 当前版本为准
Zed 的 AI 配置格式可能随版本变化。如果设置项名称和示例不同，保留这三个核心值即可：`https://puaai.xyz/v1`、API Key、模型名。
:::

## 验证

先在终端确认 Key 和模型可用：

```bash
curl https://puaai.xyz/v1/models \
  -H "Authorization: Bearer sk-your-api-key"
```

再回到 Zed 里发一条简单消息：

```text
Reply with: puaai ok
```

<div class="docs-shot-callout">
  <strong>截图建议</strong>
  <p>建议补 Zed settings 页面或 JSON 配置截图，并用遮挡层隐藏 API Key。</p>
</div>
