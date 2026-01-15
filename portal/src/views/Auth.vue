<script setup>
import { useI18n } from 'vue-i18n';
import { onMounted, ref } from 'vue';
import { LxForm, LxRow, LxTextInput } from '@wntr/lx-ui';
import useAuthStore from '@/stores/useAuthStore';
import api from '@/api';
import router from '@/router';
import useNotification from '@/stores/useNotifyStore';

const authStore = useAuthStore();
const notification = useNotification();
const t = useI18n();

const loading = ref(false);
const data = ref({
  username: '',
  password: '',
});

const invalidDefault = {
  username: null,
  password: null,
};

const invalidMessages = ref({ ...invalidDefault });

function validate() {
  invalidMessages.value = { ...invalidDefault };
  if (data.value.username === '') {
    invalidMessages.value.username = t.t('pages.auth.usernameEmpty');
    return false;
  }
  if (data.value.password === '') {
    invalidMessages.value.password = t.t('pages.auth.passwordEmpty');
    return false;
  }
  return true;
}

function login() {
  if (!validate()) {
    notification.pushError(t.t('pages.auth.errors.invalidRequestTitle'));
    return;
  }
  loading.value = true;
  const resp = api().post('/auth/login', data.value);
  resp
    .then((response) => {
      if (response.status !== 200) {
        return;
      }
      authStore.session.st = 'authorized';
      authStore.session.given_name = response.data.firstName;
      authStore.session.family_name = response.data.lastName;
      authStore.session.institution = response.data.issuer;
      router.push({ name: 'dashboard' });
    })
    .catch((error) => {
      notification.pushError(t.t('pages.auth.authError'));
    });
  resp.finally(() => {
    loading.value = false;
  });
}
</script>

<template>
  <LxForm
    :showHeader="false"
    :sticky-header="false"
    :actionDefinitions="[
      {
        id: 'login',
        name: t.t('pages.auth.loginButton'),
        kind: 'primary',
        icon: 'login',
        busy: loading,
      },
    ]"
    @buttonClick="login()"
  >
    <LxRow :label="t.t('pages.auth.username')">
      <LxTextInput
        v-model="data.username"
        :placeholder="t.t('pages.auth.usernamePlaceholder')"
        :invalid="invalidMessages.username"
        :invalidation-message="invalidMessages.username"
        @keyup.enter="login()"
      />
    </LxRow>
    <LxRow :label="t.t('pages.auth.password')">
      <LxTextInput
        type="password"
        v-model="data.password"
        :placeholder="t.t('pages.auth.passwordPlaceholder')"
        kind="password"
        :invalid="invalidMessages.password"
        :invalidation-message="invalidMessages.password"
        @keyup.enter="login()"
      />
    </LxRow>
  </LxForm>
</template>
