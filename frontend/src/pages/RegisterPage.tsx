import { useState, useEffect } from "react"
import { zodResolver } from "@hookform/resolvers/zod"
import { useForm } from "react-hook-form"
import * as z from "zod"
import { useTranslation } from "react-i18next"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { Form, FormControl, FormField, FormItem, FormLabel, FormMessage } from "@/components/ui/form"
import { Card, CardContent, CardDescription, CardFooter, CardHeader, CardTitle } from "@/components/ui/card"
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog"
import { Link } from "react-router-dom"

const APP_NAME = import.meta.env.VITE_APP_NAME || "Kaleidoscope"

export default function RegisterPage() {
  const { t } = useTranslation()
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState<string | null>(null)
  const [showSuccessDialog, setShowSuccessDialog] = useState(false)
  const [countdown, setCountdown] = useState(5)

  const registerSchema = z.object({
    username: z.string()
      .min(3, { message: t("register.usernameMin") })
      .max(20, { message: t("register.usernameMax") })
      .regex(/^[a-zA-Z0-9_]+$/, { message: t("register.usernameInvalid") }),
    email: z.string().email({ message: t("register.emailInvalid") }),
    password: z.string().min(6, { message: t("register.passwordMin") }),
  })

  type RegisterFormValues = z.infer<typeof registerSchema>

  const form = useForm<RegisterFormValues>({
    resolver: zodResolver(registerSchema),
    defaultValues: {
      username: "",
      email: "",
      password: "",
    },
  })

  useEffect(() => {
    if (!showSuccessDialog) return

    const timer = setInterval(() => {
      setCountdown((prev) => {
        if (prev <= 1) {
          clearInterval(timer)
          window.location.href = "/login"
          return 0
        }
        return prev - 1
      })
    }, 1000)

    return () => clearInterval(timer)
  }, [showSuccessDialog])

  const onSubmit = async (data: RegisterFormValues) => {
    setLoading(true)
    setError(null)
    
    try {
      const response = await fetch(`${import.meta.env.VITE_API_BASE_URL}/api/v1/users/register`, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify(data),
      })
      
      if (!response.ok) {
        const errorData = await response.json()
        throw new Error(errorData.message || t("register.registerFailed"))
      }
      
      setShowSuccessDialog(true)
    } catch (err) {
      const message = err instanceof Error ? err.message : t("register.registerError")
      setError(message)
    } finally {
      setLoading(false)
    }
  }

  return (
    <div className="flex min-h-screen items-center justify-center bg-background p-4">
      <Card className="w-full max-w-md">
        <CardHeader>
          <CardTitle className="text-2xl text-foreground">{APP_NAME} - {t("register.title")}</CardTitle>
          <CardDescription>{t("register.description")}</CardDescription>
        </CardHeader>
        <CardContent>
          <Form {...form}>
            <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-4">
              {error && (
                <div className="rounded-md bg-destructive/10 p-3 text-destructive">
                  {error}
                </div>
              )}
              <FormField
                control={form.control}
                name="username"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>{t("register.username")}</FormLabel>
                    <FormControl>
                      <Input placeholder={t("register.usernamePlaceholder")} disabled={loading} {...field} />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />
              <FormField
                control={form.control}
                name="email"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>{t("register.email")}</FormLabel>
                    <FormControl>
                      <Input placeholder={t("register.emailPlaceholder")} disabled={loading} {...field} />
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
                    <FormLabel>{t("register.password")}</FormLabel>
                    <FormControl>
                      <Input type="password" placeholder={t("register.passwordPlaceholder")} disabled={loading} {...field} />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />
              <Button type="submit" className="w-full" disabled={loading}>
                {loading ? t("register.submitting") : t("register.submit")}
              </Button>
            </form>
          </Form>
        </CardContent>
        <CardFooter className="flex justify-center">
          <Link to="/login" className="text-sm text-primary hover:underline">
            {t("register.goToLogin")}
          </Link>
        </CardFooter>
      </Card>
      <Dialog open={showSuccessDialog} onOpenChange={setShowSuccessDialog}>
        <DialogContent className="sm:max-w-[425px]">
          <DialogHeader>
            <DialogTitle>{t("register.successTitle")}</DialogTitle>
            <DialogDescription>
              {t("register.successDescription", { count: countdown })}
            </DialogDescription>
          </DialogHeader>
        </DialogContent>
      </Dialog>
    </div>
  )
}
