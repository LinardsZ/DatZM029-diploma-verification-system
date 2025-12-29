/* eslint-disable no-restricted-imports */
import { URL } from 'url';
import mocks from './mocks/middleware.mock';

const createProxyConfig = (env, baseUrl) => (proxy) => {
  proxy.on('proxyReq', async (proxyReq, req, res) => {
    const url = new URL(baseUrl);
    const basePath = url.pathname;
    const baseUrlWithoutPath = url.origin;
    const parsedUrl = new URL(basePath + req.url.replace(/^\//, ''), baseUrlWithoutPath);
    // eslint-disable-next-line no-console
    console.log(new Date().toUTCString(), 'request', req.method, parsedUrl.href);
    // ToDo: set conditions according to your project in order to correctly set cookies. By default, we check if request is for `/idauth/authorize` endpoint
    const shouldFollowWithoutProxy =
      req.url.includes('/authorize') && basePath.includes('/idauth/');
    if (env.USE_MOCK_MIDDLEWARE) {
      const mockHandler = mocks.find(
        (mock) => mock.pattern === parsedUrl.pathname && mock.method === req.method
      );
      if (mockHandler) {
        await mockHandler.handle(req, res);
        proxyReq.destroy({ stack: `mocked response for req: ${req.url}` });
      }
    } else if (shouldFollowWithoutProxy) {
      // skip proxy for authorize request, in order to correctly set cookies
      const authUrl = new URL(env.VUE_APP_AUTH_URL_PROXY);
      const targetUrl = new URL(basePath + req.url.replace(/^\//, ''), authUrl.origin);
      res.writeHead(302, { Location: targetUrl.toString() });
      res.end();
    }
  });
};

/**
 * @param { ReturnType<getEnvVariables> } env
 * @returns { import('vite').ServerOptions }
 */
export const devServerSettings = (env) => {
  const url = new URL(env.BASE_URL);
  return {
    port: url.port,
    https: url.protocol === 'https:',
    proxy: {
      '/api': {
        target: env.VUE_APP_SERVICE_URL_PROXY,
        changeOrigin: true,
        secure: false,
        rewrite: (p) => p.replace(/^\/api/, ''),
        configure: createProxyConfig(env, env.VUE_APP_SERVICE_URL_PROXY),
      },
    },
    fs: {
      allow: [
        // Allow serving files from one level up to the project root, in order to show lx fonts in serving mode
        '..',
        // Allow serving files from two levels up to the project root, in order to show lx fonts in serving mode when lx lib is referenced in the project
        '../..',
      ],
    },
  };
};
