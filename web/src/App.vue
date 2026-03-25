<template>
  <router-view />
  <ToastContainer />
  <ConfirmDialog />
</template>

<script setup lang="ts">
import { onMounted } from 'vue'
import { useThemeStore } from '@/stores/theme'
import ToastContainer from '@/components/ToastContainer.vue'
import ConfirmDialog from '@/components/ConfirmDialog.vue'

const theme = useThemeStore()

onMounted(() => {
  // Apply saved theme on startup
  theme.setMode(theme.mode)
})
</script>
