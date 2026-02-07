export default defineNuxtRouteMiddleware((to, from) => {
    const auth = useAuthStore();


    if (auth.isLoggedIn) {

        const redirect = to.query.redirect as string | undefined;


        if (redirect) {
            return navigateTo(redirect);
        } else if (from.fullPath && from.fullPath !== to.fullPath) {
            return navigateTo(from.fullPath);
        }

        return navigateTo("/dashboard");
    }
})