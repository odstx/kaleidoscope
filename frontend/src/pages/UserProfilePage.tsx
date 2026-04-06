import { useEffect, useState } from "react"
import { useNavigate } from "react-router-dom"
import { useTranslation } from "react-i18next"
import { QRCodeSVG } from "qrcode.react"
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { useAuth } from "@/contexts/AuthContext"

interface UserInfo {
  email: string
  id?: number
  uid?: string
  totp_enabled?: boolean
}

export default function UserProfilePage() {
  const navigate = useNavigate()
  const { isAuthenticated } = useAuth()
  const { t } = useTranslation()
  const [user, setUser] = useState<UserInfo | null>(null)
  const [loading, setLoading] = useState(true)
  const [totpLoading, setTotpLoading] = useState(false)
  const [totpSecret, setTotpSecret] = useState<string | null>(null)
  const [totpUrl, setTotpUrl] = useState<string | null>(null)
  const [totpCode, setTotpCode] = useState("")
  const [totpError, setTotpError] = useState<string | null>(null)

  useEffect(() => {
    if (!isAuthenticated) {
      navigate("/login")
      return
    }

    const token = localStorage.getItem("token")
    if (!token) {
      navigate("/login")
      return
    }

    const fetchUserInfo = async () => {
      try {
        const response = await fetch(`${import.meta.env.VITE_API_BASE_URL}/api/v1/users/info`, {
          method: "GET",
          headers: {
            "Content-Type": "application/json",
            "Authorization": `Bearer ${token}`,
          },
        })

        if (!response.ok) {
          throw new Error(t("profile.fetchError"))
        }

        const userData = await response.json()
        setUser(userData)
      } catch (error) {
        console.error(t("profile.fetchError"), error)
        localStorage.removeItem("token")
        navigate("/login")
      } finally {
        setLoading(false)
      }
    }

    fetchUserInfo()
  }, [navigate, isAuthenticated, t])

  const handleSetupTOTP = async () => {
    const token = localStorage.getItem("token")
    if (!token) return

    setTotpLoading(true)
    setTotpError(null)
    try {
      const response = await fetch(`${import.meta.env.VITE_API_BASE_URL}/api/v1/users/totp/setup`, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
          "Authorization": `Bearer ${token}`,
        },
      })

      if (!response.ok) {
        throw new Error(t("profile.totp.setupError"))
      }

      const data = await response.json()
      setTotpSecret(data.secret)
      setTotpUrl(data.url)
    } catch (error) {
      setTotpError(error instanceof Error ? error.message : t("profile.totp.setupError"))
    } finally {
      setTotpLoading(false)
    }
  }

  const handleVerifyTOTP = async () => {
    const token = localStorage.getItem("token")
    if (!token || !totpCode) return

    setTotpLoading(true)
    setTotpError(null)
    try {
      const response = await fetch(`${import.meta.env.VITE_API_BASE_URL}/api/v1/users/totp/verify`, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
          "Authorization": `Bearer ${token}`,
        },
        body: JSON.stringify({ code: totpCode }),
      })

      if (!response.ok) {
        const data = await response.json()
        throw new Error(data.error || t("profile.totp.verifyError"))
      }

      const enableResponse = await fetch(`${import.meta.env.VITE_API_BASE_URL}/api/v1/users/totp/enable`, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
          "Authorization": `Bearer ${token}`,
        },
      })

      if (!enableResponse.ok) {
        throw new Error(t("profile.totp.enableError"))
      }

      setTotpSecret(null)
      setTotpUrl(null)
      setTotpCode("")
      setUser(prev => prev ? { ...prev, totp_enabled: true } : null)
    } catch (error) {
      setTotpError(error instanceof Error ? error.message : t("profile.totp.verifyError"))
    } finally {
      setTotpLoading(false)
    }
  }

  const handleDisableTOTP = async () => {
    const token = localStorage.getItem("token")
    if (!token) return

    setTotpLoading(true)
    setTotpError(null)
    try {
      const response = await fetch(`${import.meta.env.VITE_API_BASE_URL}/api/v1/users/totp/disable`, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
          "Authorization": `Bearer ${token}`,
        },
      })

      if (!response.ok) {
        throw new Error(t("profile.totp.disableError"))
      }

      setUser(prev => prev ? { ...prev, totp_enabled: false } : null)
    } catch (error) {
      setTotpError(error instanceof Error ? error.message : t("profile.totp.disableError"))
    } finally {
      setTotpLoading(false)
    }
  }

  if (loading) {
    return (
      <div className="flex min-h-screen items-center justify-center bg-background">
        <div className="text-lg">{t("profile.loading")}</div>
      </div>
    )
  }

  return (
    <div className="min-h-screen bg-background p-8">
      <div className="max-w-4xl mx-auto space-y-8">
        <h1 className="text-3xl font-bold">{t("profile.title")}</h1>

        <Card>
          <CardHeader>
            <CardTitle>{t("profile.accountDetails")}</CardTitle>
            <CardDescription>{t("profile.accountDetailsDesc")}</CardDescription>
          </CardHeader>
          <CardContent>
            <div className="space-y-4">
              <div className="flex justify-between border-b pb-2">
                <span className="text-muted-foreground">{t("profile.email")}</span>
                <span className="font-medium">{user?.email}</span>
              </div>
              {user?.uid && (
                <div className="flex justify-between border-b pb-2">
                  <span className="text-muted-foreground">{t("profile.uid")}</span>
                  <span className="font-medium font-mono">{user.uid}</span>
                </div>
              )}
              {user?.id && (
                <div className="flex justify-between border-b pb-2">
                  <span className="text-muted-foreground">{t("profile.userId")}</span>
                  <span className="font-medium">{user.id}</span>
                </div>
              )}
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardHeader>
            <CardTitle>{t("profile.totp.title")}</CardTitle>
            <CardDescription>{t("profile.totp.description")}</CardDescription>
          </CardHeader>
          <CardContent>
            {totpError && (
              <div className="mb-4 rounded-md bg-destructive/10 p-3 text-destructive">
                {totpError}
              </div>
            )}
            
            {user?.totp_enabled ? (
              <div className="space-y-4">
                <div className="flex items-center justify-between">
                  <span className="text-sm text-muted-foreground">{t("profile.totp.enabled")}</span>
                  <Button variant="destructive" onClick={handleDisableTOTP} disabled={totpLoading}>
                    {totpLoading ? t("common.loading") : t("profile.totp.disable")}
                  </Button>
                </div>
              </div>
            ) : totpSecret ? (
              <div className="space-y-4">
                <div className="space-y-2 flex flex-col items-center">
                  <p className="text-sm text-muted-foreground">{t("profile.totp.scanQR")}</p>
                  <div className="p-4 bg-white rounded-lg">
                    <QRCodeSVG value={totpUrl || ""} size={200} />
                  </div>
                </div>
                <div className="space-y-2">
                  <p className="text-sm text-muted-foreground">{t("profile.totp.secret")}</p>
                  <code className="block p-2 bg-muted rounded text-sm font-mono break-all">{totpSecret}</code>
                </div>
                <div className="space-y-2">
                  <p className="text-sm text-muted-foreground">{t("profile.totp.enterCode")}</p>
                  <div className="flex gap-2">
                    <Input
                      value={totpCode}
                      onChange={(e) => setTotpCode(e.target.value)}
                      placeholder={t("profile.totp.codePlaceholder")}
                      maxLength={6}
                    />
                    <Button onClick={handleVerifyTOTP} disabled={totpLoading || totpCode.length !== 6}>
                      {totpLoading ? t("common.loading") : t("profile.totp.verify")}
                    </Button>
                  </div>
                </div>
              </div>
            ) : (
              <Button onClick={handleSetupTOTP} disabled={totpLoading}>
                {totpLoading ? t("common.loading") : t("profile.totp.setup")}
              </Button>
            )}
          </CardContent>
        </Card>
      </div>
    </div>
  )
}
