import { defineStore } from 'pinia';
import { LxAuthStore, LxAuthService } from '@wntr/lx-ui';
import { AUTH_KEY_TOKEN_SESSION, AUTH_SCOPE, APP_CONFIG } from '@/constants';
import { keepAlive, logout, session } from '@/services/authService';

export default defineStore(
  'authStore',
  LxAuthStore(
    // @ts-ignore
    (authUrl, publicUrl, clientId, scope, authSessionKey) => ({
      ...LxAuthService(authUrl, publicUrl, clientId, scope, authSessionKey),
      // ToDo: configure custom auth endpoints
      session,
      keepAlive,
      logout,
    }),
    APP_CONFIG.authUrl,
    APP_CONFIG.publicUrl,
    APP_CONFIG.clientId,
    AUTH_SCOPE,
    AUTH_KEY_TOKEN_SESSION,
    // Custom state extension
  ),
);
