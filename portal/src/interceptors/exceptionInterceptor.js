import useNotifyStore from '@/stores/useNotifyStore';
import { HttpStatusCode } from 'axios';
import i18nInstance from '@/i18n';

const dbErrorRegex = /\[err:([^[\]]+)\]\s(.+)/;
const errorRegex = /err:(.*)/;

function getDatabaseErrors(error) {
  const hasErrors = Array.isArray(error.response?.data?.errors);
  const dbErrors = hasErrors && error.response?.data?.errors.filter((x) => x.type === 'ErrExec');
  if (dbErrors && dbErrors.length > 0) {
    return dbErrors.map((fullError) => {
      const match = fullError?.message?.match(dbErrorRegex);
      if (match && match.length > 2) {
        const code = match[1];
        const message = match[2];
        return { code, message };
      }
      return { code: 'unknown', message: fullError };
    });
  }
  return undefined;
}

/**
 * @param {import('axios').AxiosError} error
 * @param {import('vue-i18n').VueI18nTranslation} $t
 * @returns {string | undefined}
 */
function getDatabaseErrorMessage(error, $t) {
  const dbErrors = getDatabaseErrors(error);
  if (dbErrors && dbErrors.length > 0) {
    return dbErrors
      .map((x) => {
        // if message is in translation file in format `errors.{code}` (without `err:` prefix) then use it
        const translation = $t(`errors.${x.code}`);
        if (translation !== `errors.${x.code}`) {
          return translation;
        }
        return x.message;
      })
      .join('\n');
  }
  return undefined;
}

/**
 * @param {import('axios').AxiosError} error
 * @param {import('vue-i18n').VueI18nTranslation} $t
 * @returns {string | undefined}
 */
function getErrorMessageFromText(error, $t) {
  const match = error.response?.data?.match(errorRegex);
  if (match && match.length === 2) {
    // if message is in translation file in format `errors.{code}` (without `err:` prefix) then use it
    const translation = $t(`errors.${match[1]}`);
    if (translation !== `errors.${match[1]}`) {
      return translation;
    }
  }
  return undefined;
}

/**
 * @param {import('axios').AxiosError} error
 */
export function exceptionInterceptor(error) {
  if (!error.response && error?.code === 'ERR_NETWORK') {
    const notify = useNotifyStore();
    const $t = i18nInstance.global.t;
    notify.pushError(getErrorMessageFromText(error, $t) || $t('errors.networkError'));
    return Promise.reject(error);
  }
  const { status } = error.response;

  if (status === HttpStatusCode.BadRequest) {
    const notify = useNotifyStore();
    const $t = i18nInstance.global.t;
    notify.pushError(getErrorMessageFromText(error, $t) || $t('errors.http400'));
  } else if (status === HttpStatusCode.UnprocessableEntity || status === HttpStatusCode.NotFound) {
    const notify = useNotifyStore();
    const $t = i18nInstance.global.t;
    const dbError = getDatabaseErrorMessage(error, $t);
    notify.pushError(dbError || $t('errors.http422'));
  } else if (status === HttpStatusCode.InternalServerError) {
    const notify = useNotifyStore();
    const $t = i18nInstance.global.t;
    notify.pushError(error.response?.data || $t('errors.http500'));
  }

  return Promise.reject(error);
}
