---
home: true
title: puaai 文档
icon: material-symbols:docs-rounded
heroText: puaai 接入文档
tagline: 先看计费口径，再拿 Key、绑分组和发第一条请求。
heroImage: logo.png
bgImage: hero-grid.svg
heroImageStyle:
  width: 104px
  borderRadius: 28px
  boxShadow: 0 26px 60px rgba(15, 118, 110, 0.24)
actions:
  - text: 先看计费说明
    link: /reference/pricing.html
    type: primary
    icon: material-symbols:calculate-rounded
  - text: 先获取 API Key
    link: /guide/token.html
    type: default
    icon: mdi:key-variant
  - text: 发送第一个请求
    link: /guide/first-request.html
    type: default
    icon: material-symbols:rocket-launch-rounded
  - text: 查看接口参考
    link: /reference/endpoints.html
    type: default
    icon: material-symbols:api-rounded
features:
  - title: 先把价格看懂
    icon: material-symbols:calculate-rounded
    details: 先区分充值汇率、站内额度、分组倍率和输入输出价格，再决定怎么充值和选分组。
    link: /reference/pricing.html
  - title: 新用户先跑通
    icon: material-symbols:bolt-rounded
    details: 先创建 Key，再给 Key 绑定分组，再用模型列表自检，最后再去接客户端。
    link: /guide/token.html
  - title: 面向真实接入
    icon: material-symbols:terminal-rounded
    details: 文档默认围绕 Base URL、协议路径、最小请求和排错顺序来组织。
    link: /guide/first-request.html
highlights:
  - header: 30 秒接入顺序
    description: 首次接入先把价格、Key、分组和最小请求这四步跑通，剩下的细节按需再查。
    color: '#0f766e'
    type: order
    highlights:
      - title: 先看计费说明
        icon: material-symbols:calculate-rounded
        details: 先分清刀、汇率、倍率和实际扣费之间是什么关系。
      - title: 创建一把 API Key
        icon: mdi:key-variant
        details: 登录后台，先创建一把测试用 Key。
      - title: 给 Key 绑定可用分组
        icon: material-symbols:group-work-rounded
        details: 分组决定平台边界、模型可见性和计费规则。
      - title: 再发第一条正式请求
        icon: material-symbols:send-rounded
        details: Claude、OpenAI、Gemini 走的入口不同，先用最小请求验证协议没填错。
  - header: 你最常会查的内容
    description: 这套文档按“先接入，再排障，再查规则”组织，不按内部实现铺陈。
    color: '#14532d'
    highlights:
      - title: 计费与价格
        icon: material-symbols:calculate-rounded
        details: 把刀、汇率、倍率和实际扣费的关系讲清楚。
        link: /reference/pricing.html
      - title: 获取 API Key
        icon: mdi:key-outline
        details: 给第一次来的用户看最短路径。
        link: /guide/token.html
      - title: 第一个请求
        icon: material-symbols:terminal
        details: 直接复制最小 curl 样例做联通性验证。
        link: /guide/first-request.html
      - title: 分组说明
        icon: material-symbols:account-tree-rounded
        details: 明确分组和平台、模型、计费之间的关系。
        link: /reference/groups.html
      - title: 模型与自检
        icon: material-symbols:deployed-code-history-rounded
        details: 用当前 Key 查实时可用模型，不靠静态列表猜。
        link: /reference/models.html
      - title: 常见问题
        icon: material-symbols:help-center-rounded
        details: 报错时按接入路径排查，而不是零散搜答案。
        link: /reference/faq.html
---

::: important 第一次访问建议先看三件事
- 先看 [计费与价格说明](./reference/pricing.md)，分清充值汇率、站内额度和倍率
- API Base URL 填 `https://puaai.xyz`
- 文档站 `https://puaai.xyz/docs` 只负责说明，不是网关入口
:::

<div class="docs-command-bar">
  <div class="docs-command-pill">
    <span class="docs-command-label">Pricing</span>
    <span>先读计费说明</span>
  </div>
  <div class="docs-command-pill">
    <span class="docs-command-label">Base URL</span>
    <code>https://puaai.xyz</code>
  </div>
  <div class="docs-command-pill">
    <span class="docs-command-label">Docs</span>
    <code>https://puaai.xyz/docs</code>
  </div>
  <div class="docs-command-pill">
    <span class="docs-command-label">Search</span>
    <span>Ctrl + K</span>
  </div>
</div>

## 一眼看懂

<div class="docs-home-grid">
  <a class="docs-home-card" href="./reference/pricing.html">
    <p class="docs-home-kicker">Step 01</p>
    <h3>先把价格口径看明白</h3>
    <p>先理解刀、充值汇率、分组倍率和输入输出价格，再决定怎么充值和选分组。</p>
  </a>
  <a class="docs-home-card" href="./guide/token.html">
    <p class="docs-home-kicker">Step 02</p>
    <h3>先拿到能用的 Key</h3>
    <p>进入用户后台创建 API Key，并确认它已经绑定一个有效分组。</p>
  </a>
  <a class="docs-home-card" href="./guide/first-request.html">
    <p class="docs-home-kicker">Step 03</p>
    <h3>先跑最小请求</h3>
    <p>用 curl 验证 Base URL、鉴权和协议入口都正确，再去接客户端或业务代码。</p>
  </a>
</div>

<div class="docs-split-panel">
  <div class="docs-split-copy">
    <p class="docs-home-kicker">How To Read</p>
    <h2>先看价格，再决定怎么用</h2>
    <p>这套文档把计费放到了最前面。第一次来的用户先看价格口径，再去拿 Key、绑分组和发请求；已经在用的用户直接查分组、模型和 FAQ。</p>
    <div class="docs-mini-list">
      <div class="docs-mini-item">
        <strong>新用户入口</strong>
        <span>计费说明 → 获取 API Key → 第一条请求</span>
      </div>
      <div class="docs-mini-item">
        <strong>老用户入口</strong>
        <span>分组说明 → 计费说明 → 模型与自检 → FAQ</span>
      </div>
      <div class="docs-mini-item">
        <strong>移动端阅读</strong>
        <span>顶部导航、左侧栏和目录都做了折叠与卡片化处理。</span>
      </div>
    </div>
  </div>
  <div class="docs-split-art">
    <img src="/flow-overview.svg" alt="puaai docs quickstart overview" />
  </div>
</div>

## 你真正会用到的地址

| 名称 | 地址 | 用途 |
| --- | --- | --- |
| 文档站 | `https://puaai.xyz/docs` | 查看接入说明和参考文档 |
| 平台首页 | `https://puaai.xyz` | 注册、登录、进入用户后台 |
| API Base URL | `https://puaai.xyz` | 绝大多数客户端应该填写这个值 |
| Key 用量查询 | `https://puaai.xyz/key-usage` | 公开查询某把 Key 的余额和用量 |

::: tip 为什么文档一直强调“先查模型列表”
puaai 的模型可见性不是固定表，而是由 Key 当前绑定的分组和该分组能调度到的上游共同决定。先查模型列表，比在客户端里盲填模型名更稳。
:::

::: tip 为什么这次把计费放到最前面
对第一次接触中转站的用户来说，最容易混淆的不是接口怎么调，而是“刀”“倍率”“缓存计费”和实际人民币成本是什么关系。先把这层看明白，后面的充值、选分组和接入动作才不容易误判。
:::

<div class="docs-stat-strip">
  <div class="docs-stat-card">
    <span class="docs-stat-value">1</span>
    <span class="docs-stat-label">统一 Base URL</span>
  </div>
  <div class="docs-stat-card">
    <span class="docs-stat-value">4</span>
    <span class="docs-stat-label">主要兼容层</span>
  </div>
  <div class="docs-stat-card">
    <span class="docs-stat-value">Ctrl + K</span>
    <span class="docs-stat-label">站内搜索</span>
  </div>
</div>

## 常用入口速查

<div class="docs-endpoint-grid">
  <div class="docs-endpoint-card">
    <code>/v1/messages</code>
    <p>Claude / Anthropic 兼容入口</p>
  </div>
  <div class="docs-endpoint-card">
    <code>/v1/chat/completions</code>
    <p>OpenAI Chat Completions 兼容入口</p>
  </div>
  <div class="docs-endpoint-card">
    <code>/v1/responses</code>
    <p>OpenAI Responses 兼容入口</p>
  </div>
  <div class="docs-endpoint-card">
    <code>/v1/models</code>
    <p>Claude / OpenAI 兼容分组的模型列表</p>
  </div>
  <div class="docs-endpoint-card">
    <code>/v1beta/models</code>
    <p>Gemini 原生兼容层的模型列表</p>
  </div>
  <div class="docs-endpoint-card">
    <code>/sora/v1/chat/completions</code>
    <p>Sora 专用入口，和普通文本模型分开理解</p>
  </div>
</div>

<div class="docs-split-panel docs-split-panel-reverse">
  <div class="docs-split-art">
    <img src="/reference-map.svg" alt="puaai reference map" />
  </div>
  <div class="docs-split-copy">
    <p class="docs-home-kicker">Reference</p>
    <h2>需要查的时候，定位要快</h2>
    <p>左侧导航和页面标题故意做得更短、更像产品文档。你应该能在几秒内区分“获取 Key”“协议入口”“分组规则”“模型自检”和“FAQ”分别在哪一层。</p>
    <div class="docs-mini-list">
      <div class="docs-mini-item">
        <strong>导航层级更短</strong>
        <span>把原本泛泛的文案收成用户会直接点击的词。</span>
      </div>
      <div class="docs-mini-item">
        <strong>提示块有优先级</strong>
        <span>`important` 讲硬规则，`tip` 讲最佳实践，`warning` 讲容易踩坑的点。</span>
      </div>
      <div class="docs-mini-item">
        <strong>目录更像索引</strong>
        <span>右侧 TOC 和页面图标让阅读时更容易快速扫读。</span>
      </div>
    </div>
  </div>
</div>

## 按目标阅读

<div class="docs-home-grid docs-home-grid-compact">
  <a class="docs-home-card" href="./reference/pricing.html">
    <p class="docs-home-kicker">先看</p>
    <h3>怎么理解计费价格</h3>
    <p>先看充值汇率、站内额度、分组倍率和输入输出价格怎么对应。</p>
  </a>
  <a class="docs-home-card" href="./guide/token.html">
    <p class="docs-home-kicker">新用户</p>
    <h3>怎么拿令牌</h3>
    <p>直接看创建 Key、分组绑定和首次自检的最短路径。</p>
  </a>
  <a class="docs-home-card" href="./guide/first-request.html">
    <p class="docs-home-kicker">接入</p>
    <h3>怎么发第一条请求</h3>
    <p>按 Claude、OpenAI、Gemini 三类入口复制最小样例。</p>
  </a>
  <a class="docs-home-card" href="./reference/groups.html">
    <p class="docs-home-kicker">规则</p>
    <h3>怎么理解分组</h3>
    <p>看平台边界、倍率、专属限制和订阅模式到底影响什么。</p>
  </a>
  <a class="docs-home-card" href="./reference/pricing.html">
    <p class="docs-home-kicker">计费</p>
    <h3>怎么理解刀和倍率</h3>
    <p>把充值汇率、站内额度、分组倍率和输入输出价格放在同一页说明。</p>
  </a>
  <a class="docs-home-card" href="./reference/models.html">
    <p class="docs-home-kicker">自检</p>
    <h3>怎么查当前可用模型</h3>
    <p>不要维护静态模型表，用当前 Key 直接请求模型列表。</p>
  </a>
  <a class="docs-home-card" href="./reference/endpoints.html">
    <p class="docs-home-kicker">参考</p>
    <h3>接口和 Base URL 填什么</h3>
    <p>这里集中列出 `/v1`、`/v1beta`、`/antigravity` 和 `/sora` 路径。</p>
  </a>
  <a class="docs-home-card" href="./reference/faq.html">
    <p class="docs-home-kicker">排障</p>
    <h3>报错时先查什么</h3>
    <p>从 Key、分组、模型名和协议入口四条线排，不要先怀疑客户端。</p>
  </a>
</div>
