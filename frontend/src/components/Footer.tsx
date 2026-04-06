import { useEffect, useState } from 'react';
import { useTranslation } from 'react-i18next';
import { getSystemInfo } from '@/utils/api';
import type { SystemInfo } from '@/utils/api';

export function Footer() {
  const { t } = useTranslation();
  const [systemInfo, setSystemInfo] = useState<SystemInfo | null>(null);
  const [error, setError] = useState<boolean>(false);

  const isDebugMode = window.location.hostname.includes('local');

  useEffect(() => {
    if (!isDebugMode) return;
    getSystemInfo()
      .then(setSystemInfo)
      .catch(() => setError(true));
  }, [isDebugMode]);

  if (!isDebugMode) return null;

  return (
    <footer className="fixed bottom-0 left-0 right-0 bg-muted border-t border-border py-2 px-4">
      <div className="max-w-7xl mx-auto flex items-center justify-center text-sm text-muted-foreground gap-4">
        {error ? (
          <span>{t("footer.error")}</span>
        ) : systemInfo ? (
          <>
            <span>{t("footer.version")} <strong>{systemInfo.version}</strong></span>
            <span>{t("footer.buildId")} <strong>{systemInfo.build_id}</strong></span>
              <a
                href={`${import.meta.env.VITE_API_BASE_URL || 'http://localhost:8000'}${systemInfo.openapi_path}`}
                target="_blank"
                rel="noopener noreferrer"
                className="text-primary hover:underline"
              >
              OpenAPI
            </a>
          </>
        ) : (
          <span>{t("footer.loading")}</span>
        )}
      </div>
    </footer>
  );
}
