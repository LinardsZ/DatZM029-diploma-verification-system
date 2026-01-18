<script setup>
import { useI18n } from 'vue-i18n';
import { LxButton } from '@wntr/lx-ui';
import useAuthStore from '@/stores/useAuthStore';

const authStore = useAuthStore();
const i18n = useI18n();
</script>

<template>
  <LxButton
    :href="{ name: 'dashboard' }"
    :label="i18n.t('pages.home.buttons.loginAuthorized')"
    icon="dashboard"
  />
  <LxButton
    v-if="!authStore.isAuthorized"
    :href="{ name: 'auth' }"
    kind="secondary"
    :label="i18n.t('pages.home.buttons.loginUnauthorized')"
    icon="next"
  />
</template>
