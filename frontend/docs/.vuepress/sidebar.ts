export const sidebar = {
  '/': [
    {
      text: '快速开始',
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
          text: '第一个请求',
          icon: 'material-symbols:terminal-rounded',
          link: '/guide/first-request.html'
        },
        {
          text: '分组与模型',
          icon: 'material-symbols:hub-rounded',
          link: '/guide/groups-and-models.html'
        }
      ]
    },
    {
      text: 'CLI 配置',
      icon: 'material-symbols:terminal-rounded',
      expanded: true,
      collapsible: false,
      children: [
        {
          text: '通用前置步骤',
          icon: 'material-symbols:checklist-rounded',
          link: '/cli/overview.html'
        },
        {
          text: 'Claude Code',
          icon: 'material-symbols:smart-toy-rounded',
          link: '/integration/claude-code.html'
        },
        {
          text: 'Codex',
          icon: 'material-symbols:code-rounded',
          link: '/integration/codex.html'
        },
        {
          text: 'Gemini CLI',
          icon: 'material-symbols:diamond-rounded',
          link: '/cli/gemini.html'
        },
        {
          text: 'Windows / WSL',
          icon: 'material-symbols:desktop-windows-rounded',
          link: '/cli/wsl.html'
        }
      ]
    },
    {
      text: '第三方调用',
      icon: 'material-symbols:extension-rounded',
      expanded: true,
      collapsible: false,
      children: [
        {
          text: '总览',
          icon: 'material-symbols:apps-rounded',
          link: '/third-party/overview.html'
        },
        {
          text: 'Curl / API 测试',
          icon: 'material-symbols:http-rounded',
          link: '/third-party/curl.html'
        },
        {
          text: 'OpenCode',
          icon: 'material-symbols:deployed-code-rounded',
          link: '/integration/opencode.html'
        },
        {
          text: 'Kilo Code',
          icon: 'material-symbols:code-blocks-rounded',
          link: '/third-party/kilo-code.html'
        },
        {
          text: 'Zed',
          icon: 'material-symbols:edit-square-rounded',
          link: '/third-party/zed.html'
        },
        {
          text: 'Cline / Roo Code',
          icon: 'material-symbols:developer-mode-rounded',
          link: '/third-party/cline-roo.html'
        },
        {
          text: 'Open WebUI / LobeChat',
          icon: 'material-symbols:chat-rounded',
          link: '/third-party/webui-lobechat.html'
        },
        {
          text: 'Hermes Agent',
          icon: 'material-symbols:automation-rounded',
          link: '/third-party/hermes.html'
        }
      ]
    },
    {
      text: '参考资料',
      icon: 'material-symbols:menu-book-rounded',
      expanded: true,
      collapsible: false,
      children: [
        {
          text: 'API 地址',
          icon: 'material-symbols:route-rounded',
          link: '/reference/endpoints.html'
        },
        {
          text: '模型与自检',
          icon: 'material-symbols:view-list-rounded',
          link: '/reference/models.html'
        },
        {
          text: '分组说明',
          icon: 'material-symbols:group-work-rounded',
          link: '/reference/groups.html'
        },
        {
          text: '计费与价格说明',
          icon: 'material-symbols:calculate-rounded',
          link: '/reference/pricing.html'
        },
        {
          text: '生图接口',
          icon: 'material-symbols:image-rounded',
          link: '/integration/image-generation.html'
        },
        {
          text: '常见问题',
          icon: 'material-symbols:help-rounded',
          link: '/reference/faq.html'
        }
      ]
    }
  ]
}
