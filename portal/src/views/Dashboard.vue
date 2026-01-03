<script setup>
import { LxTile } from '@wntr/lx-ui';
import { useI18n } from 'vue-i18n';

const t = useI18n();
</script>

<template>
  <div class="lx-dashboard">
    <LxTile
      icon="add"
      :label="t.t('pages.newCredential.title')"
      :to="{ name: 'newCredential' }"
    />
    <LxTile
      icon="search"
      :label="t.t('pages.verification.title')"
      :to="{ name: 'verification' }"
    />
    <LxTile
      icon="inspection"
      :label="t.t('pages.verificationFull.title')"
      :to="{ name: 'verificationFull' }"
    />
  </div>
</template>
