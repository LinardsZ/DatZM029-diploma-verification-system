import { removeSessionKey } from '@/utils/sessionUtils';
import { HttpStatusCode } from 'axios';

/**
 * @param {import('axios').AxiosError} error
 */
export async function unauthorizedInterceptor(error) {
  if (!error.response && error?.code === 'ERR_NETWORK') {
    return Promise.reject(error);
  }
  const { status } = error.response;

  if (status === HttpStatusCode.Unauthorized) {
    [() => sessionStorage.clear(), () => removeSessionKey()].forEach((f) => {
      try {
        f();
      } catch (e) {
        // try next function even if one fails
        console.error('Error while resetting session storage', e);
      }
    });
    window.location.reload();
  }

  return Promise.reject(error);
}
