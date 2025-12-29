import { createApp } from 'vue';
import { createPinia } from 'pinia';
import VueClickAway from 'vue3-click-away';
import App from '@/App.vue';
import router from '@/router';
import events from '@/router/events';
import i18n from '@/i18n';
import { createLx } from '@wntr/lx-ui';
import { APP_CONFIG, AUTH_KEY_TOKEN_SESSION, SYSTEM_NAME } from '@/constants';

import '@wntr/lx-ui/dist/styles/lx-reset.css';
import '@wntr/lx-ui/dist/styles/lx-fonts-carbon.css';
import '@wntr/lx-ui/dist/styles/lx-pt-carbon.css';
import '@wntr/lx-ui/dist/styles/lx-ut-carbon-light.css';
import '@wntr/lx-ui/dist/styles/lx-ut-carbon-dark.css';
import '@wntr/lx-ui/dist/styles/lx-ut-carbon-contrast.css';

import '@wntr/lx-ui/dist/styles/lx-buttons.css';
import '@wntr/lx-ui/dist/styles/lx-data-grid.css';
import '@wntr/lx-ui/dist/styles/lx-inputs.css';
import '@wntr/lx-ui/dist/styles/lx-steps.css';
import '@wntr/lx-ui/dist/styles/lx-forms.css';
import '@wntr/lx-ui/dist/styles/lx-notifications.css';
import '@wntr/lx-ui/dist/styles/lx-modal.css';
import '@wntr/lx-ui/dist/styles/lx-loaders.css';
import '@wntr/lx-ui/dist/styles/lx-lists.css';
import '@wntr/lx-ui/dist/styles/lx-expanders.css';
import '@wntr/lx-ui/dist/styles/lx-tabs.css';
import '@wntr/lx-ui/dist/styles/lx-animations.css';
import '@wntr/lx-ui/dist/styles/lx-master-detail.css';
import '@wntr/lx-ui/dist/styles/lx-ratings.css';
import '@wntr/lx-ui/dist/styles/lx-day-input.css';
import '@wntr/lx-ui/dist/styles/lx-map.css';
import '@wntr/lx-ui/dist/styles/lx-shell-grid.css';
import '@wntr/lx-ui/dist/styles/lx-shell-grid-public.css';
import '@wntr/lx-ui/dist/styles/lx-forms-grid.css';
import '@wntr/lx-ui/dist/styles/lx-treelist.css';
import '@wntr/lx-ui/dist/styles/lx-date-pickers.css';
import '@wntr/lx-ui/dist/styles/lx-data-visualizer.css';
import '@wntr/lx-ui/dist/styles/lx-stack.css';

import '@/assets/styles.css';

const myApp = createApp(App);
myApp.use(createPinia());
events(router);
myApp.use(router);

myApp.use(i18n);
myApp.use(VueClickAway);
myApp.use(createLx, {
  systemId: SYSTEM_NAME,
  authSessionKey: AUTH_KEY_TOKEN_SESSION,
  authUrl: APP_CONFIG.authUrl,
  authClientId: APP_CONFIG.clientId,
  publicUrl: APP_CONFIG.publicUrl,
  environment: APP_CONFIG.environment,
  iconSet: 'phosphor',
});
myApp.mount('#app');
