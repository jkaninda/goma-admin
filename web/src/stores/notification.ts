import { defineStore } from 'pinia'
import { ref } from 'vue'

export interface Notification {
  id: number
  message: string
  type: 'success' | 'error' | 'info'
}

let nextId = 0

export const useNotificationStore = defineStore('notification', () => {
  const notifications = ref<Notification[]>([])

  function show(message: string, type: Notification['type'] = 'info', duration?: number) {
    const id = ++nextId
    notifications.value.push({ id, message, type })

    const timeout = duration ?? (type === 'error' ? 6000 : 4000)
    setTimeout(() => dismiss(id), timeout)

    return id
  }

  function success(message: string) {
    return show(message, 'success')
  }

  function error(message: string) {
    return show(message, 'error')
  }

  function info(message: string) {
    return show(message, 'info')
  }

  function dismiss(id: number) {
    const index = notifications.value.findIndex((n) => n.id === id)
    if (index !== -1) {
      notifications.value.splice(index, 1)
    }
  }

  return { notifications, show, success, error, info, dismiss }
})
