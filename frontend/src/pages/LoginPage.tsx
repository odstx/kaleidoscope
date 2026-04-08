import { useState, useEffect } from "react"
import { zodResolver } from "@hookform/resolvers/zod"
import { useForm } from "react-hook-form"
import * as z from "zod"
import { useTranslation } from "react-i18next"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { Form, FormControl, FormField, FormItem, FormLabel, FormMessage } from "@/components/ui/form"
import { Label } from "@/components/ui/label"
import { Card, CardContent, CardDescription, CardFooter, CardHeader, CardTitle } from "@/components/ui/card"
import { Link, useNavigate } from "react-router-dom"
import { useAuth } from "@/contexts/AuthContext"
import { fetchFrontendConfig, getOidcAuthUrl } from "@/utils/oidc"

const APP_NAME = import.meta.env.VITE_APP_NAME || "Kaleidoscope"

export default function LoginPage() {
  const { t } = useTranslation()
  const navigate = useNavigate()
  const { login } = useAuth()
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState<string | null>(null)
  const [totpRequired, setTotpRequired] = useState(false)
  const [totpCode, setTotpCode] = useState("")
  const [pendingCredentials, setPendingCredentials] = useState<{ email: string; password: string } | null>(null)
  const [oidcEnabled, setOidcEnabled] = useState(false)

  useEffect(() => {
    console.log("Fetching frontend config...")
    fetchFrontendConfig().then((config) => {
      console.log("Config fetched:", config)
      setOidcEnabled(config.enabled && config.issuerUrl && config.clientId)
    }).catch((err) => {
      console.error("Config fetch error:", err)
    })
  }, [])

  const handleOidcLogin = async () => {
    try {
      const url = await getOidcAuthUrl()
      window.location.href = url
    } catch (err) {
      const message = err instanceof Error ? err.message : 'Failed to initiate OIDC login'
      setError(message)
    }
  }

  const loginSchema = z.object({
    email: z.string().email({ message: t('login.emailInvalid') }),
    password: z.string().min(1, { message: t('login.passwordRequired') }),
  })

  type LoginFormValues = z.infer<typeof loginSchema>

  const form = useForm<LoginFormValues>({
    resolver: zodResolver(loginSchema),
    defaultValues: {
      email: "",
      password: "",
    },
  })

  const handleLogin = async (email: string, password: string, totp?: string) => {
    const response = await fetch(`${import.meta.env.VITE_API_BASE_URL}/api/v1/users/login`, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({ email, password, totp_code: totp }),
    })

    if (!response.ok) {
      const errorData = await response.json()
      if (errorData.totp_required) {
        setTotpRequired(true)
        setPendingCredentials({ email, password })
        return
      }
      throw new Error(errorData.error || t('login.loginFailed'))
    }

    const result = await response.json()
    login(result.token)
    navigate("/dashboard")
  }

  const onSubmit = async (data: LoginFormValues) => {
    setLoading(true)
    setError(null)

    try {
      await handleLogin(data.email, data.password)
    } catch (err) {
      const message = err instanceof Error ? err.message : t('login.loginError')
      setError(message)
    } finally {
      setLoading(false)
    }
  }

  const onSubmitTotp = async () => {
    if (!pendingCredentials || !totpCode) return

    setLoading(true)
    setError(null)

    try {
      await handleLogin(pendingCredentials.email, pendingCredentials.password, totpCode)
    } catch (err) {
      const message = err instanceof Error ? err.message : t('login.loginError')
      setError(message)
    } finally {
      setLoading(false)
    }
  }

  return (
    <div className="flex min-h-screen items-center justify-center bg-background p-4">
      <Card className="w-full max-w-md">
        <CardHeader>
          <CardTitle className="text-2xl text-foreground">{APP_NAME} - {t('login.title')}</CardTitle>
          <CardDescription>{t('login.description')}</CardDescription>
        </CardHeader>
        <CardContent>
          {!totpRequired ? (
            <>
              <Form {...form}>
                <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-4">
                  {error && (
                    <div className="rounded-md bg-destructive/10 p-3 text-destructive">
                      {error}
                    </div>
                  )}
                  <FormField
                    control={form.control}
                    name="email"
                    render={({ field }) => (
                      <FormItem>
                        <FormLabel>{t('login.email')}</FormLabel>
                        <FormControl>
                          <Input placeholder={t('login.emailPlaceholder')} {...field} />
                        </FormControl>
                        <FormMessage />
                      </FormItem>
                    )}
                  />
                  <FormField
                    control={form.control}
                    name="password"
                    render={({ field }) => (
                      <FormItem>
                        <FormLabel>{t('login.password')}</FormLabel>
                        <FormControl>
                          <Input type="password" placeholder={t('login.passwordPlaceholder')} {...field} />
                        </FormControl>
                        <FormMessage />
                      </FormItem>
                    )}
                  />
                  <Button type="submit" className="w-full" disabled={loading}>
                    {loading ? t('login.submitting') : t('login.submit')}
                  </Button>
                </form>
              </Form>
              {oidcEnabled && (
                <>
                  <div className="relative my-4">
                    <div className="absolute inset-0 flex items-center">
                      <span className="w-full border-t" />
                    </div>
                    <div className="relative flex justify-center text-xs uppercase">
                      <span className="bg-card px-2 text-muted-foreground">Or</span>
                    </div>
                  </div>
                  <Button
                    type="button"
                    variant="outline"
                    className="w-full"
                    onClick={handleOidcLogin}
                    disabled={loading}
                  >
                    {t('login.continueWithOidc')}
                  </Button>
                </>
              )}
            </>
          ) : (
            <div className="space-y-4">
              {error && (
                <div className="rounded-md bg-destructive/10 p-3 text-destructive">
                  {error}
                </div>
              )}
              <div className="space-y-2">
                <Label>{t('login.totpCode')}</Label>
                <Input
                  placeholder={t('login.totpPlaceholder')}
                  value={totpCode}
                  onChange={(e) => setTotpCode(e.target.value)}
                  maxLength={6}
                />
              </div>
              <Button onClick={onSubmitTotp} className="w-full" disabled={loading || totpCode.length !== 6}>
                {loading ? t('login.submitting') : t('login.submit')}
              </Button>
            </div>
          )}
        </CardContent>
        <CardFooter className="flex flex-col gap-2">
          <div className="text-center w-full">
            <Link to="/forgot-password" className="text-sm text-primary hover:underline">
              {t('login.forgotPassword')}
            </Link>
          </div>
          <div className="text-center w-full">
            <Link to="/register" className="text-sm text-primary hover:underline">
              {t('login.goToRegister')}
            </Link>
          </div>
        </CardFooter>
      </Card>
    </div>
  )
}