export const sidebar = {
  '/': [
    {
      text: '使用文档',
      icon: 'material-symbols:docs-rounded',
      expanded: true,
      collapsible: false,
      children: [
        {
          text: '快速开始',
          icon: 'material-symbols:bolt-rounded',
          link: '/guide/quick-start.html'
        },
        {
          text: 'API Keys 说明',
          icon: 'mdi:key-variant',
          link: '/guide/token.html'
        },
        {
          text: 'API 地址',
          icon: 'material-symbols:route-rounded',
          link: '/reference/endpoints.html'
        },
        {
          text: '计费与价格说明',
          icon: 'material-symbols:calculate-rounded',
          link: '/reference/pricing.html'
        }
      ]
    },
    {
      text: '快速接入',
      icon: 'material-symbols:terminal-rounded',
      expanded: true,
      collapsible: false,
      children: [
        {
          text: 'Claude Code 接入指南',
          icon: 'material-symbols:smart-toy-rounded',
          link: '/integration/claude-code.html'
        },
        {
          text: 'Codex 接入指南',
          icon: 'material-symbols:code-rounded',
          link: '/integration/codex.html'
        },
        {
          text: 'OpenCode 接入指南',
          icon: 'material-symbols:deployed-code-rounded',
          link: '/integration/opencode.html'
        }
      ]
    }
  ]
}
