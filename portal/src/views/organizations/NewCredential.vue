<script setup>
import { ref, computed, watch } from 'vue';
import {
  LxForm,
  LxRow,
  LxTextInput,
  LxDateTimePicker,
  LxFileUploader,
  LxValuePicker,
  LxStateDisplay,
} from '@wntr/lx-ui';
import { useI18n } from 'vue-i18n';
import { useRouter } from 'vue-router';
import useNotifyStore from '@/stores/useNotifyStore';
import { getFileHash } from '@/utils/generalUtils';

const t = useI18n();
const router = useRouter();
const notify = useNotifyStore();

const diplomaFile = ref();
const loading = ref(false);

const inputData = ref({
  id: null, // TODO: remove
  graduatePublicKey: null, // TODO: remove
  issuerSignature: null, // TODO: remove
  issuerPublicKey: null, // TODO: remove

  diplomaHash: null,
  diplomaMetadata: {
    universityName: null,
    degreeName: null,
    issueDate: null,
    expiryDate: null,
  },
  status: null,
  credentialType: null,
});

const statusItmes = [
  { id: 'Valid', name: 'Valid' },
  { id: 'Invalid', name: 'Invalid' },
];

const statusDict = [
  {
    value: 'Valid',
    displayName: t.t('pages.newCredential.form.valid'),
    displayType: 'new',
  },
  {
    value: 'Invalid',
    displayName: t.t('pages.newCredential.form.invalid'),
    displayType: 'inactive',
  },
];

const credentialTypeItems = [
  { id: 'Diploma', name: t.t('pages.newCredential.form.diploma') },
  { id: 'Certificate', name: t.t('pages.newCredential.form.certificate') },
];

watch(diplomaFile, async (newValue) => {
  if (newValue && newValue.length > 0) {
    const res = await getFileHash(newValue);

    inputData.value.diplomaHash = res;
  } else {
    inputData.value.diplomaHash = null;
  }
});

function saveDocument() {
  return new Promise((resolve) => {
    setTimeout(() => {
      const mockResponse = {
        status: 201,
        message: "Document saved successfully!",
      };
      resolve(mockResponse);
    }, 2000);
  });
}

async function formActionClick(id) {
  if (id === 'save') {
    // TODO: add validation
    loading.value = true;
    const res = await saveDocument();
    if (res.status === 201) {
      notify.pushSuccess(t.t('pages.newCredential.form.documentCreated'));
      router.push({ name: 'dashboard' });
    } else {
      notify.pushError(t.t('pages.newCredential.form.creationFailed'));
    }
    loading.value = false;
  } else {
    // TODO: can add route guard and confirmStore
    router.push({ name: 'dashboard' });
  }
}
</script>
<template>
  <div>
    <!-- TODO: add componente translations -->
    <LxForm
      :showHeader="false"
      :columnCount="2"
      requiredMode="required-asterisk"
      :actionDefinitions="[
        {
          id: 'save',
          name: t.t('pages.newCredential.form.save'),
          icon: 'save',
          kind: 'primary',
          busy: loading,
        },
        {
          id: 'cancel',
          name: t.t('pages.newCredential.form.cancel'),
          icon: 'cancel',
          kind: 'secondary',
          disabled: loading,
        },
      ]"
      @button-click="formActionClick"
    >
      <LxRow :label="t.t('pages.newCredential.form.id')" :required="true">
        <template #info>TODO: get from session</template>
        <LxTextInput v-model="inputData.id" :disabled="loading" />
      </LxRow>
      <LxRow
        :label="t.t('pages.newCredential.form.graduatePublicKey')"
        :required="true"
      >
        <template #info>TODO: get from session</template>
        <LxTextInput
          v-model="inputData.graduatePublicKey"
          :disabled="loading"
        />
      </LxRow>
      <LxRow
        :label="t.t('pages.newCredential.form.issuerSignature')"
        :required="true"
      >
        <template #info>TODO: get from session</template>
        <LxTextInput v-model="inputData.issuerSignature" :disabled="loading" />
      </LxRow>
      <LxRow
        :label="t.t('pages.newCredential.form.issuerPublicKey')"
        :required="true"
      >
        <template #info>TODO: get from session</template>
        <LxTextInput v-model="inputData.issuerPublicKey" :disabled="loading" />
      </LxRow>

      <LxRow
        :label="t.t('pages.newCredential.form.universityName')"
        :required="true"
      >
        <template #info>TODO: get from session</template>
        <LxTextInput
          v-model="inputData.diplomaMetadata.universityName"
          :disabled="loading"
        />
      </LxRow>
      <LxRow
        :label="t.t('pages.newCredential.form.degreeName')"
        :required="true"
      >
        <LxTextInput
          v-model="inputData.diplomaMetadata.degreeName"
          :disabled="loading"
        />
      </LxRow>

      <LxRow :label="t.t('pages.newCredential.form.status')" :required="true">
        <LxValuePicker
          v-model="inputData.status"
          :items="statusItmes"
          variant="tags-custom"
          :disabled="loading"
        >
          <template #customItem="item">
            <LxStateDisplay :value="item.name" :dictionary="statusDict" />
          </template>
        </LxValuePicker>
      </LxRow>
      <LxRow
        :label="t.t('pages.newCredential.form.credentialType')"
        :required="true"
      >
        <LxValuePicker
          v-model="inputData.credentialType"
          :items="credentialTypeItems"
          :disabled="loading"
        />
      </LxRow>

      <LxRow
        :label="t.t('pages.newCredential.form.issueDate')"
        :required="true"
      >
        {{ t.locale }}
        <LxDateTimePicker
          v-model="inputData.diplomaMetadata.issueDate"
          :disabled="loading"
          :locale="{
            locale: t.locale === 'lv' ? 'lv' : 'en',
          }"
          :texts="t.tm('pages.verification.dateTimePicker')"
        />
      </LxRow>
      <LxRow :label="t.t('pages.newCredential.form.expiryDate')">
        <LxDateTimePicker
          v-model="inputData.diplomaMetadata.expiryDate"
          :disabled="loading"
          :locale="{
            locale: t.locale === 'lv' ? 'lv' : 'en',
          }"
          :texts="t.tm('pages.verification.dateTimePicker')"
        />
      </LxRow>

      <LxRow
        :label="t.t('pages.newCredential.form.degreeFile')"
        :required="true"
      >
        <LxFileUploader
          v-model="diplomaFile"
          :allowedFileExtensions="['.pdf']"
          :disabled="loading"
          :texts="t.tm('pages.verification.fileUploader')"
        />
      </LxRow>
    </LxForm>
  </div>
</template>
