<script setup>
import { LxErrorPage } from '@wntr/lx-ui';
import { useI18n } from 'vue-i18n';
import { useRouter } from 'vue-router';
import { onMounted, onUnmounted, computed } from 'vue';
import useViewStore from '@/stores/useViewStore';

const viewStore = useViewStore();
const router = useRouter();
const $t = useI18n().t;

function action(actionName) {
  if (actionName === 'authorize') {
    router.push({
      name: 'home',
      query: { returnPath: router.currentRoute.value?.query?.returnPath },
    });
    return;
  }
  router.go(-1);
}
onMounted(() => {
  viewStore.hideHeader();
});
onUnmounted(() => {
  viewStore.showHeader();
});
const actionDefinitions = computed(() => [
  { id: 'back', name: $t('pages.error.goBack'), icon: 'undo' },
  { id: 'authorize', name: $t('pages.error.authorizeButton'), icon: 'next' },
]);
</script>
<template>
  <LxErrorPage
    kind="401"
    @actionClick="action"
    :action-definitions="actionDefinitions"
  ></LxErrorPage>
</template>
