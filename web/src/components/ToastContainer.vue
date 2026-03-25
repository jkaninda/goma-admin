<template>
  <div v-if="notifications.length" class="toast-container">
    <div
      v-for="notification in notifications"
      :key="notification.id"
      class="toast"
      :class="`toast-${notification.type}`"
      @click="dismiss(notification.id)"
    >
      {{ notification.message }}
    </div>
  </div>
</template>

<script setup lang="ts">
import { storeToRefs } from 'pinia'
import { useNotificationStore } from '@/stores/notification'

const store = useNotificationStore()
const { notifications } = storeToRefs(store)
const { dismiss } = store
</script>
