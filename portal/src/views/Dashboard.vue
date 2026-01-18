<script setup>
import { LxTile } from '@wntr/lx-ui';
import { useI18n } from 'vue-i18n';
import useAuthStore from '@/stores/useAuthStore';

const authStore = useAuthStore();

const t = useI18n();
</script>

<template>
  <div class="lx-dashboard">
    <LxTile
      v-if="authStore.isAuthorized"
      icon="registry"
      :label="t.t('pages.credentials.title')"
      :to="{ name: 'credentials' }"
    />
    <LxTile
      icon="inspection"
      :label="t.t('pages.verification.title')"
      :to="{ name: 'verificationFull' }"
    />
  </div>
</template>
