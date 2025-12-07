import axios, { AxiosError, HttpStatusCode } from 'axios';
import { APP_CONFIG } from '@/constants';
import { exceptionInterceptor } from '@/interceptors/exceptionInterceptor';
import { unauthorizedInterceptor } from '@/interceptors/unauthorizedInterceptor';
import { getSessionKey } from '@/utils/sessionUtils';

function unauthorizedResp() {
  return Promise.reject(
    new AxiosError('no token', null, null, null, {
      status: HttpStatusCode.Unauthorized,
      data: null,
      statusText: null,
      headers: null,
      config: null,
    })
  );
}

/**
 *
 * @param {Object} [options]
 * @param {boolean} [options.useExceptionInterceptor] - Whether to use exception interceptor or not. Default is true.
 * @param {boolean} [options.allowAnonymous] - Whether to allow anonymous requests or not. Default is false.
 */
export const apiAuth = (options = {}) => {
  const { useExceptionInterceptor = true, allowAnonymous = false } = options || {};
  const token = getSessionKey();
  if (!token && !allowAnonymous) {
    return {
      get: unauthorizedResp,
      post: unauthorizedResp,
      put: unauthorizedResp,
      patch: unauthorizedResp,
      delete: unauthorizedResp,
    };
  }
  const http = axios.create({
    baseURL: APP_CONFIG.authUrl,
    withCredentials: false,
    headers: {
      Accept: 'application/json',
      'Content-Type': 'application/json',
      Authorization: `Bearer ${getSessionKey()}`,
    },
  });
  http.interceptors.response.use(null, unauthorizedInterceptor);
  http.interceptors.response.use(null, useExceptionInterceptor && exceptionInterceptor);
  return http;
};
