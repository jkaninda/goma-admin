export const sidebarOpen = ref(true);
export const mobileDrawerVisible = ref(false);

export function toggleSidebar() {
    sidebarOpen.value = !sidebarOpen.value;
}

export function openMobileDrawer() {
    mobileDrawerVisible.value = true;
}

export function closeMobileDrawer() {
    mobileDrawerVisible.value = false;
}
