<script setup>
import { useI18n } from 'vue-i18n';
import { LxButton } from '@wntr/lx-ui';
import useAuthStore from '@/stores/useAuthStore';
import { AUTH_TYPE } from '@/constants';

const authStore = useAuthStore();
const i18n = useI18n();

function login() {
  /** @todo: change this for real login from authStore
  /* @example: authStore.login('/my-callback-url');
  */
  const callbackPath = null;
  authStore.login(callbackPath, AUTH_TYPE, 'demo');
}
</script>

<template>
  <LxButton
    :href="{ name: 'dashboard' }"
    :label="i18n.t('pages.home.buttons.loginAuthorized')"
    icon="dashboard"
  />
  <!-- <LxButton
    v-if="!authStore.isAuthorized"
    @click="login"
    kind="secondary"
    :label="i18n.t('pages.home.buttons.loginUnauthorized')"
    icon="next"
  /> -->
</template>
