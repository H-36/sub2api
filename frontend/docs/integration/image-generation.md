---
title: 生图接口接入指南
icon: material-symbols:image-rounded
description: 使用 puaai 的 OpenAI 兼容生图接口生成或编辑图片。
---

# 生图接口接入指南

puaai 支持 OpenAI 兼容的图片生成接口。日常接入建议优先使用 `/v1/images/generations`，需要图片编辑时使用 `/v1/images/edits`。

::: tip 先确认模型
不同 API Key 绑定的分组可能不同。接入前建议先请求 `/v1/models`，确认当前 Key 能看到 `gpt-image-2` 或其他 `gpt-image-*` 模型。
:::

## 接口地址

| 场景 | Method | Path |
| --- | --- | --- |
| 生成图片 | `POST` | `/v1/images/generations` |
| 编辑图片 | `POST` | `/v1/images/edits` |
| 高级 Responses 工具调用 | `POST` | `/v1/responses` |

Base URL 是 `https://puaai.xyz`，不要把 `/docs` 写进 API Base URL。

## 图片生成

最小请求只需要 `model` 和 `prompt`。推荐显式传入 `response_format: "b64_json"`，客户端拿到 base64 后自行保存为图片。

```bash
curl https://puaai.xyz/v1/images/generations \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer sk-your-api-key" \
  -d '{
    "model": "gpt-image-2",
    "prompt": "生成一张浅色背景的产品海报：一只红色马克杯放在木桌上。",
    "size": "1024x1024",
    "quality": "low",
    "output_format": "png",
    "response_format": "b64_json"
  }'
```

返回示例：

```json
{
  "created": 1776984091,
  "data": [
    {
      "b64_json": "iVBORw0KGgo..."
    }
  ],
  "background": "auto",
  "output_format": "png",
  "quality": "auto",
  "size": "1024x1024",
  "model": "gpt-image-2"
}
```

## Python 示例

```python
import base64
import os

import requests

api_key = os.environ["PUAAI_API_KEY"]

resp = requests.post(
    "https://puaai.xyz/v1/images/generations",
    headers={
        "Authorization": f"Bearer {api_key}",
        "Content-Type": "application/json",
    },
    json={
        "model": "gpt-image-2",
        "prompt": "生成一张浅色背景的产品海报：一只红色马克杯放在木桌上。",
        "size": "1024x1024",
        "quality": "low",
        "output_format": "png",
        "response_format": "b64_json",
    },
    timeout=180,
)
resp.raise_for_status()

result = resp.json()
image_base64 = result["data"][0]["b64_json"]

with open("output.png", "wb") as f:
    f.write(base64.b64decode(image_base64))

print("saved output.png")
```

## 常用参数

| 参数 | 类型 | 说明 |
| --- | --- | --- |
| `model` | string | 图片模型，例如 `gpt-image-2`。 |
| `prompt` | string | 生图提示词，必填。 |
| `size` | string | 图片尺寸，常用 `1024x1024`、`1536x1024`、`1024x1536`，也可使用 `auto`。 |
| `quality` | string | 图片质量，常用 `low`、`medium`、`high` 或 `auto`。 |
| `response_format` | string | 返回格式，推荐 `b64_json`。如使用 `url`，网关可能返回 data URL。 |
| `output_format` | string | 输出格式，常用 `png`、`jpeg`、`webp`。 |
| `background` | string | 背景模式，例如 `auto`、`opaque`、`transparent`。 |
| `n` | number | 生成数量。部分后端通道可能只返回 1 张，批量生成建议客户端循环请求。 |
| `stream` | boolean | 是否使用流式返回。图片接口支持流式事件，普通接入可以先用非流式。 |

## 高级参数

| 参数 | 类型 | 说明 |
| --- | --- | --- |
| `moderation` | string | 内容审核策略，例如 `auto`。 |
| `style` | string | 风格参数，是否生效取决于上游模型。 |
| `output_compression` | number | 输出压缩质量，主要用于 `jpeg`、`webp`。 |
| `partial_images` | number | 流式生成时返回中间预览图的数量。 |
| `input_fidelity` | string | 图片编辑时的输入保真度，例如 `high`。 |

## 图片编辑

JSON 方式可以传远程图片地址：

```bash
curl https://puaai.xyz/v1/images/edits \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer sk-your-api-key" \
  -d '{
    "model": "gpt-image-2",
    "prompt": "把背景替换成夜晚城市霓虹灯，主体保持不变。",
    "images": [
      {
        "image_url": "https://example.com/source.png"
      }
    ],
    "mask": {
      "image_url": "https://example.com/mask.png"
    },
    "size": "1024x1024",
    "output_format": "png",
    "response_format": "b64_json"
  }'
```

如果图片在本地，可以用 multipart：

```bash
curl https://puaai.xyz/v1/images/edits \
  -H "Authorization: Bearer sk-your-api-key" \
  -F "model=gpt-image-2" \
  -F "prompt=把背景替换成夜晚城市霓虹灯，主体保持不变。" \
  -F "image=@source.png" \
  -F "mask=@mask.png" \
  -F "output_format=png" \
  -F "response_format=b64_json"
```

## 流式返回

图片接口使用 `stream: true` 时会返回 SSE。生成过程中可能出现预览图事件，最终图在 completed 事件中返回。

```bash
curl https://puaai.xyz/v1/images/generations \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer sk-your-api-key" \
  -d '{
    "model": "gpt-image-2",
    "prompt": "生成一张极简风格的蓝色几何海报。",
    "stream": true,
    "partial_images": 1,
    "response_format": "b64_json"
  }'
```

常见事件：

| 事件 | 说明 |
| --- | --- |
| `image_generation.partial_image` | 中间预览图，读取 `b64_json`。 |
| `image_generation.completed` | 最终图片，读取 `b64_json` 或 `url`。 |
| `error` | 上游或网关错误。 |

## Responses 生图

`/v1/responses` 不是专用图片接口，但可以通过 `image_generation` 工具触发生图。这个方式适合已经使用 Responses API 的客户端。

```bash
curl https://puaai.xyz/v1/responses \
  -H "Content-Type: application/json" \
  -H "Accept: text/event-stream" \
  -H "Authorization: Bearer sk-your-api-key" \
  -d '{
    "model": "gpt-5.4-mini",
    "input": [
      {
        "type": "message",
        "role": "user",
        "content": [
          {
            "type": "input_text",
            "text": "生成一张浅色背景的产品海报：一只红色马克杯放在木桌上。"
          }
        ]
      }
    ],
    "stream": true,
    "tool_choice": {
      "type": "image_generation"
    },
    "tools": [
      {
        "type": "image_generation",
        "action": "generate",
        "model": "gpt-image-2",
        "size": "1024x1024",
        "quality": "low",
        "output_format": "png"
      }
    ]
  }'
```

::: warning Responses 的模型选择
`/v1/responses` 顶层 `model` 应该填写支持 Responses 的文本模型，例如 `gpt-5.4-mini`。图片模型应放在 `tools[0].model` 中，不要把 `gpt-image-2` 写到顶层 `model`。
:::

## 排查建议

- 返回 `401`：检查 `Authorization: Bearer sk-your-api-key` 是否正确。
- 返回模型不可用：先请求 `/v1/models`，确认当前 Key 所属分组包含 `gpt-image-*`。
- 返回体很大：图片结果通常是 base64，客户端超时时间建议设置到 180 秒以上。
- 没拿到图片：非流式读取 `data[0].b64_json`；流式读取 `image_generation.completed` 事件里的 `b64_json`。
