---
title: puaai 文档
icon: material-symbols:docs-rounded
description: 文档
home: true
heroText: puaai 文档
tagline: 从 API Key 到 CLI、第三方工具和生产调用，一套路径跑通 puaai。
actions:
  - text: 快速开始
    link: /guide/quick-start.html
    type: primary
  - text: CLI 配置
    link: /cli/overview.html
  - text: 第三方调用
    link: /third-party/overview.html
---

# puaai 文档

欢迎使用 puaai 文档。puaai 提供 OpenAI、Anthropic、Gemini 等多协议兼容接入，建议先完成快速开始，再按你使用的工具选择对应接入页。

文档地址是 `https://puaai.xyz/docs`，实际 API Base URL 是 `https://puaai.xyz`。

## 快速导航

<div class="docs-home-grid">
  <a class="docs-home-card" href="./guide/quick-start.html">
    <p class="docs-home-kicker">START</p>
    <h3>快速开始</h3>
    <p>完成登录、创建 API Key、绑定分组、查询模型和第一个请求。</p>
  </a>
  <a class="docs-home-card" href="./cli/overview.html">
    <p class="docs-home-kicker">CLI</p>
    <h3>CLI 配置</h3>
    <p>配置 Claude Code、Codex、Gemini CLI，以及 Windows / WSL 环境。</p>
  </a>
  <a class="docs-home-card" href="./third-party/overview.html">
    <p class="docs-home-kicker">TOOLS</p>
    <h3>第三方调用</h3>
    <p>接入 OpenCode、Kilo Code、Zed、Cline、Open WebUI、LobeChat 等工具。</p>
  </a>
  <a class="docs-home-card" href="./third-party/curl.html">
    <p class="docs-home-kicker">TEST</p>
    <h3>Curl / API 测试</h3>
    <p>复制最小请求，先验证 Base URL、Key、模型和协议入口。</p>
  </a>
  <a class="docs-home-card" href="./reference/pricing.html">
    <p class="docs-home-kicker">BILLING</p>
    <h3>计费与价格</h3>
    <p>理解站内额度、充值比例、分组倍率、输入输出和缓存计费。</p>
  </a>
  <a class="docs-home-card" href="./reference/faq.html">
    <p class="docs-home-kicker">FAQ</p>
    <h3>常见问题</h3>
    <p>按 Key、分组、模型名、协议入口快速排查常见错误。</p>
  </a>
</div>

## 推荐路径

<div class="docs-step-list">
  <div class="docs-step-item">
    <span>1</span>
    <div>
      <strong>先创建测试 Key</strong>
      <p>进入主站 API Keys 页面，创建一把测试 Key，并选择状态正常的分组。</p>
    </div>
  </div>
  <div class="docs-step-item">
    <span>2</span>
    <div>
      <strong>再查询模型列表</strong>
      <p>不要凭记忆填写模型名，先用当前 Key 查询 `/v1/models`。</p>
    </div>
  </div>
  <div class="docs-step-item">
    <span>3</span>
    <div>
      <strong>最后配置工具</strong>
      <p>根据你的客户端选择 Claude Code、Codex、Gemini CLI 或第三方工具页面。</p>
    </div>
  </div>
</div>
