<script setup lang="ts">
import { onMounted } from 'vue';
import { useRouter } from 'vue-router';
import { useI18n } from 'vue-i18n';
import { LxSiteMap } from '@wntr/lx-ui';

import useViewStore from '@/stores/useViewStore';

const viewStore = useViewStore();
const i18n = useI18n();
const { t } = useI18n();
const router = useRouter();
const routes = router.getRoutes();

onMounted(() => {
  viewStore.title = i18n.t('pages.siteMap.title');
  viewStore.description = i18n.t('pages.siteMap.description');
  viewStore.goBack = true;
  viewStore.backName = i18n.t('pages.dashboard.title');
});
</script>
<template>
  <div class="top-routes">
    <LxSiteMap
      :routes="routes"
      :label="i18n.t('pages.siteMap.top3Routes')"
      category="useful"
      :translator="t"
    ></LxSiteMap>
  </div>
</template>
