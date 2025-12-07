export const APP_CONFIG = {
  // @ts-ignore
  apiUrl: window.config.serviceUrl,
  // @ts-ignore
  authUrl: window.config.authUrl,
  // @ts-ignore
  publicUrl: window.config.publicUrl,
  // @ts-ignore
  clientId: window.config.clientId,
  // @ts-ignore
  environment: window.config.environment,
  defaultLocale: 'lv',
  supportedLocales: ['lv'],
  fallbackLocale: 'lv',
};

export const AUTH_KEY_TOKEN_SESSION = '$REPO_NAME_LOWER-sessionkey';
export const SYSTEM_NAME = '$REPO_NAME_LOWER';

export const AUTH_SCOPE = 'vpm';
export const AUTH_TYPE = 'VPM';
