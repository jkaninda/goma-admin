<template>
  <div>
    <div class="page-header">
      <h1>Profile</h1>
    </div>

    <div class="section">
      <!-- My Profile -->
      <div class="card">
        <div class="card-header">
          <h2>My Profile</h2>
        </div>
        <div class="card-body">
          <form @submit.prevent="handleProfileUpdate">
            <div class="form-group">
              <label class="form-label" for="profile-name">Name</label>
              <input
                id="profile-name"
                v-model="profileForm.name"
                type="text"
                class="form-input"
                placeholder="Your name"
                required
              />
            </div>
            <div class="form-group">
              <label class="form-label" for="profile-email">Email</label>
              <input
                id="profile-email"
                v-model="profileForm.email"
                type="email"
                class="form-input"
                placeholder="you@example.com"
                required
              />
            </div>
            <button type="submit" class="btn btn-primary" :disabled="savingProfile">
              {{ savingProfile ? 'Saving...' : 'Save' }}
            </button>
          </form>
        </div>
      </div>

      <!-- Change Password -->
      <div class="card" style="margin-top: 24px;">
        <div class="card-header">
          <h2>Change Password</h2>
        </div>
        <div class="card-body">
          <form @submit.prevent="handlePasswordUpdate">
            <div class="form-group">
              <label class="form-label" for="current-password">Current Password</label>
              <input
                id="current-password"
                v-model="passwordForm.currentPassword"
                type="password"
                class="form-input"
                placeholder="Enter current password"
                required
              />
            </div>
            <div class="form-group">
              <label class="form-label" for="new-password">New Password</label>
              <input
                id="new-password"
                v-model="passwordForm.newPassword"
                type="password"
                class="form-input"
                placeholder="Enter new password"
                required
              />
            </div>
            <div class="form-group">
              <label class="form-label" for="confirm-password">Confirm Password</label>
              <input
                id="confirm-password"
                v-model="passwordForm.confirmPassword"
                type="password"
                class="form-input"
                placeholder="Confirm new password"
                required
              />
            </div>
            <button type="submit" class="btn btn-primary" :disabled="savingPassword">
              {{ savingPassword ? 'Updating...' : 'Update' }}
            </button>
          </form>
        </div>
      </div>

      <!-- Danger Zone -->
      <div class="danger-zone">
        <h3>Danger Zone</h3>
        <p>Once you delete your account, there is no going back. Please be certain.</p>
        <button class="btn btn-danger" disabled>Delete Account</button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { reactive, ref, onMounted } from 'vue'
import { useAuthStore } from '@/stores/auth'
import { useNotificationStore } from '@/stores/notification'
import api from '@/api/client'

const authStore = useAuthStore()
const notifications = useNotificationStore()

const savingProfile = ref(false)
const savingPassword = ref(false)

const profileForm = reactive({
  name: '',
  email: '',
})

const passwordForm = reactive({
  currentPassword: '',
  newPassword: '',
  confirmPassword: '',
})

onMounted(() => {
  if (authStore.user) {
    profileForm.name = authStore.user.name
    profileForm.email = authStore.user.email
  }
})

async function handleProfileUpdate() {
  savingProfile.value = true
  try {
    await api.put('/profile', {
      name: profileForm.name,
      email: profileForm.email,
    })
    if (authStore.user) {
      authStore.user.name = profileForm.name
      authStore.user.email = profileForm.email
      localStorage.setItem('user', JSON.stringify(authStore.user))
    }
    notifications.success('Profile updated successfully.')
  } catch (err: unknown) {
    const axiosErr = err as { response?: { data?: { message?: string } } }
    notifications.error(axiosErr.response?.data?.message || 'Failed to update profile.')
  } finally {
    savingProfile.value = false
  }
}

async function handlePasswordUpdate() {
  if (passwordForm.newPassword !== passwordForm.confirmPassword) {
    notifications.error('New passwords do not match.')
    return
  }
  savingPassword.value = true
  try {
    await api.put('/profile/password', {
      current_password: passwordForm.currentPassword,
      new_password: passwordForm.newPassword,
    })
    passwordForm.currentPassword = ''
    passwordForm.newPassword = ''
    passwordForm.confirmPassword = ''
    notifications.success('Password updated successfully.')
  } catch (err: unknown) {
    const axiosErr = err as { response?: { data?: { message?: string } } }
    notifications.error(axiosErr.response?.data?.message || 'Failed to update password.')
  } finally {
    savingPassword.value = false
  }
}
</script>
