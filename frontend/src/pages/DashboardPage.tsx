import { useNavigate } from "react-router-dom"
import { useTranslation } from "react-i18next"
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card"
import { useAuth } from "@/contexts/AuthContext"

export default function DashboardPage() {
  const navigate = useNavigate()
  const { isAuthenticated } = useAuth()
  const { t } = useTranslation()

  if (!isAuthenticated) {
    navigate("/login")
    return null
  }

  return (
    <div className="min-h-screen bg-background p-8">
      <div className="max-w-4xl mx-auto space-y-8">
        <h1 className="text-3xl font-bold text-foreground">{t("dashboard.title")}</h1>

        <Card>
          <CardHeader>
            <CardTitle>{t("dashboard.welcome")}</CardTitle>
            <CardDescription>{t("dashboard.welcomeDesc")}</CardDescription>
          </CardHeader>
          <CardContent>
            <p className="text-muted-foreground">
              {t("dashboard.welcomeMessage")}
            </p>
          </CardContent>
        </Card>
      </div>
    </div>
  )
}
