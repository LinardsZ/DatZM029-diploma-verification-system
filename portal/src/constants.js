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
  supportedLocales: ['lv', 'en'],
  fallbackLocale: 'en',
};

export const AUTH_KEY_TOKEN_SESSION = 'blockchain-sessionkey';
export const SYSTEM_NAME = 'blockchain';
export const AUTH_SCOPE = 'vpm';
export const AUTH_TYPE = 'VPM';
