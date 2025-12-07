<script setup lang="ts">
import router from '@/router';
import useAuthStore from '@/stores/useAuthStore';
import useNotifyStore from '@/stores/useNotifyStore';
import useViewStore from '@/stores/useViewStore';
import { LxForm, LxRow, LxTextInput, LxList } from '@wntr/lx-ui';
import { onMounted, ref } from 'vue';
import { useI18n } from 'vue-i18n';

const authStore = useAuthStore();

const notification = useNotifyStore();
const t = useI18n();
const viewStore = useViewStore();

const logout = async () => {
  try {
    const resp = await authStore.logout();
    if (resp.status === 200 && resp.data !== '') {
      window.location.href = resp.data;
    } else {
      await router.push({ name: 'home' });
      notification.pushSuccess(t.t('shell.notifications.logOut'));
    }
  } catch (err) {
    notification.pushError(`Neparedzēta kļūda: ${err.message}`);
  }
};
function buttonClicked(actionName) {
  if (actionName === 'logout') {
    logout();
  }
}

const formButtons = ref([
  {
    id: 'logout',
    name: t.t('shell.logOut'),
    icon: 'logout',
    kind: 'tertiary',
    destructive: true,
  },
]);

const scopeItems = authStore.session.scope.map((item, index) => ({
  id: index,
  name: item,
}));

onMounted(() => {
  viewStore.goBack = true;
});
</script>
<template>
  <LxForm
    :action-definitions="formButtons"
    :show-header="true"
    :showPreHeaderInfo="true"
    :column-count="2"
    @buttonClick="buttonClicked"
  >
    <template #pre-header>Sesijas darbības laiks: {{ authStore.session.secondsToLive }}</template>
    <template #header>{{ authStore.fullName }}</template>

    <LxRow :label="$t('pages.userProfile.personalIdLabel')">
      <LxTextInput v-model="authStore.session.code" :read-only="true" />
    </LxRow>
    <LxRow :label="$t('pages.userProfile.roleLabel')">
      <LxTextInput v-model="authStore.session.role" :read-only="true" />
    </LxRow>
    <LxRow :label="$t('pages.userProfile.organizationLabel')">
      <LxTextInput v-model="authStore.session.org_name" :read-only="true" />
    </LxRow>
    <LxRow :label="$t('pages.userProfile.scope')">
      <LxList id="id" :items="scopeItems"> </LxList>
    </LxRow>
  </LxForm>
</template>
