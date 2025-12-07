<script setup>
import { watch } from 'vue';
import { useRoute } from 'vue-router';
import { useI18n } from 'vue-i18n';
import useAppStore from '@/stores/useAppStore';
import Error404 from '@/views/Error404.vue';

const i18n = useI18n();
const route = useRoute();
const appStore = useAppStore();

watch(
  route,
  () => {
    let title = i18n.t('title.prefix');
    if (typeof route.meta.title === 'function') {
      title += route.meta.title(i18n);
    } else {
      title = i18n.t(route.meta?.title?.toString() || 'title.default');
    }
    document.title = title;
  },
  { immediate: true },
);
</script>
<template>
  <Error404 v-if="appStore.showError" />
  <router-view v-else />
</template>
