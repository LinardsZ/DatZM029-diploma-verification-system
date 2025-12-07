import { AUTH_KEY_TOKEN_SESSION } from '@/constants';

const storage = () => sessionStorage;

export const setSessionKey = (key) => {
  storage().setItem(AUTH_KEY_TOKEN_SESSION, key);
};

export const getSessionKey = () => storage().getItem(AUTH_KEY_TOKEN_SESSION);

export const removeSessionKey = () => storage().removeItem(AUTH_KEY_TOKEN_SESSION);
