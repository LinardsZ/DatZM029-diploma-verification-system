import { apiAuth } from '@/apiAuth';
import { APP_CONFIG } from '@/constants';

export function session() {
  return apiAuth().get('/api/1.0/session');
}

export function keepAlive() {
  return apiAuth().get('/api/1.0/session/keep-alive');
}

export async function logout() {
  await apiAuth().delete('/api/1.0/session');
}

/**
 * @typedef {Object} TokenIssueResponse
 * @property {string} token_type - The type of token (e.g., access_token, refresh_token, id_token).
 * @property {number} expires_in - The number of seconds until the token expires (optional for some token types).
 * @property {string} access_token - The access token (if applicable).
 */

/**
 * @param {string} ott - The one-time token (OTT) received from the authorization server.
 * @returns { Promise<import('axios').AxiosResponse<TokenIssueResponse>> } A promise that resolves to the token issue response.
 * @throws {Error} If there is an error during the token issue process.
 */
export async function startSession(ott) {
  return apiAuth({ allowAnonymous: true }).post('/1.0/token', {
    code: ott,
    client_id: APP_CONFIG.clientId,
    grant_type: 'authorization_code',
  });
}
