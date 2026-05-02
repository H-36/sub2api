import { docsMeta } from './docsMeta'

export const navbar = [
  {
    text: '首页',
    icon: 'material-symbols:rocket-launch-rounded',
    link: '/'
  },
  {
    text: '快速开始',
    icon: 'material-symbols:bolt-rounded',
    link: '/guide/quick-start.html'
  },
  {
    text: 'CLI 配置',
    icon: 'material-symbols:terminal-rounded',
    link: '/cli/overview.html'
  },
  {
    text: '第三方调用',
    icon: 'material-symbols:extension-rounded',
    link: '/third-party/overview.html'
  },
  {
    text: '返回主站',
    icon: 'material-symbols:public-rounded',
    link: docsMeta.siteUrl
  }
]
