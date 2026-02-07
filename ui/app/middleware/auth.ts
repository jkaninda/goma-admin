export default defineNuxtRouteMiddleware((to, from) => {
    const authStore = useAuthStore();
    if (!authStore.isLoggedIn) {
        return navigateTo(`/auth/login?redirect=${to.fullPath}`);
    }
});
