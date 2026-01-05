<script setup>
import { useI18n } from 'vue-i18n';
import { onMounted, ref } from 'vue';
import { LxForm, LxRow, LxTextInput } from '@wntr/lx-ui';
import useAuthStore from '@/stores/useAuthStore';
import api from '@/api';
import router from '@/router';
import useViewStore from '@/stores/useViewStore';

const authStore = useAuthStore();
const i18n = useI18n();

const data = ref({
  username: '',
  password: '',
});

function login() {
  const resp = api().post('/auth/login', data.value);
  resp.then((response) => {
    if (response.status !== 200) {
      return;
    }
    authStore.session.st = 'authorized';
    authStore.session.given_name = response.data.firstName;
    authStore.session.family_name = response.data.lastName;
    authStore.session.institution = response.data.issuer;
    router.push({ name: 'dashboard' });
  });
}

onMounted(() => {
  useViewStore().setPageTitle('Autentifikācija');
});
</script>

<template>
  <LxForm
    :show-header="false"
    :action-definitions="[{ id: 'login', name: 'Login', kind: 'primary' }]"
    @button-click="login()"
  >
    <LxRow label="Username">
      <LxTextInput v-model="data.username" placeholder="Enter your username" />
    </LxRow>
    <LxRow label="Password">
      <LxTextInput
        type="password"
        v-model="data.password"
        placeholder="Enter your password"
        kind="password"
      />
    </LxRow>
  </LxForm>
</template>
