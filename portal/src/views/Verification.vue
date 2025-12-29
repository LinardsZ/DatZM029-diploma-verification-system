<script setup>
import { ref, watch } from 'vue';
import { LxFileUploader, LxButton, LxInfoBox } from '@wntr/lx-ui';
import { getFileHash } from '@/utils/generalUtils';
import { useI18n } from 'vue-i18n';

const t = useI18n();

const file = ref();
const fileHash = ref();

const success = ref(false);
const error = ref(false);
const loading = ref(false);

watch(file, async (newValue) => {
  if (newValue && newValue.length > 0) {
    const res = await getFileHash(newValue);

    fileHash.value = res;
  } else {
    fileHash.value = null;
  }
});

function verifyCall() {
  return new Promise((resolve) => {
    setTimeout(() => {
      resolve({
        status: 300,
      });
    }, 2000);
  });
}

async function verifyFile() {
  loading.value = true;
  try {
    const res = await verifyCall();

    if (res.status === 200) {
      success.value = true;
    } else {
      error.value = true;
    }
  } catch (e) {
    console.error(e);
  } finally {
    loading.value = false;
  }
}

function reset() {
  success.value = false;
  error.value = false;
  file.value = null;
}
</script>

<template>
  <div>
    <LxFileUploader
      v-if="!success && !error"
      v-model="file"
      :disabled="loading"
      :allowedFileExtensions="['.pdf']"
      :texts="t.tm('pages.verification.fileUploader')"
    />
    <LxButton
      v-if="!success && !error"
      customClass="verify-button"
      :label="t.t('pages.verification.verify')"
      icon="search"
      :disabled="!fileHash"
      :busy="loading"
      @click="verifyFile"
    />
    <LxInfoBox
      v-if="success || error"
      :label="
        success
          ? t.t('pages.verification.successLabel')
          : t.t('pages.verification.errorLabel')
      "
      :description="
        success
          ? t.t('pages.verification.successDescription')
          : t.t('pages.verification.errorDescription')
      "
      :variant="success ? 'success' : 'error'"
      kind="button"
      :actionDefinitions="[
        { id: 'reset', name: t.t('pages.verification.reset'), icon: 'reset' },
      ]"
      @actionClick="reset"
    />
  </div>
</template>
