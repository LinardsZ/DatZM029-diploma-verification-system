/* eslint-disable no-console -- for debugging */
/* eslint-disable no-restricted-imports -- vite doesn`t use aliases in imports */
/* eslint-disable import/no-extraneous-dependencies -- vite mostly should use dev dependencies */
import { defineConfig, loadEnv } from 'vite';
import vue from '@vitejs/plugin-vue';
import path from 'path';
import { createHtmlPlugin } from 'vite-plugin-html';
import mkcert from 'vite-plugin-mkcert';
import { fileURLToPath } from 'url';
import dns from 'dns';
import viteCompression from 'vite-plugin-compression';
import { lxViteSecureHeadersPlugin } from '@wntr/lx-ui/vite';
import packageJson from './package.json';
import { devServerSettings } from './vite.config.server.mjs';

// set ip4 as default dns lookup in order to support localhost
// https://vitejs.dev/config/server-options.html#server-host
dns.setDefaultResultOrder('ipv4first');

const CONFIG = {
  appName: '$REPO_NAME',
  appDescription: '$REPO_DESCRIPTION',
  defaultUrl: 'https://localhost:44341/',
  defaultClient: '$REPO_NAME_LOWER',
};

const getEnvVariables = (mode, serving) => {
  // Load env file based on `mode` in the current working directory.
  // Set the third parameter to '' to load all env regardless of the `VITE_` prefix.
  const env = loadEnv(mode, process.cwd(), '');
  const envVariables = {
    VUE_APP_VERSION: env.CUSTOM_VERSION || packageJson.version,
    VUE_APP_ENVIRONMENT: serving ? 'local' : env.NODE_ENV,
    VUE_APP_NAME: CONFIG.appName,
    VUE_APP_DESCRIPTION: CONFIG.appDescription,
    VUE_APP_SERVICE_URL: env.SERVICE_URL,
    VUE_APP_AUTH_URL: env.AUTH_URL,
    VUE_APP_CLIENT_ID: env.CLIENT_ID,
    BASE_PATH: env.BASE_PATH,
    BASE_URL: env.PUBLIC_URL,

    // dev only
    VUE_APP_SERVICE_URL_PROXY: '',
    VUE_APP_AUTH_URL_PROXY: '',
  };
  if (serving) {
    envVariables.BASE_URL = env.PUBLIC_URL || CONFIG.defaultUrl;
    envVariables.BASE_PATH = envVariables.BASE_PATH || '/';
    envVariables.VUE_APP_VERSION += '-serve';
    envVariables.VUE_APP_SERVICE_URL = '/api/';
    envVariables.VUE_APP_SERVICE_URL_PROXY = env.SERVICE_URL;
    envVariables.VUE_APP_AUTH_URL = '/idauth/';
    envVariables.VUE_APP_AUTH_URL_PROXY = env.AUTH_URL;
    envVariables.USE_MOCK_MIDDLEWARE = env.USE_MOCK_MIDDLEWARE === 'true';
    envVariables.VUE_APP_CLIENT_ID = envVariables.VUE_APP_CLIENT_ID || CONFIG.defaultClient;
  } else {
    envVariables.VUE_APP_ENVIRONMENT = '{{ENVIRONMENT}}';
    envVariables.VUE_APP_SERVICE_URL = envVariables.VUE_APP_SERVICE_URL || '{{SERVICE_URL}}';
    envVariables.VUE_APP_AUTH_URL = envVariables.VUE_APP_AUTH_URL || '{{AUTH_URL}}';
    envVariables.BASE_PATH = envVariables.BASE_PATH || '//BASE_PATH//';
    envVariables.BASE_URL = envVariables.BASE_URL || '{{PUBLIC_URL}}';
    envVariables.VUE_APP_CLIENT_ID = envVariables.VUE_APP_CLIENT_ID || CONFIG.defaultClient;
  }
  return envVariables;
};

/**
 * @returns { import('vite').UserConfigExport }
 * @description https://vitejs.dev/config/
 */
export default defineConfig((command) => {
  const serving = command?.command === 'serve' && command?.mode === 'development';
  const envVariables = getEnvVariables(command.mode, serving);
  return {
    base: envVariables.BASE_PATH,
    resolve: {
      alias: {
        '@': fileURLToPath(new URL('./src', import.meta.url)),
        '/lx-fonts': fileURLToPath(
          new URL('./node_modules/@wntr/lx-ui/dist/lx-fonts', import.meta.url)
        ),
      },
    },
    plugins: [
      vue(),
      createHtmlPlugin({
        minify: true,
        inject: {
          data: envVariables,
        },
      }),
      serving && mkcert(),
      viteCompression({
        algorithm: 'gzip',
        ext: '.gz',
      }),
      lxViteSecureHeadersPlugin(),
    ],
    build: {
      // https://vitejs.dev/config/#build-target
      target: ['es2020'],
      outDir: './dist',
      sourcemap: true, // in prod env you can build with `pnpm run build --sourcemap=false` in order to disable sourcemaps
    },
    server: serving ? devServerSettings(envVariables) : {},
    test: {
      globals: true,
      setupFiles: [path.resolve(__dirname, './tests/setup.js')],
      environment: 'jsdom',
      isolate: false,
      include: ['tests/unit/**/*.js'],
    },
  };
});
