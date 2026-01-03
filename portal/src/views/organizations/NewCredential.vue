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
import { postCredential } from '@/services/credentialService';

const t = useI18n();
const router = useRouter();
const notify = useNotifyStore();

const diplomaFile = ref();
const loading = ref(false);

const inputData = ref({
  diplomaHash: null,
  graduatePublicKey: 'MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA1234...', // TODO: get from idk where
  issuerId: 'lu', // TODO: get from session
  issuerSignature: '3045022100abcd...', // TODO: get from idk where
  diplomaMetadata: {
    universityName: 'Latvijas Universitāte', // TODO: get from session
    degreeName: null,
    issueDate: null,
    expiryDate: null,
  },
  status: null,
  credentialType: null,
});

const statusItems = [
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

const defaultInvalidFields = {
  diplomaMetadata: {
    universityName: null,
    degreeName: null,
    issueDate: null,
  },
  status: null,
  credentialType: null,
  diplomaFile: null,
};

const invalidFields = ref({
  ...defaultInvalidFields,
});

function validateModel() {
  invalidFields.value = {
    ...defaultInvalidFields,
    diplomaMetadata: {
      ...defaultInvalidFields.diplomaMetadata,
    },
  };

  let hasErrors = false;
  if (!inputData.value.diplomaMetadata.degreeName) {
    invalidFields.value.diplomaMetadata.degreeName = t.t(
      'pages.newCredential.form.invalidMessages.degreeName',
    );
    hasErrors = true;
  }
  if (!inputData.value.diplomaMetadata.issueDate) {
    invalidFields.value.diplomaMetadata.issueDate = t.t(
      'pages.newCredential.form.invalidMessages.issueDate',
    );
    hasErrors = true;
  }
  if (!inputData.value.status) {
    invalidFields.value.status = t.t(
      'pages.newCredential.form.invalidMessages.status',
    );
    hasErrors = true;
  }
  if (!inputData.value.credentialType) {
    invalidFields.value.credentialType = t.t(
      'pages.newCredential.form.invalidMessages.credentialType',
    );
    hasErrors = true;
  }
  if (!diplomaFile.value || diplomaFile.value.length === 0) {
    invalidFields.value.diplomaFile = t.t(
      'pages.newCredential.form.invalidMessages.degreeFile',
    );
    hasErrors = true;
  }

  return hasErrors;
}

async function saveDocument() {
  let res = null;
  try {
    res = await postCredential(inputData.value);
    return res;
  } catch (error) {
    res = error;
  }
  return res;
}

async function formActionClick(id) {
  if (id === 'save') {
    // TODO: add validation
    if (validateModel()) {
      notify.pushError(
        t.t('pages.newCredential.form.invalidMessages.fillAllRequired'),
      );
      return;
    }
    loading.value = true;
    const res = await saveDocument();
    console.log('saveDocument res:', res);
    if (res?.status === 201) {
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

const locale = computed(() => (t?.locale.value === 'lv' ? 'lv-LV' : 'en-EN'));
</script>
<template>
  <div>
    <!-- TODO: add components translations -->
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
      <LxRow :label="t.t('pages.newCredential.form.issuerId')" :required="true">
        <template #info>TODO: get from session</template>
        <LxTextInput v-model="inputData.issuerId" :disabled="loading" />
      </LxRow>
      <LxRow
        :label="t.t('pages.newCredential.form.issuerSignature')"
        :required="true"
      >
        <template #info>TODO: get from session</template>
        <LxTextInput v-model="inputData.issuerSignature" :disabled="loading" />
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
        columnSpan="2"
      >
        <LxTextInput
          v-model="inputData.diplomaMetadata.degreeName"
          :disabled="loading"
          :invalid="invalidFields.diplomaMetadata.degreeName"
          :invalidationMessage="invalidFields.diplomaMetadata.degreeName"
        />
      </LxRow>

      <LxRow :label="t.t('pages.newCredential.form.status')" :required="true">
        <LxValuePicker
          v-model="inputData.status"
          :items="statusItems"
          variant="tags-custom"
          :disabled="loading"
          :invalid="invalidFields.status"
          :invalidationMessage="invalidFields.status"
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
          :invalid="invalidFields.credentialType"
          :invalidationMessage="invalidFields.credentialType"
        />
      </LxRow>

      <LxRow
        :label="t.t('pages.newCredential.form.issueDate')"
        :required="true"
      >
        <LxDateTimePicker
          v-model="inputData.diplomaMetadata.issueDate"
          :disabled="loading"
          :locale="{
            locale: locale,
          }"
          :maxDate="inputData.diplomaMetadata.expiryDate"
          :invalid="invalidFields.diplomaMetadata.issueDate"
          :invalidationMessage="invalidFields.diplomaMetadata.issueDate"
          :texts="t.tm('pages.verification.dateTimePicker')"
        />
      </LxRow>
      <LxRow :label="t.t('pages.newCredential.form.expiryDate')">
        <LxDateTimePicker
          v-model="inputData.diplomaMetadata.expiryDate"
          :disabled="loading"
          :locale="{
            locale: locale,
          }"
          :minDate="inputData.diplomaMetadata.issueDate"
          :texts="t.tm('pages.verification.dateTimePicker')"
        />
      </LxRow>

      <LxRow
        :label="t.t('pages.newCredential.form.degreeFile')"
        :required="true"
        class="file-uploader-row"
      >
        <LxFileUploader
          v-model="diplomaFile"
          :allowedFileExtensions="['.pdf']"
          :disabled="loading"
          :invalid="invalidFields.diplomaFile"
          :invalidationMessage="invalidFields.diplomaFile"
          :texts="t.tm('pages.verification.fileUploader')"
        />
        <p class="lx-invalidation-message" v-if="invalidFields.diplomaFile">
          {{ invalidFields.diplomaFile }}
        </p>
      </LxRow>
    </LxForm>
  </div>
</template>
