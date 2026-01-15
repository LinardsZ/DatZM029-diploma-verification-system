<script setup>
import { ref, onMounted, computed } from 'vue';
import {
  LxList,
  LxStack,
  LxStateDisplay,
  LxIcon,
  lxDateUtils,
} from '@wntr/lx-ui';
import {
  getCredentialList,
  revokeCredential,
} from '@/services/credentialService';
import { useI18n } from 'vue-i18n';
import useConfirmStore from '@/stores/useConfirmStore';
import useNotifyStore from '@/stores/useNotifyStore';
import { useRouter } from 'vue-router';
import useAuthStore from '@/stores/useAuthStore';
import { useWindowSize } from '@vueuse/core';

const t = useI18n();
const confirm = useConfirmStore();
const notify = useNotifyStore();
const router = useRouter();
const authStore = useAuthStore();

const { width } = useWindowSize();

const listData = ref([]);
const loading = ref(false);

async function loadList() {
  loading.value = true;
  try {
    const res = await getCredentialList(authStore.session.institution.id);
    listData.value = res.data.credentials;
  } catch (e) {
    console.error(e);
  } finally {
    loading.value = false;
  }
}

const statusDict = [
  {
    value: 'Valid',
    displayName: t.t('pages.newCredential.form.valid'),
    displayType: 'new',
  },
  {
    value: 'Revoked',
    displayName: t.t('pages.newCredential.form.invalid'),
    displayType: 'inactive',
  },
];

const listDisplay = computed(() =>
  listData.value.map((item) => ({
    ...item,
    isActive: item.status == 'Valid',
  })),
);

async function revoke(itemId) {
  try {
    const credentialId = itemId.split('_')[1];
    confirm.$state.confirmDialogState.primaryBusy = true;
    const res = await revokeCredential(credentialId);
    confirm.$state.confirmDialogState.primaryBusy = false;
    if (res.status !== 204) {
      return;
    }
    notify.pushSuccess(t.t('pages.credentials.revokeSuccess'));

    loadList();
  } catch (e) {
    console.error(e);
    notify.pushError(t.t('pages.credentials.revokeError'));
  }
}

function revokeItem(_, itemId) {
  confirm.pushSimple(
    t.t('pages.credentials.revokeTitle'),
    t.t('pages.credentials.revokeDesc'),
    () => revoke(itemId),
  );
}

onMounted(() => {
  loadList();
});
</script>
<template>
  <div class="cred-list-wrapper">
    <LxList
      :items="listDisplay"
      :loading="loading"
      listType="1"
      :actionDefinitions="[
        {
          id: 'revoke',
          label: t.t('pages.credentials.revoke'),
          icon: 'block',
          visibleByAttribute: 'isActive',
        },
      ]"
      :toolbarActionDefinitions="[
        {
          id: 'add',
          name: t.t('pages.credentials.addCredential'),
          icon: 'add',
          kind: 'primary',
        },
      ]"
      @actionClick="revokeItem"
      @toolbar-action-click="router.push({ name: 'newCredential' })"
    >
      <template #customItem="item">
        <LxStack
          kind="compact"
          :orientation="width > 1000 ? 'horizontal' : 'vertical'"
          verticalAlignment="center"
          mode="grid"
          :horizontalConfig="['*', 'auto']"
          :verticalConfig="width > 1000 ? ['*'] : ['*', 'auto']"
        >
          <div style="display: flex; align-items: center">
            <LxIcon
              :value="
                item.credentialType === 'Diploma' ? 'diploma' : 'contract'
              "
            />

            <div style="padding-left: 0.75rem">
              <p class="lx-primary">{{ item.diplomaMetadata.degreeName }}</p>
            </div>
          </div>
          <div
            style="
              display: grid;
              grid-template-columns: 1fr 1fr;
              align-items: center;
              gap: 1rem;
              width: 15rem;
            "
          >
            <LxStack
              orientation="horizontal"
              verticalAlignment="center"
              kind="compact"
              class="date-wrapper"
            >
              <LxIcon value="calendar" />
              <p class="lx-secondary">
                {{ lxDateUtils.formatDate(item.diplomaMetadata.issueDate) }}
              </p>
            </LxStack>

            <LxStateDisplay :value="item?.status" :dictionary="statusDict" />
          </div>
        </LxStack>
      </template>
    </LxList>
  </div>
</template>
