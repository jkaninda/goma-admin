<script setup lang="ts">
import { mobileDrawerVisible, closeMobileDrawer } from "@/composables/useLayout";

const config = useRuntimeConfig()
const route = useRoute();

const menuItems = [
  { label: "common.dashboard", icon: "pi pi-home", to: "/" },
  { label: "common.routes", icon: "pi pi-compass", to: "/routes" },
  { label: "common.middlewares", icon: "pi pi-filter", to: "/middlewares" },
  { label: "common.instances", icon: "pi pi-server", to: "/instances" },
  { label: "common.analytics", icon: "pi pi-chart-line", to: "/analytics" },
  { label: "common.configuration", icon: "pi pi-cog", to: "/configuration" },
  { label: "common.history", icon: "pi pi-history", to: "/history" },
];

function isActive(path: string) {
  return route.path === path;
}
</script>

<template>
  <Drawer
    v-model:visible="mobileDrawerVisible"
    position="left"
    class="w-72"
    @hide="closeMobileDrawer"
    :header="config.public.appName"
  >
    <nav class="space-y-1 mt-2">
      <RouterLink
        v-for="item in menuItems"
        :key="item.to"
        :to="item.to"
        @click="closeMobileDrawer"
        class="flex items-center gap-3 px-4 py-3 rounded-md transition"
        :class="isActive(item.to)
          ? 'bg-primary text-primary-contrast'
          : 'hover:bg-primary/10 hover:text-primary'"
      >
        <i :class="item.icon" class="text-lg" />
        <span class="text-sm font-medium">{{ $t(item.label) }}</span>
      </RouterLink>
    </nav>
    
  </Drawer>
</template>
