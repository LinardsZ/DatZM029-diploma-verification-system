import useAuthStore from '@/stores/useAuthStore';
import useAppStore from '@/stores/useAppStore';
import { lxFlowUtils } from '@wntr/lx-ui';
import useRights from '@/hooks/useRights';

/**
 * @param { import('vue-router').RouteRecordNormalized & { meta: import('@/router/routes').CustomMetaProps } } record
/* @param {ReturnType<typeof import('@/hooks/useRights').default>} rights
*/
function checkRouteAccess(record, rights) {
  if (typeof record.meta.access === 'function') {
    return record.meta.access(rights);
  }
  if (record.meta.access === undefined) {
    return true;
  }
  throw new Error('Invalid access property in route');
}
/**
 * @param { import('vue-router').RouteLocationNormalized } to
 * @param { import('vue-router').RouteLocationNormalized } from
 * @param { import('vue-router').NavigationGuardNext} next
 */
function routeCheckCallback(to, from, next) {
  const rights = useRights();
  const onlyAnonymous = to.matched.some((record) => record.meta.onlyAnonymous);
  const routesWithAccessControl = to.matched.filter((route) => route.meta.access);
  const canAccessRoute =
    routesWithAccessControl.length === 0 ||
    routesWithAccessControl.some((route) => checkRouteAccess(route, rights));
  if (canAccessRoute && !onlyAnonymous) {
    next();
    return;
  }
  next({
    params: { pathMatch: to.path.substring(1) },
    query: { returnPath: to.path },
    replace: true,
    name: 'forbidden',
  });
}

export default (router) => {
  router.beforeEach(async (to, from, next) => {
    const authStore = useAuthStore();
    const appStore = useAppStore();
    await lxFlowUtils.beforeEach(to, from, next, appStore, authStore, routeCheckCallback);
  });
  router.afterEach(async (to, from) => {
    const appStore = useAppStore();
    await lxFlowUtils.afterEach(to, from, appStore);
    lxFlowUtils.removeFocus();
  });
};
