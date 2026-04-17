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
  agentEnabled: boolean
  enableRegistration: boolean
}

let cachedConfig: OidcConfig | null = null
let cachedAgentEnabled: boolean | null = null
let cachedEnableRegistration: boolean | null = null

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
    cachedAgentEnabled = data.agentEnabled ?? true
    cachedEnableRegistration = data.enableRegistration ?? true
    return cachedConfig
  } catch {
    cachedConfig = {
      enabled: false,
      issuerUrl: '',
      clientId: '',
      redirectUri: '',
    }
    cachedAgentEnabled = true
    cachedEnableRegistration = true
    return cachedConfig
  }
}

export async function fetchAgentEnabled(): Promise<boolean> {
  if (cachedAgentEnabled !== null) {
    return cachedAgentEnabled
  }
  await fetchFrontendConfig()
  return cachedAgentEnabled ?? true
}

export async function fetchEnableRegistration(): Promise<boolean> {
  if (cachedEnableRegistration !== null) {
    return cachedEnableRegistration
  }
  await fetchFrontendConfig()
  return cachedEnableRegistration ?? true
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