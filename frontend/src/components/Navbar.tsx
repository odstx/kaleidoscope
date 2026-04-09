import { useAuth } from '@/contexts/AuthContext'
import { useTheme } from '@/contexts/ThemeContext'
import { useTranslation } from 'react-i18next'
import { Link } from 'react-router-dom'
import { useEffect, useState } from 'react'
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
import {
  Drawer,
  DrawerContent,
  DrawerHeader,
  DrawerTitle,
  DrawerTrigger,
} from '@/components/ui/drawer'
import { Menu } from 'lucide-react'
import { apiCall } from '@/utils/api'

interface App {
  id: number
  name: string
  description: string
  icon: string
  url: string
  enabled: boolean
  order: number
}

export function Navbar() {
  const { isAuthenticated, logout } = useAuth()
  const { t, i18n } = useTranslation()
  const { theme, setTheme } = useTheme()
  const [apps, setApps] = useState<App[]>([])
  const [open, setOpen] = useState(false)

  useEffect(() => {
    if (isAuthenticated) {
      apiCall<App[]>('/api/v1/apps', { method: 'GET' })
        .then(setApps)
        .catch(() => setApps([]))
    }
  }, [isAuthenticated])

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

  const handleAppClick = () => {
    setOpen(false)
  }

  return (
    <nav className="bg-background border-b border-border px-4 py-3">
      <div className="max-w-7xl mx-auto flex items-center justify-between">
        <div className="flex items-center gap-2">
          <Drawer open={open} onOpenChange={setOpen}>
            <DrawerTrigger asChild>
              <Button variant="ghost" size="icon">
                <Menu className="h-5 w-5" />
              </Button>
            </DrawerTrigger>
            <DrawerContent side="left" className="w-72">
              <DrawerHeader>
                <DrawerTitle>{t('nav.apps')}</DrawerTitle>
              </DrawerHeader>
              <div className="p-4">
                {apps.length === 0 ? (
                  <p className="text-muted-foreground text-sm">{t('nav.noApps')}</p>
                ) : (
                  <div className="space-y-2">
                    {apps
                      .sort((a, b) => a.order - b.order)
                      .filter((app) => app.enabled)
                      .map((app) => (
                        <a
                          key={app.id}
                          href={app.url}
                          onClick={handleAppClick}
                          className="flex items-center gap-3 p-2 rounded-md hover:bg-accent transition-colors"
                        >
                          {app.icon ? (
                            <img src={app.icon} alt={app.name} className="w-8 h-8" />
                          ) : (
                            <div className="w-8 h-8 bg-primary/10 rounded-md flex items-center justify-center">
                              <span className="text-primary font-semibold">{app.name[0]}</span>
                            </div>
                          )}
                          <div className="flex-1 min-w-0">
                            <p className="font-medium truncate">{app.name}</p>
                            {app.description && (
                              <p className="text-xs text-muted-foreground truncate">{app.description}</p>
                            )}
                          </div>
                        </a>
                      ))}
                  </div>
                )}
              </div>
            </DrawerContent>
          </Drawer>
          <Link to="/dashboard" className="text-xl font-semibold text-foreground">
            {import.meta.env.VITE_APP_NAME}
          </Link>
        </div>
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
