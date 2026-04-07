import { useState } from "react"
import { zodResolver } from "@hookform/resolvers/zod"
import { useForm } from "react-hook-form"
import * as z from "zod"
import { useTranslation } from "react-i18next"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { Form, FormControl, FormField, FormItem, FormLabel, FormMessage } from "@/components/ui/form"
import { Card, CardContent, CardDescription, CardFooter, CardHeader, CardTitle } from "@/components/ui/card"
import { Link } from "react-router-dom"

const APP_NAME = import.meta.env.VITE_APP_NAME || "Kaleidoscope"

export default function ForgotPasswordPage() {
  const { t } = useTranslation()
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState<string | null>(null)
  const [success, setSuccess] = useState(false)

  const forgotPasswordSchema = z.object({
    email: z.string().email({ message: t('forgotPassword.emailInvalid') }),
  })

  type ForgotPasswordFormValues = z.infer<typeof forgotPasswordSchema>

  const form = useForm<ForgotPasswordFormValues>({
    resolver: zodResolver(forgotPasswordSchema),
    defaultValues: {
      email: "",
    },
  })

  const onSubmit = async (data: ForgotPasswordFormValues) => {
    setLoading(true)
    setError(null)

    try {
      const response = await fetch(`${import.meta.env.VITE_API_BASE_URL}/api/v1/users/forgot-password`, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({ email: data.email }),
      })

      if (!response.ok) {
        const errorData = await response.json()
        throw new Error(errorData.message || t('forgotPassword.requestFailed'))
      }

      setSuccess(true)
    } catch (err) {
      const message = err instanceof Error ? err.message : t('forgotPassword.requestError')
      setError(message)
    } finally {
      setLoading(false)
    }
  }

  return (
    <div className="flex min-h-screen items-center justify-center bg-background p-4">
      <Card className="w-full max-w-md">
        <CardHeader>
          <CardTitle className="text-2xl text-foreground">{`${APP_NAME} - ${t('forgotPassword.title')}`}</CardTitle>
          <CardDescription>{t('forgotPassword.description')}</CardDescription>
        </CardHeader>
        <CardContent>
          {success ? (
            <div className="space-y-4">
              <div className="rounded-md bg-green-500/10 p-3 text-green-600 dark:text-green-400">
                {t('forgotPassword.successMessage')}
              </div>
              <Link to="/login">
                <Button className="w-full">{t('forgotPassword.backToLogin')}</Button>
              </Link>
            </div>
          ) : (
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
                      <FormLabel>{t('forgotPassword.email')}</FormLabel>
                      <FormControl>
                        <Input placeholder={t('forgotPassword.emailPlaceholder')} {...field} />
                      </FormControl>
                      <FormMessage />
                    </FormItem>
                  )}
                />
                <Button type="submit" className="w-full" disabled={loading}>
                  {loading ? t('forgotPassword.submitting') : t('forgotPassword.submit')}
                </Button>
              </form>
            </Form>
          )}
        </CardContent>
        <CardFooter className="flex justify-center">
          <Link to="/login" className="text-sm text-primary hover:underline">
            {t('forgotPassword.backToLogin')}
          </Link>
        </CardFooter>
      </Card>
    </div>
  )
}
