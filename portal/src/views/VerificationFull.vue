<script setup>
import { computed, ref, watch } from 'vue';
import {
  LxFileUploader,
  LxButton,
  LxInfoBox,
  LxForm,
  LxRow,
  LxSteps,
  LxTextInput,
  lxDateUtils,
  LxTextArea,
} from '@wntr/lx-ui';
import { getFileHash } from '@/utils/generalUtils';
import { useI18n } from 'vue-i18n';
import {
  verifyDiplomaHash,
  verifyDiplomaSignature,
} from '@/services/credentialService';
import useNotifyStore from '@/stores/useNotifyStore';

const t = useI18n();
const notification = useNotifyStore();

const file = ref();
const fileHash = ref();

const success = ref(false);
const error = ref(false);
const loading = ref(false);

const stepModel = ref('default');
const graduateSignature = ref(
  'MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA1234...',
);
const credentialId = ref(null);
const fullDiplomaData = ref(null);

function getState(stepId) {
  if (stepModel.value === 'default') {
    if (stepId === 'default') {
      return 'current';
    }
    return null;
  }
  if (stepModel.value === 'sign') {
    if (stepId === 'default') {
      return 'complete';
    }
    if (stepId === 'sign') {
      return 'current';
    }
    return null;
  }
  if (stepId !== 'view') {
    return 'complete';
  }
  return 'current';
}

const index = computed(() => [
  {
    code: 'default',
    name: t.t('pages.verificationFull.form.upload'),
    state: getState('default'),
  },
  {
    code: 'sign',
    name: t.t('pages.verificationFull.form.sign'),
    state: getState('sign'),
  },
  {
    code: 'view',
    name: t.t('pages.verificationFull.form.view'),
    state: getState('view'),
  },
]);

watch(file, async (newValue) => {
  if (newValue && newValue.length > 0) {
    const res = await getFileHash(newValue);

    fileHash.value = res;
  } else {
    fileHash.value = null;
  }
});

function verifyCall() {
  return verifyDiplomaHash({
    diplomaHash: fileHash.value,
  });
}

async function verifyFile() {
  loading.value = true;
  try {
    const res = await verifyCall();
    console.log('verifyCall res:', res);
    if (res.status === 200) {
      credentialId.value = res.data.credentialId;
      success.value = true;
      stepModel.value = 'sign';
      notification.pushSuccess(t.t('pages.verification.successLabel'));
    } else {
      error.value = true;
    }
  } catch (e) {
    console.error(e);
    error.value = true;
  } finally {
    loading.value = false;
  }
}

function reset() {
  success.value = false;
  error.value = false;
  file.value = null;
}

function resetFull() {
  success.value = false;
  error.value = false;
  file.value = null;
  stepModel.value = 'default';
  graduateSignature.value =
    'MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA1234...';
  credentialId.value = null;
  fullDiplomaData.value = null;
}

async function verifyFull() {
  loading.value = true;
  try {
    const res = await verifyDiplomaSignature({
      credentialId: credentialId.value,
      graduateSignature: graduateSignature.value,
    });
    fullDiplomaData.value = res.data.credential;
    stepModel.value = 'view';
    notification.pushSuccess('Dokumenta paraksts veiksmīgi pārbaudīts.');
  } catch (e) {
    console.error(e);
    notification.pushError('Kļūda pārbaudot dokumenta parakstu.');
  } finally {
    loading.value = false;
  }
}
</script>

<template>
  <div>
    <LxForm :showHeader="false" :columnCount="2">
      <LxRow columnSpan="2" :hideLabel="true">
        <LxSteps v-model="stepModel" :items="index" />
        <div v-if="stepModel === 'default'">
          <LxFileUploader
            v-if="!success && !error"
            v-model="file"
            :disabled="loading"
            :allowedFileExtensions="['.pdf']"
            :texts="t.tm('pages.verification.fileUploader')"
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
              {
                id: 'reset',
                name: t.t('pages.verification.reset'),
                icon: 'reset',
              },
            ]"
            @actionClick="reset"
          />
        </div>
        <div v-if="stepModel === 'sign'">
          Luuudzu paraksti mani:
          <p>{{ credentialId }}</p>
          <LxTextArea v-model="graduateSignature" />
        </div>
      </LxRow>
      <template v-if="stepModel === 'view'">
        <LxRow label="id">
          <p class="lx-data">{{ fullDiplomaData.id || '—' }}</p>
        </LxRow>
        <LxRow label="diplomaHash">
          <p class="lx-data">{{ fullDiplomaData.diplomaHash || '—' }}</p>
        </LxRow>
        <LxRow label="graduatePublicKey">
          <p class="lx-data">{{ fullDiplomaData.graduatePublicKey || '—' }}</p>
        </LxRow>
        <LxRow label="issuerId">
          <p class="lx-data">{{ fullDiplomaData.issuerId || '—' }}</p>
        </LxRow>
        <LxRow label="issuerSignature">
          <p class="lx-data">{{ fullDiplomaData.issuerSignature || '—' }}</p>
        </LxRow>
        <LxRow label="universityName">
          <p class="lx-data">
            {{ fullDiplomaData.diplomaMetadata.universityName || '—' }}
          </p>
        </LxRow>
        <LxRow label="degreeName">
          <p class="lx-data">
            {{ fullDiplomaData.diplomaMetadata.degreeName || '—' }}
          </p>
        </LxRow>
        <LxRow label="issueDate">
          <p class="lx-data">
            {{
              lxDateUtils.formatDate(
                fullDiplomaData.diplomaMetadata.issueDate,
              ) || '—'
            }}
          </p>
        </LxRow>
        <LxRow label="expiryDate">
          <p class="lx-data">
            {{
              lxDateUtils.formatDate(
                fullDiplomaData.diplomaMetadata.expiryDate,
              ) || '—'
            }}
          </p>
        </LxRow>
        <LxRow label="status">
          <p class="lx-data">{{ fullDiplomaData.status || '—' }}</p>
        </LxRow>
        <LxRow label="credentialType">
          <p class="lx-data">{{ fullDiplomaData.credentialType || '—' }}</p>
        </LxRow>
      </template>
      <template #footer>
        <div class="form-custom-footer">
          <LxButton
            v-if="!success && !error"
            :label="t.t('pages.verification.verify')"
            icon="search"
            :disabled="!fileHash"
            :busy="loading"
            @click="verifyFile"
          />
          <LxButton
            v-if="stepModel === 'sign'"
            :label="t.t('pages.verificationFull.form.verifyOwner')"
            icon="inspection"
            :disabled="stepModel !== 'sign' || !graduateSignature"
            @click="verifyFull"
          />

          <LxButton
            v-if="stepModel === 'view'"
            :label="t.t('pages.verificationFull.form.reset')"
            icon="reset"
            @click="resetFull"
          />
        </div>
      </template>
    </LxForm>
  </div>
</template>
