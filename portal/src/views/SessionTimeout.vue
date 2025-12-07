<script setup>
import useViewStore from '@/stores/useViewStore';
import { onMounted, onUnmounted } from 'vue';
import { LxErrorPage } from '@wntr/lx-ui';

import { useRouter } from 'vue-router';

const router = useRouter();

const viewStore = useViewStore();

function action() {
  router.push({ name: 'home' });
}

onMounted(() => {
  viewStore.hideHeader();
});
onUnmounted(() => {
  viewStore.showHeader();
});
</script>
<template>
  <LxErrorPage
    kind="sessionTimeout"
    @actionClick="action"
    :action-definitions="[{ id: 'back', name: $t('pages.error.goHome'), icon: 'undo' }]"
    :title="$t('pages.sessionTimeout.title')"
    :description="$t('pages.sessionTimeout.description')"
  ></LxErrorPage>
</template>
