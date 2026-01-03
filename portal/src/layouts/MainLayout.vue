<script setup>
import { computed, ref, onMounted } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import { useI18n } from 'vue-i18n';
import { LxShell, LxIcon } from '@wntr/lx-ui';
import { invoke, until, useIdle, useIntervalFn } from '@vueuse/core';
import { SYSTEM_NAME } from '@/constants';

import LoginView from '@/views/Login.vue';
import useErrors from '@/hooks/useErrors';
import useAuthStore from '@/stores/useAuthStore';
import useAppStore from '@/stores/useAppStore';
import useNotifyStore from '@/stores/useNotifyStore';
import useConfirmStore from '@/stores/useConfirmStore';
import useViewStore from '@/stores/useViewStore';
import CoverBackground from '@/components/CoverBackground.vue';

const authStore = useAuthStore();
const notify = useNotifyStore();
const viewStore = useViewStore();
const errors = useErrors();
const router = useRouter();
const confirmStore = useConfirmStore();
const appStore = useAppStore();

const secondsToIdle = 10;
const secondsCheckApiInterval = 30;

const { idle } = useIdle(secondsToIdle * 1000);

const idleModalOpened = ref(false);

// ToDo: develop login & get session
// eslint-disable-next-line no-unused-vars

const i18n = useI18n();
const { t } = useI18n();
const route = useRoute();
const routes = router.getRoutes();
const shellMode = computed(() => {
  let ret = 'public';
  if (route.name === 'home') {
    ret = 'cover';
  }
  return ret;
});

const nav = [
  {
    label: i18n.t('pages.dashboard.title'),
    icon: 'dashboard',
    to: { name: 'dashboard' },
  },
  {
    label: i18n.t('pages.newCredential.title'),
    icon: 'add',
    to: { name: 'newCredential' },
  },
  {
    label: i18n.t('pages.verification.title'),
    icon: 'search',
    to: { name: 'verification' },
  },
];

const systemName = computed(() => i18n.t('title.shortName'));

const pageTitle = computed(() => {
  if (typeof route.meta.title === 'function') {
    return route.meta.title(i18n);
  }
  if (typeof route.meta.title === 'string') {
    return viewStore?.pageTitle || i18n.t(route.meta.title);
  }
  return '';
});

const pageDescription = computed(() => {
  if (typeof route.meta.description === 'function') {
    return route.meta.description(i18n);
  }
  if (typeof route.meta.description === 'string') {
    return viewStore?.pageDescription || i18n.t(route.meta.description);
  }
  return '';
});

const breadcrumbs = computed(() => {
  const ret = [];

  if (route.meta.breadcrumbs) {
    // @ts-ignore
    route.meta.breadcrumbs.forEach((item) => {
      ret.push({
        label: viewStore?.backRouteName || i18n.t(item.text),
        to: item.to,
      });
    });
  }
  return ret;
});

const showBackButton = computed(() =>
  viewStore?.canGoBack === false
    ? viewStore.canGoBack
    : breadcrumbs.value.length > 0,
);

const selectedNavItems = computed(() => {
  const ret = {};
  ret[router.currentRoute.value.name] = true;
  if (route.meta?.breadcrumbs) {
    // @ts-ignore
    route.meta?.breadcrumbs.forEach((item) => {
      ret[item.to?.name] = true;
    });
  }
  return ret;
});

const selectedLanguage = ref({ id: 'lv', name: t('shell.languages.lv') });
const languages = computed(() => [
  { id: 'lv', name: t('shell.languages.lv') },
  { id: 'en', name: t('shell.languages.en') },
]);

function changeLang(value) {
  i18n.locale.value = value.id;
  selectedLanguage.value = value;
  localStorage.setItem(`${SYSTEM_NAME}-locale`, value.id);
}

function goBack(path) {
  if (path !== -1) {
    router.push(path);
  } else {
    router.back();
  }
}
function goHome(path) {
  router.push(path);
}

onMounted(() => {
  const res = localStorage.getItem(`${SYSTEM_NAME}-locale`);
  if (res) {
    changeLang(languages.value.find((x) => x.id === res));
  }
  if (authStore.session.active) {
    authStore.keepAlive();
  }
});

const userInfo = computed(() => {
  if (authStore.isAuthorized) {
    return {
      firstName: authStore.session?.given_name,
      lastName: authStore.session?.family_name,
      description: authStore.session?.role,
      institution: null,
    };
  }
  return null;
});

const closeModal = () => {
  idleModalOpened.value = false;
};

const openModal = () => {
  idleModalOpened.value = true;
};

async function logout() {
  try {
    const resp = await authStore.logout();
    if (resp?.status === 200 && resp?.data) {
      window.location.href = resp.data;
    } else {
      notify.pushSuccess(i18n.t('shell.notifications.logOut'));
    }
  } catch (err) {
    const error = errors.get(err);
    if (error.status !== 401 && error.data) {
      notify.pushError(error.data);
    }
  } finally {
    closeModal();
    router.push({ name: 'home' });
  }
}

function primary() {
  logout();
  confirmStore.$state.isOpen = false;
}
function secondary() {
  confirmStore.$state.isOpen = false;
}

function openConfirmModal() {
  confirmStore.push(
    i18n.t('shell.notifications.logoutConfirmTitle'),
    i18n.t('shell.notifications.logoutConfirm'),
    i18n.t('shell.notifications.logoutConfirmYes'),
    i18n.t('shell.notifications.logoutConfirmNo'),
    primary,
    secondary,
  );
}

function confirmModalClosed() {
  confirmStore.$state.isOpen = false;
}

async function getSession() {
  try {
    await authStore.fetchSession();
  } catch (err) {
    const error = errors.get(err);
    if (error.status === 401) {
      logout();
    } else if (error.data) {
      notify.pushError(error.data);
    }
  }
}

async function callKeepAlive() {
  try {
    await authStore.keepAlive();
  } catch (err) {
    const error = errors.get(err);
    if (error.status === 401) {
      logout();
    } else if (error.data) {
      notify.pushError(error.data);
    }
  }
}

const checkApiSession = () => {
  if (idle.value || idleModalOpened.value) {
    getSession();
  } else {
    callKeepAlive();
  }
};

useIntervalFn(() => {
  if (!authStore.session.active) {
    if (idleModalOpened.value) {
      closeModal();
      router.push({ name: 'sessionTimeout' });
    }
    return;
  }
  if (authStore.session.secondsToLive < 1) {
    logout();
    closeModal();
    return;
  }
  if (authStore.session.secondsToLive < authStore.session.secondsToCountdown) {
    if (!idleModalOpened.value) {
      openModal();
    }
  } else if (idleModalOpened.value) {
    closeModal();
    return;
  }
  const refreshIntervals =
    authStore.session.secondsToLive % secondsCheckApiInterval === 0;
  const refreshBeforeWarn =
    authStore.session.secondsToLive - 3 <
      authStore.session.secondsToCountdown && !idle.value;
  const refreshBeforeLogout = authStore.session.secondsToLive === 3;
  if (refreshIntervals || refreshBeforeWarn || refreshBeforeLogout) {
    checkApiSession();
  }
  authStore.session.secondsToLive -= 1;
}, 1000);

async function continueSession() {
  try {
    await authStore.keepAlive();
    notify.pushSuccess(i18n.t('shell.notifications.sessionContinued'));
  } catch (err) {
    notify.pushError(i18n.t('shell.notifications.sessionContinuedFailed'));
    if (err.response.status === 401) {
      logout();
    }
  } finally {
    closeModal();
  }
}

invoke(async () => {
  // @ts-ignore
  await until(() => authStore.showSessionEndCountdown).toBe(true);
  notify.pushWarning(i18n.t('shell.notifications.sessionEndingSoon'));
});

function idleModalPrimary() {
  continueSession();
}
function idleModalSecondary() {
  logout();
}
</script>
<template>
  <div>
    <div>
      <LxShell
        :system-name="i18n.t('title.fullName')"
        :system-subheader="i18n.t('title.subheader')"
        :system-name-short="systemName"
        :user-info="userInfo"
        :nav-items="nav"
        :nav-items-selected="selectedNavItems"
        :mode="shellMode"
        :page-label="pageTitle"
        :pageDescription="pageDescription"
        :page-back-button-visible="showBackButton"
        :page-breadcrumbs="breadcrumbs"
        :page-index-path="{ name: 'home' }"
        :has-cover-logo="false"
        :cover-image="null"
        :cover-image-dark="null"
        :cover-logo="null"
        :has-theme-picker="true"
        :has-language-picker="true"
        :languages="languages"
        :selected-language="selectedLanguage"
        :navigating="appStore.$state.isNavigating"
        :showIdleModal="idleModalOpened"
        :showIdleBadge="
          authStore.session.secondsToLive <
            authStore.session.secondsToCountdown &&
          !authStore.session.isSessionExtendable
        "
        :secondsToLive="authStore.session.secondsToLive"
        :confirmDialogData="confirmStore"
        :confirmPrimaryButtonBusy="false"
        :confirmPrimaryButtonDestructive="true"
        v-model:notifications="notify.notifications"
        :hideNavBar="!viewStore?.isNavBarShown"
        :headerNavDisable="viewStore.blockNav"
        :hideHeaderText="!viewStore?.isHeaderShown"
        @update:selected-language="changeLang"
        @confirmModalClosed="confirmModalClosed"
        @go-home="goHome"
        @go-back="goBack"
        @log-out="openConfirmModal"
        @idleModalPrimary="idleModalPrimary"
        @idleModalSecondary="idleModalSecondary"
      >
        <template #backdrop>
          <CoverBackground />
        </template>
        <template #coverArea>
          <div class="lx-button-set">
            <LoginView />
          </div>
        </template>
        <template #footer>
          <div class="footer-sitemap">
            <div class="footer-sitemap-left">
              <div class="footer-logo">
                <LxIcon value="link" />
                <div>{{ i18n.t("title.fullName") }}</div>
              </div>
              <div class="footer-sitemap-text">
                {{ i18n.t("title.subheader") }}
              </div>
            </div>
            <div class="footer-sitemap-right">
              <p class="footer-sitemap-text">
                {{ i18n.t("shell.footer.footerText") }}
              </p>
              <p class="footer-sitemap-text">
                {{ i18n.t("shell.footer.versionText") }} 0.1.11 20.12.2022
              </p>
            </div>
          </div>
        </template>
        <template #logo>
          <!-- TODO: add better -->
          <LxIcon value="link" />
        </template>

        <template #logoSmall>
          <!-- TODO: add better -->
          <LxIcon value="link" />
        </template>
        <router-view />
      </LxShell>
    </div>
  </div>
</template>
