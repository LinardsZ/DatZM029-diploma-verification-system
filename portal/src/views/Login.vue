<script setup>
import { useI18n } from 'vue-i18n';
import { LxButton } from '@wntr/lx-ui';
import useAuthStore from '@/stores/useAuthStore';
import router from '@/router';

const authStore = useAuthStore();
const i18n = useI18n();

function login() {
  router.push({ name: 'auth' });
}
</script>

<template>
  <LxButton
    :href="{ name: 'dashboard' }"
    :label="i18n.t('pages.home.buttons.loginAuthorized')"
    icon="dashboard"
  />
  <LxButton
    v-if="!authStore.isAuthorized"
    @click="login"
    kind="secondary"
    :label="i18n.t('pages.home.buttons.loginUnauthorized')"
    icon="next"
  />
</template>
