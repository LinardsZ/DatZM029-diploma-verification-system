import useAuthStore from '@/stores/useAuthStore';

export default () => {
  const authStore = useAuthStore();
  // ToDo: add your scopes and check permissions here
  // const authScopes = {
  //   /** Game list scope, just example scope in order to test permissions */
  //   gameList: 'game/list',
  // };
  // canViewDashboard: () => authStore.session?.scope.some((scope) => scope.includes(authScopes.gameList))
  return {
    // Temporary function to test permissions
    canViewDashboard: () => authStore.session?.active,
  };
};
