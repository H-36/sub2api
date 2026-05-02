export const INTERRUPTED_TASK_ERROR =
  '请求中断：页面离开或刷新导致本地请求停止，生成结果无法继续接收。'

export function isInterruptedTaskError(error: string | null | undefined): boolean {
  const text = (error || '').trim()
  return text === INTERRUPTED_TASK_ERROR || text === '请求中断' || text.includes('请求中断')
}
