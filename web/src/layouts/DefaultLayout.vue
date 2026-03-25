<template>
  <div class="min-h-screen bg-gray-50 dark:bg-gray-900">
    <Sidebar :open="sidebarOpen" @close="sidebarOpen = false" />

    <div class="lg:pl-64">
      <header class="sticky top-0 z-30 bg-white dark:bg-gray-800 border-b border-gray-200 dark:border-gray-700">
        <div class="flex items-center justify-between h-14 px-4">
          <button @click="sidebarOpen = true" class="lg:hidden text-gray-500 hover:text-gray-700 dark:text-gray-400 dark:hover:text-gray-200">
            <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 6h16M4 12h16M4 18h16" />
            </svg>
          </button>
          <div class="flex-1" />
          <div class="flex items-center gap-3">
            <ThemeToggle />
            <div class="flex items-center gap-2">
              <div class="w-8 h-8 rounded-full bg-primary-600 flex items-center justify-center text-white text-sm font-medium">
                {{ userInitials }}
              </div>
              <button @click="authStore.logout()" class="text-sm text-gray-500 hover:text-gray-700 dark:text-gray-400 dark:hover:text-gray-200">
                Logout
              </button>
            </div>
          </div>
        </div>
      </header>

      <main class="p-6">
        <router-view />
      </main>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { useAuthStore } from '@/stores/auth'
import Sidebar from '@/components/Sidebar.vue'
import ThemeToggle from '@/components/ThemeToggle.vue'

const authStore = useAuthStore()
const sidebarOpen = ref(false)

const userInitials = computed(() => {
  const name = authStore.user?.name || authStore.user?.email || '?'
  return name.charAt(0).toUpperCase()
})
</script>
