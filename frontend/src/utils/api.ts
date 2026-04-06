export interface SystemInfo {
  version: string;
  build_id: string;
  build_time: string;
  git_commit: string;
  openapi_path: string;
}

export async function getSystemInfo(): Promise<SystemInfo> {
  const baseUrl = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8000';
  const response = await fetch(`${baseUrl}/api/v1/system/info`);
  
  if (!response.ok) {
    throw new Error(`Failed to fetch system info: ${response.status}`);
  }
  
  return response.json();
}

export interface TOTPSetupResponse {
  secret: string;
  url: string;
}

export interface LoginResponse {
  message: string;
  user: {
    id: number;
    email: string;
    totp_enabled: boolean;
  };
  token: string;
}

export interface LoginError {
  error: string;
  totp_required?: boolean;
}

const getBaseUrl = (): string => {
  return import.meta.env.VITE_API_BASE_URL || 'http://localhost:8000';
};

const getAuthHeaders = (): HeadersInit => {
  const token = localStorage.getItem('token');
  return {
    'Content-Type': 'application/json',
    ...(token ? { 'Authorization': `Bearer ${token}` } : {}),
  };
};

export async function setupTOTP(): Promise<TOTPSetupResponse> {
  const response = await fetch(`${getBaseUrl()}/api/v1/users/totp/setup`, {
    method: 'GET',
    headers: getAuthHeaders(),
  });

  if (!response.ok) {
    const error = await response.json();
    throw new Error(error.error || 'Failed to setup TOTP');
  }

  return response.json();
}

export async function verifyTOTP(code: string): Promise<void> {
  const response = await fetch(`${getBaseUrl()}/api/v1/users/totp/verify`, {
    method: 'POST',
    headers: getAuthHeaders(),
    body: JSON.stringify({ code }),
  });

  if (!response.ok) {
    const error = await response.json();
    throw new Error(error.error || 'Failed to verify TOTP');
  }
}

export async function enableTOTP(): Promise<void> {
  const response = await fetch(`${getBaseUrl()}/api/v1/users/totp/enable`, {
    method: 'POST',
    headers: getAuthHeaders(),
  });

  if (!response.ok) {
    const error = await response.json();
    throw new Error(error.error || 'Failed to enable TOTP');
  }
}

export async function disableTOTP(): Promise<void> {
  const response = await fetch(`${getBaseUrl()}/api/v1/users/totp/disable`, {
    method: 'POST',
    headers: getAuthHeaders(),
  });

  if (!response.ok) {
    const error = await response.json();
    throw new Error(error.error || 'Failed to disable TOTP');
  }
}
