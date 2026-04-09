export interface OidcConfig {
  enabled: boolean
  issuerUrl: string
  clientId: string
  redirectUri: string
}

export interface FrontendConfig {
  oidcClientId: string
  oidcEnabled: boolean
  oidcIssuerUrl: string
  oidcRedirectUri: string
}

let cachedConfig: OidcConfig | null = null

export async function fetchFrontendConfig(): Promise<OidcConfig> {
  if (cachedConfig) {
    return cachedConfig
  }

  try {
    const response = await fetch(`${import.meta.env.VITE_API_BASE_URL}/api/v1/system/config`)
    if (!response.ok) {
      throw new Error('Failed to fetch config')
    }

    const data: FrontendConfig = await response.json()
    cachedConfig = {
      enabled: data.oidcEnabled,
      issuerUrl: data.oidcIssuerUrl || '',
      clientId: data.oidcClientId || '',
      redirectUri: data.oidcRedirectUri || '',
    }
    return cachedConfig
  } catch {
    cachedConfig = {
      enabled: false,
      issuerUrl: '',
      clientId: '',
      redirectUri: '',
    }
    return cachedConfig
  }
}



export async function getOidcAuthUrl(): Promise<string> {
  const response = await fetch(`${import.meta.env.VITE_API_BASE_URL}/api/v1/users/oidc/login`, {
    method: 'GET',
    headers: {
      'Content-Type': 'application/json',
    },
  })

  if (!response.ok) {
    const error = await response.json().catch(() => ({ error: 'Failed to get OIDC auth URL' }))
    throw new Error(error.error || 'Failed to get OIDC auth URL')
  }

  const data = await response.json()
  return data.url
}

export async function handleOidcCallback(code: string): Promise<{ token: string }> {
  const response = await fetch(`${import.meta.env.VITE_API_BASE_URL}/api/v1/users/oidc/callback`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify({ code }),
  })

  if (!response.ok) {
    const error = await response.json().catch(() => ({ error: 'OIDC callback failed' }))
    throw new Error(error.error || 'OIDC callback failed')
  }

  const data = await response.json()
  return { token: data.token }
}