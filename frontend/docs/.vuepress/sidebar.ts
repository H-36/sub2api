export const sidebar = {
  '/': [
    {
      text: '开始使用',
      icon: 'material-symbols:rocket-launch-rounded',
      link: '/',
      expanded: true,
      collapsible: false,
      children: [
        {
          text: '概览',
          icon: 'material-symbols:docs-rounded',
          link: '/'
        },
        {
          text: '计费与价格',
          icon: 'material-symbols:calculate-rounded',
          link: '/reference/pricing.html'
        },
        {
          text: '获取 API Key',
          icon: 'mdi:key-variant',
          link: '/guide/token.html'
        },
        {
          text: '第一个请求',
          icon: 'material-symbols:terminal-rounded',
          link: '/guide/first-request.html'
        }
      ]
    },
    {
      text: '接入与价格',
      icon: 'material-symbols:api-rounded',
      expanded: true,
      collapsible: false,
      children: [
        {
          text: '接口与 Base URL',
          icon: 'material-symbols:route-rounded',
          link: '/reference/endpoints.html'
        },
        {
          text: '分组说明',
          icon: 'material-symbols:group-work-rounded',
          link: '/reference/groups.html'
        },
        {
          text: '模型与自检',
          icon: 'material-symbols:view-list-rounded',
          link: '/reference/models.html'
        }
      ]
    },
    {
      text: '排障与查询',
      icon: 'material-symbols:help-center-rounded',
      collapsible: false,
      children: [
        {
          text: '常见问题',
          icon: 'material-symbols:quiz-rounded',
          link: '/reference/faq.html'
        }
      ]
    }
  ]
}
