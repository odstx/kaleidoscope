import { useEffect, useState } from "react"
import { useNavigate } from "react-router-dom"
import { useTranslation } from "react-i18next"
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card"
import { useAuth } from "@/contexts/AuthContext"

interface UserInfo {
  email: string
  id?: number
}

export default function UserProfilePage() {
  const navigate = useNavigate()
  const { isAuthenticated } = useAuth()
  const { t } = useTranslation()
  const [user, setUser] = useState<UserInfo | null>(null)
  const [loading, setLoading] = useState(true)

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
              {user?.id && (
                <div className="flex justify-between border-b pb-2">
                  <span className="text-muted-foreground">{t("profile.userId")}</span>
                  <span className="font-medium">{user.id}</span>
                </div>
              )}
            </div>
          </CardContent>
        </Card>
      </div>
    </div>
  )
}
