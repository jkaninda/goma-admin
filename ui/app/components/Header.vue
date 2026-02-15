<script setup lang="ts">
import { toggleSidebar, openMobileDrawer } from "@/composables/useLayout";
const { t } = useI18n();


const userMenu = ref();
const notifications = ref(3);

const userItems = [
  { label: t("common.profile"), icon: "pi pi-user" },
  { label: t("common.settings"), icon: "pi pi-cog" },
  { separator: true },
  { label: t("common.logout"), icon: "pi pi-sign-out" },
];

function toggleUserMenu(event: Event) {
  userMenu.value.toggle(event);
}

function handleToggleSidebar() {
  if (window.innerWidth < 1024) {
    openMobileDrawer();
  } else {
    toggleSidebar();
  }
}


</script>

<template>
  <header class="fixed top-0 left-0 right-0 z-50 h-10 shadow-sm">
    <Toolbar unstyled class="h-full flex items-center justify-between px-3 bg-primary">
      <template #start>
        <div class="flex items-center gap-2">
          <Button
            icon="pi pi-bars"
            text
            rounded
            size="small"
            aria-label="common.toggle_sidebar"
            @click="handleToggleSidebar "
            class="lg:hidden text-primary-contrast! hover:bg-white/10!"
          />

          <img src="/favicon.ico" class="w-5 h-5 opacity-90" alt="logo" />
        </div>
      </template>

      <template #end>
        <div class="flex items-center gap-1 sm:gap-2">
          <ThemeToggle variant="toolbar" />
          <LanguageSwitcher />

          <Button
            text
            rounded
            size="small"
            aria-label="common.notifications"
            class="text-primary-contrast! hover:bg-white/10!"
          >
            <span class="relative inline-flex">
              <i class="pi pi-bell text-base" />

              <Badge
                v-if="notifications"
                :value="notifications"
                severity="danger"
                class="absolute -top-1.5 -right-1 text-[8px]"
              />
            </span>
          </Button>

          <Button
            text
            rounded
            size="small"
            aria-label="common.user_menu"
            aria-haspopup="true"
            aria-controls="user_menu"
            @click="toggleUserMenu"
            class="text-primary-contrast! hover:bg-white/10!"
          >
            <i class="pi pi-user text-base" />
          </Button>

          <Menu ref="userMenu" :model="userItems" id="user_menu" :popup="true" />
        </div>
      </template>
    </Toolbar>
  </header>
</template>
