<script setup lang="ts">
import { sidebarOpen } from "@/composables/useLayout";
const config = useRuntimeConfig()

const year = new Date().getFullYear();

const footerClasses = computed(() => ({
  "lg:ml-56 lg:w-[calc(100%-14rem)]": sidebarOpen.value,
  "lg:ml-16 lg:w-[calc(100%-4rem)]": !sidebarOpen.value,
}));
</script>

<template>
  <footer
    class="fixed bottom-0 left-0 right-0 h-12 px-4 flex items-center justify-between bg-surface-0 dark:bg-surface-900 border-t border-surface-200 dark:border-surface-800 transition-all duration-300"
    :class="footerClasses"
  >
    <div class="text-sm text-surface-600 dark:text-surface-400">
      © {{ year +" "+config.public.appName}}
    </div>

    <div class="hidden md:flex text-xs text-surface-500">{{ config.public.appVersion }}</div>

    <div class="flex gap-1 text-xs">
      <Button :label="$t('common.docs')" variant="link" size="small" />
      <Button :label="$t('common.support')" variant="link" size="small" />
      <Button :label="$t('common.github')" variant="link" size="small" />
    </div>
  </footer>
</template>
