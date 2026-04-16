export type ThemeMode = 'dark' | 'light'

const THEME_STORAGE_KEY = 'theme'

export function getThemeMode(): ThemeMode {
  return localStorage.getItem(THEME_STORAGE_KEY) === 'light' ? 'light' : 'dark'
}

export function applyTheme(mode: ThemeMode): boolean {
  const isDark = mode === 'dark'
  document.documentElement.classList.toggle('dark', isDark)
  return isDark
}

export function initThemeMode(): boolean {
  return applyTheme(getThemeMode())
}

export function toggleThemeMode(currentIsDark: boolean): boolean {
  const nextMode: ThemeMode = currentIsDark ? 'light' : 'dark'
  localStorage.setItem(THEME_STORAGE_KEY, nextMode)
  return applyTheme(nextMode)
}
