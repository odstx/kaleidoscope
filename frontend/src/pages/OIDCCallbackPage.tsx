import { useEffect, useState, useRef } from "react"
import { useNavigate, useSearchParams, Link } from "react-router-dom"
import { useTranslation } from "react-i18next"
import { useAuth } from "@/contexts/AuthContext"
import { handleOidcCallback } from "@/utils/oidc"
import { Card, CardContent, CardDescription, CardTitle } from "@/components/ui/card"

export default function OIDCCallbackPage() {
  const { t } = useTranslation()
  const [searchParams] = useSearchParams()
  const navigate = useNavigate()
  const { login } = useAuth()
  const [status, setStatus] = useState<"loading" | "error" | "success">("loading")
  const [errorMessage, setErrorMessage] = useState<string | null>(null)
  const processed = useRef(false)

  useEffect(() => {
    if (processed.current) return
    processed.current = true

    const code = searchParams.get("code")
    if (!code) {
      setTimeout(() => {
        setErrorMessage("No authorization code received")
        setStatus("error")
      }, 0)
      return
    }

    const executeCallback = async () => {
      try {
        const result = await handleOidcCallback(code)
        login(result.token)
        setTimeout(() => {
          setStatus("success")
          navigate("/dashboard")
        }, 0)
      } catch (err) {
        const message = err instanceof Error ? err.message : "OIDC authentication failed"
        setTimeout(() => {
          setErrorMessage(message)
          setStatus("error")
        }, 0)
      }
    }

    executeCallback()
  }, [searchParams, login, navigate])

  if (status === "error") {
    return (
      <div className="flex min-h-screen items-center justify-center bg-background p-4">
        <Card className="w-full max-w-md">
          <CardHeader>
            <CardTitle className="text-2xl text-destructive">{t('oidc.callback.errorTitle')}</CardTitle>
            <CardDescription>{errorMessage}</CardDescription>
          </CardHeader>
          <CardContent>
            <Link to="/login" className="text-primary hover:underline">
              {t('oidc.callback.returnToLogin')}
            </Link>
          </CardContent>
        </Card>
      </div>
    )
  }

  return (
    <div className="flex min-h-screen items-center justify-center bg-background p-4">
      <Card className="w-full max-w-md">
        <CardHeader>
          <CardTitle className="text-2xl text-foreground">{t('oidc.callback.title')}</CardTitle>
          <CardDescription>{t('oidc.callback.description')}</CardDescription>
        </CardHeader>
      </Card>
    </div>
  )
}