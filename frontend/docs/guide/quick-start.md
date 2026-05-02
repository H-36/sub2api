---
title: 快速开始
icon: material-symbols:rocket-launch-rounded
description: 用最短路径完成 puaai 首次接入。
---

# 快速开始

第一次接入 puaai，按下面 4 步走就够了。

## 1. 登录主站

访问 `https://puaai.xyz`，登录并确认账户状态可用。

如果你还不清楚充值、站内额度和倍率的关系，先看 <a href="../reference/pricing.html">计费与价格说明</a>。

## 2. 创建 API Key

进入用户后台的 **API Keys** 页面，创建一把测试用 Key。

首次创建建议：

- 名称写清用途
- 选择一个可用分组
- 自定义 Key 留空时由系统生成
- 先设置额度限制和过期时间

更详细说明见 <a href="./token.html">API Keys 说明</a>。

## 3. 先查模型列表

不要先手填模型名，先确认这把 Key 当前真正可用的模型。

```bash
curl https://puaai.xyz/v1/models \
  -H "Authorization: Bearer sk-your-api-key"
```

Gemini CLI 这类特殊客户端，按对应接入页单独配置即可。

## 4. 开始接入

- 普通 API：看 <a href="../reference/endpoints.html">API 地址</a>
- Claude Code：看 <a href="../integration/claude-code.html">Claude Code 接入指南</a>
- Codex：看 <a href="../integration/codex.html">Codex 接入指南</a>
- OpenCode：看 <a href="../integration/opencode.html">OpenCode 接入指南</a>
