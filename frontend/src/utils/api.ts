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
