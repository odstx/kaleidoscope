import { useAuth } from '@/contexts/AuthContext'
import { useTheme } from '@/contexts/ThemeContext'
import { useTranslation } from 'react-i18next'
import { Link } from 'react-router-dom'
import { Button } from '@/components/ui/button'
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
  DropdownMenuSub,
  DropdownMenuSubContent,
  DropdownMenuSubTrigger,
  DropdownMenuSeparator,
} from '@/components/ui/dropdown-menu'

export function Navbar() {
  const { isAuthenticated, logout } = useAuth()
  const { t, i18n } = useTranslation()
  const { theme, setTheme } = useTheme()

  if (!isAuthenticated) return null

  const changeLanguage = (lng: string) => {
    i18n.changeLanguage(lng)
  }

  const getThemeLabel = () => {
    switch (theme) {
      case 'light':
        return t('theme.light')
      case 'dark':
        return t('theme.dark')
      default:
        return t('theme.auto')
    }
  }

  return (
    <nav className="bg-background border-b border-border px-4 py-3">
      <div className="max-w-7xl mx-auto flex items-center justify-between">
        <Link to="/dashboard" className="text-xl font-semibold text-foreground">
          {import.meta.env.VITE_APP_NAME}
        </Link>
        <div className="flex items-center gap-2">
          <DropdownMenu>
            <DropdownMenuTrigger asChild>
              <Button variant="ghost">{t("nav.menu")}</Button>
            </DropdownMenuTrigger>
            <DropdownMenuContent align="end">
              <DropdownMenuItem asChild>
                <Link to="/profile" className="text-foreground">{t("nav.profile")}</Link>
              </DropdownMenuItem>
              <DropdownMenuSub>
                <DropdownMenuSubTrigger>
                  {t('language.title')} ({i18n.language === 'zh' ? '中文' : 'EN'})
                </DropdownMenuSubTrigger>
                <DropdownMenuSubContent>
                  <DropdownMenuItem onClick={() => changeLanguage('zh')}>
                    {t('language.zh')}
                  </DropdownMenuItem>
                  <DropdownMenuItem onClick={() => changeLanguage('en')}>
                    {t('language.en')}
                  </DropdownMenuItem>
                </DropdownMenuSubContent>
              </DropdownMenuSub>
              <DropdownMenuSub>
                <DropdownMenuSubTrigger>
                  {t('theme.title')} ({getThemeLabel()})
                </DropdownMenuSubTrigger>
                <DropdownMenuSubContent>
                  <DropdownMenuItem onClick={() => setTheme('auto')}>
                    {t('theme.auto')}
                  </DropdownMenuItem>
                  <DropdownMenuItem onClick={() => setTheme('light')}>
                    {t('theme.light')}
                  </DropdownMenuItem>
                  <DropdownMenuItem onClick={() => setTheme('dark')}>
                    {t('theme.dark')}
                  </DropdownMenuItem>
                </DropdownMenuSubContent>
              </DropdownMenuSub>
              <DropdownMenuSeparator />
              <DropdownMenuItem onClick={logout}>
                {t("nav.logout")}
              </DropdownMenuItem>
            </DropdownMenuContent>
          </DropdownMenu>
        </div>
      </div>
    </nav>
  )
}
