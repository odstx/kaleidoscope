import { useEffect, useState, useRef } from "react"
import { useNavigate, useParams } from "react-router-dom"
import { useAuth } from "@/contexts/AuthContext"

export default function MicroAppPage() {
  const { appname } = useParams<{ appname: string }>()
  const navigate = useNavigate()
  const { isAuthenticated } = useAuth()
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)
  const containerRef = useRef<HTMLDivElement>(null)
  const loadedRef = useRef(false)

  const scriptUrl = appname ? `/app/${appname}.js` : null
  const tagName = appname ? `${appname}-app` : null

  useEffect(() => {
    if (!isAuthenticated) {
      navigate("/login")
      return
    }

    if (!appname || !scriptUrl || !tagName) {
      return
    }

    if (loadedRef.current) return

    const loadMicroApp = async () => {
      try {
        setLoading(true)

        const existingScript = document.querySelector(`script[src="${scriptUrl}"]`)
        if (!existingScript) {
          await new Promise<void>((resolve, reject) => {
            const script = document.createElement("script")
            script.src = scriptUrl
            script.type = "module"
            script.onload = () => resolve()
            script.onerror = () => reject(new Error(`Failed to load ${appname}`))
            document.head.appendChild(script)
          })
        }

        await customElements.whenDefined(tagName)

        setLoading(false)
        loadedRef.current = true
      } catch (err) {
        setError(err instanceof Error ? err.message : "Failed to load micro app")
        setLoading(false)
      }
    }

    loadMicroApp()
  }, [isAuthenticated, navigate, appname, scriptUrl, tagName])

  if (!isAuthenticated) {
    return null
  }

  if (!appname || !scriptUrl || !tagName) {
    return (
      <div className="min-h-screen bg-background p-8 flex items-center justify-center">
        <div className="text-center">
          <p className="text-destructive mb-4">App name is required</p>
        </div>
      </div>
    )
  }

  if (loading) {
    return (
      <div className="min-h-screen bg-background p-8 flex items-center justify-center">
        <div className="text-center">
          <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-primary mx-auto mb-4"></div>
          <p className="text-muted-foreground">Loading {appname}... </p>
        </div>
      </div>
    )
  }

  if (error) {
    return (
      <div className="min-h-screen bg-background p-8 flex items-center justify-center">
        <div className="text-center">
          <p className="text-destructive mb-4">{error}</p>
          <button
            onClick={() => window.location.reload()}
            className="px-4 py-2 bg-primary text-primary-foreground rounded-md hover:bg-primary/90"
          >
            Retry
          </button>
        </div>
      </div>
    )
  }

  return (
    <div className="min-h-screen bg-background" ref={containerRef}>
      <div ref={(el) => {
        if (el && !el.firstChild) {
          const customElement = document.createElement(tagName)
          el.appendChild(customElement)
        }
      }} />
    </div>
  )
}
