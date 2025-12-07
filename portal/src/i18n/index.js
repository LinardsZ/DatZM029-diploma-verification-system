import { createI18n } from 'vue-i18n';
import lv from '@/locales/lv.json';
import { APP_CONFIG } from '@/constants';

const i18n = createI18n({
  locale: APP_CONFIG.defaultLocale,
  fallbackLocale: APP_CONFIG.fallbackLocale,
  legacy: false,
  messages: {
    lv,
  },
});

export default i18n;
