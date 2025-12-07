<script setup>
import { onMounted } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import { LxLoader } from '@wntr/lx-ui';
import useAuthStore from '@/stores/useAuthStore';
import useNotifyStore from '@/stores/useNotifyStore';
import { useI18n } from 'vue-i18n';

import useAppStore from '@/stores/useAppStore';
import useErrors from '@/hooks/useErrors';
import { HttpStatusCode } from 'axios';

const route = useRoute();
const router = useRouter();
const appStore = useAppStore();
const authStore = useAuthStore();
const notify = useNotifyStore();
const translate = useI18n();
const errors = useErrors();
/**
 * @type {Record<string, () => {title: string, description: string}>}
 * @description possible error codes returned by the server: [invalid_callback, invalid_token, server_error, invalid_client, no_session, invalid_rights, blocked, invalid_request
 */
const errCodeMessage = {
  unknown: () => ({
    title: translate.t('pages.auth.errors.unknownTitle'),
    description: translate.t('pages.auth.errors.unknownDescription'),
  }),
  invalid_rights: () => ({
    title: translate.t('pages.auth.errors.invalidRightsTitle'),
    description: translate.t('pages.auth.errors.invalidRightsDescription'),
  }),
  blocked: () => ({
    title: translate.t('pages.auth.errors.blockedTitle'),
    description: translate.t('pages.auth.errors.blockedDescription'),
  }),
  invalid_request: () => ({
    title: translate.t('pages.auth.errors.invalidRequestTitle'),
    description: translate.t('pages.auth.errors.invalidRequestDescription'),
  }),
};

function handleUnknownError() {
  const err = errCodeMessage.unknown();
  appStore.setError(err.title);
}

function handleError(errorCode) {
  const err = (errCodeMessage[errorCode] || errCodeMessage.unknown)();
  appStore.setError(err.title);
}

function handleAuthError(err) {
  const error = errors.get(err);
  if (error.status === HttpStatusCode.Unauthorized) {
    authStore.logout();
  } else if (error.data) {
    notify.pushError(error.data);
  }
  router.replace({ name: 'error' });
}

async function handleSuccessfulAuth() {
  if (authStore.session.st && authStore.isAuthorized) {
    const returnPath = await authStore.getReturnPath();
    if (returnPath) {
      router.replace(returnPath);
      authStore.clearReturnPath();
      return;
    }
    router.replace({ name: 'dashboard' });
  } else if (authStore.session.st) {
    notify.pushError(translate.t('shell.notifications.unknownState'));
  }
}

async function handleAuthCode(code) {
  await authStore.setSessionKey(code);
  // ToDo: delete or use code below to get token from one time token
  // const token = await startSession(code);
  // await authStore.setSessionKey(token.data.access_token);
  try {
    await authStore.fetchSession();
    handleSuccessfulAuth();
  } catch (err) {
    handleAuthError(err);
  }
}

onMounted(async () => {
  if (route.query?.error) {
    handleError(route.query.error.toString());
    return;
  }

  if (route.query?.code) {
    await handleAuthCode(route.query.code.toString());
    return;
  }

  handleUnknownError();
});
</script>
<template>
  <div class="lx-plate">
    <LxLoader loading size="l" />
    <p class="lx-description">{{ $t('pages.auth.description') }}</p>
  </div>
</template>
