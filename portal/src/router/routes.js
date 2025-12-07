/**
 * @typedef {Object} CustomMetaProps
 * @property {string} [title]
 * @property {boolean} [canGoBack] - if true, the route will have a back button
 * @property {string} [backRouteName] - name of the route to go back to
 * @property {boolean} [anonymous] - if true, the route is accessible also for anonymous users
 * @property {boolean} [onlyAnonymous] - if true, the route is accessible only for anonymous users
 * @property {{ text: string, to: import('vue-router').RouteLocationRaw }[]} [breadcrumbs] - array of breadcrumbs
 * @property {(rights: ReturnType<typeof import('@/hooks/useRights').default>) => boolean} [access] - route access function that returns true if the user has access to the route
 */
/**
 * @typedef {import('vue-router').RouteRecordRaw & { meta?: CustomMetaProps, children?: CustomRoute[] }} CustomRoute
 */

/** @type {CustomRoute[]} */
const routes = [
  {
    path: '/',
    name: 'home',
    component: () => import('@/layouts/MainLayout.vue'),
    meta: {
      title: 'pages.home.title',
    },
    children: [
      {
        path: '/dashboard',
        name: 'dashboard',
        meta: {
          title: 'pages.dashboard.title',
          // access: (rights) => rights.canViewDashboard(),
          category: 'useful',
        },
        component: () => import('@/views/Dashboard.vue'),
      },
      {
        path: '/auth-done',
        name: 'authDone',
        component: () => import('@/views/AuthDone.vue'),
        meta: {
          title: 'pages.auth.title',
          category: 'useful',
          anonymous: true,
        },
      },
      {
        path: '/sessionTimeout',
        name: 'sessionTimeout',
        component: () => import('@/views/SessionTimeout.vue'),
        meta: {
          title: 'pages.sessionTimeout.title',
          category: 'useful',
          anonymous: true,
        },
      },
      {
        path: '/error',
        name: 'error',
        alias: 'notFound',
        component: () => import('@/views/Error404.vue'),
        meta: {
          anonymous: true,
          category: 'useful',
          hideFromIndex: true,
        },
      },
      {
        path: '/forbidden',
        name: 'forbidden',
        component: () => import('@/views/Error403.vue'),
        meta: { anonymous: true },
      },
      {
        path: '/not-authorized',
        name: 'notAuthorized',
        component: () => import('@/views/Error401.vue'),
        meta: { anonymous: true },
      },
    ],
  },
];

export default routes;
