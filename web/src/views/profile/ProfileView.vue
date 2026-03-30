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
          <div v-if="isOAuthUser" class="oauth-notice">
            Your account is managed by an external provider. Name and email cannot be changed here.
          </div>
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
                :disabled="isOAuthUser"
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
                :disabled="isOAuthUser"
              />
            </div>
            <button v-if="!isOAuthUser" type="submit" class="btn btn-primary" :disabled="savingProfile">
              {{ savingProfile ? 'Saving...' : 'Save' }}
            </button>
          </form>
        </div>
      </div>

      <!-- Two-Factor Authentication -->
      <div class="card" style="margin-top: 24px;">
        <div class="card-header">
          <h2>Two-Factor Authentication</h2>
          <span v-if="twoFactorEnabled" class="badge badge-success">Enabled</span>
          <span v-else class="badge badge-secondary">Disabled</span>
        </div>
        <div class="card-body">
          <!-- 2FA Not Enabled -->
          <template v-if="!twoFactorEnabled && !show2FASetup">
            <p class="tfa-description">Add an extra layer of security to your account by requiring a code from your authenticator app.</p>
            <button class="btn btn-primary" :disabled="tfaLoading" @click="startSetup2FA">
              {{ tfaLoading ? 'Setting up...' : 'Enable 2FA' }}
            </button>
          </template>

          <!-- 2FA Setup Flow -->
          <template v-if="show2FASetup">
            <div class="tfa-setup">
              <p class="tfa-description">Scan this QR code with your authenticator app (Google Authenticator, Authy, etc.):</p>
              <div class="tfa-qr">
                <img :src="`https://api.qrserver.com/v1/create-qr-code/?size=200x200&data=${encodeURIComponent(tfaURL)}`" alt="QR Code" width="200" height="200" />
              </div>
              <div class="tfa-secret-group">
                <label class="form-label">Or enter this secret manually:</label>
                <code class="tfa-secret">{{ tfaSecret }}</code>
              </div>
              <form @submit.prevent="verify2FA" style="margin-top: 16px">
                <div class="form-group">
                  <label class="form-label" for="tfa-code">Verification Code</label>
                  <input
                    id="tfa-code"
                    v-model="tfaCode"
                    type="text"
                    class="form-input totp-input"
                    placeholder="000000"
                    maxlength="6"
                    inputmode="numeric"
                    autocomplete="one-time-code"
                  />
                  <small class="form-hint">Enter the 6-digit code from your authenticator app to verify setup</small>
                </div>
                <div style="display: flex; gap: 8px">
                  <button type="submit" class="btn btn-primary" :disabled="tfaLoading">
                    {{ tfaLoading ? 'Verifying...' : 'Verify & Enable' }}
                  </button>
                  <button type="button" class="btn btn-secondary" @click="cancel2FASetup">Cancel</button>
                </div>
              </form>
            </div>
          </template>

          <!-- 2FA Enabled -->
          <template v-if="twoFactorEnabled && !show2FADisable">
            <p class="tfa-description">Two-factor authentication is currently enabled on your account.</p>
            <button class="btn btn-danger" @click="show2FADisable = true">Disable 2FA</button>
          </template>

          <!-- 2FA Disable Flow -->
          <template v-if="show2FADisable">
            <form @submit.prevent="disable2FA">
              <p class="tfa-description">Enter a code from your authenticator app to confirm disabling 2FA.</p>
              <div class="form-group">
                <label class="form-label" for="tfa-disable-code">Authentication Code</label>
                <input
                  id="tfa-disable-code"
                  v-model="tfaDisableCode"
                  type="text"
                  class="form-input totp-input"
                  placeholder="000000"
                  maxlength="6"
                  inputmode="numeric"
                  autocomplete="one-time-code"
                />
              </div>
              <div style="display: flex; gap: 8px">
                <button type="submit" class="btn btn-danger" :disabled="tfaDisableLoading">
                  {{ tfaDisableLoading ? 'Disabling...' : 'Confirm Disable' }}
                </button>
                <button type="button" class="btn btn-secondary" @click="show2FADisable = false; tfaDisableCode = ''">Cancel</button>
              </div>
            </form>
          </template>
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
import { profileApi } from '@/api/profile'
import api from '@/api/client'

const authStore = useAuthStore()
const notifications = useNotificationStore()

const savingProfile = ref(false)
const savingPassword = ref(false)
const twoFactorEnabled = ref(false)
const isOAuthUser = ref(false)

const profileForm = reactive({
  name: '',
  email: '',
})

const passwordForm = reactive({
  currentPassword: '',
  newPassword: '',
  confirmPassword: '',
})

onMounted(async () => {
  if (authStore.user) {
    profileForm.name = authStore.user.name
    profileForm.email = authStore.user.email
  }
  try {
    const res = await profileApi.getProfile()
    twoFactorEnabled.value = res.data.two_factor_enabled
    isOAuthUser.value = !!res.data.oauth_provider
  } catch { /* ignore */ }
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

// 2FA
const show2FASetup = ref(false)
const tfaSecret = ref('')
const tfaURL = ref('')
const tfaCode = ref('')
const tfaLoading = ref(false)
const tfaDisableCode = ref('')
const show2FADisable = ref(false)
const tfaDisableLoading = ref(false)

async function startSetup2FA() {
  tfaLoading.value = true
  try {
    const res = await profileApi.setup2FA()
    tfaSecret.value = res.data.secret
    tfaURL.value = res.data.url
    show2FASetup.value = true
  } catch (err: unknown) {
    const axiosErr = err as { response?: { data?: { message?: string } } }
    notifications.error(axiosErr.response?.data?.message || 'Failed to setup 2FA')
  } finally {
    tfaLoading.value = false
  }
}

async function verify2FA() {
  if (!tfaCode.value || tfaCode.value.length !== 6) {
    notifications.error('Please enter a valid 6-digit code')
    return
  }
  tfaLoading.value = true
  try {
    await profileApi.verify2FA(tfaCode.value)
    twoFactorEnabled.value = true
    show2FASetup.value = false
    tfaCode.value = ''
    tfaSecret.value = ''
    tfaURL.value = ''
    notifications.success('Two-factor authentication enabled')
  } catch (err: unknown) {
    const axiosErr = err as { response?: { data?: { message?: string } } }
    notifications.error(axiosErr.response?.data?.message || 'Invalid code. Please try again.')
  } finally {
    tfaLoading.value = false
  }
}

async function disable2FA() {
  if (!tfaDisableCode.value || tfaDisableCode.value.length !== 6) {
    notifications.error('Please enter a valid 6-digit code')
    return
  }
  tfaDisableLoading.value = true
  try {
    await profileApi.disable2FA(tfaDisableCode.value)
    twoFactorEnabled.value = false
    show2FADisable.value = false
    tfaDisableCode.value = ''
    notifications.success('Two-factor authentication disabled')
  } catch (err: unknown) {
    const axiosErr = err as { response?: { data?: { message?: string } } }
    notifications.error(axiosErr.response?.data?.message || 'Invalid code. Please try again.')
  } finally {
    tfaDisableLoading.value = false
  }
}

function cancel2FASetup() {
  show2FASetup.value = false
  tfaCode.value = ''
  tfaSecret.value = ''
  tfaURL.value = ''
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

<style scoped>
.tfa-description {
  font-size: 13px;
  color: var(--text-secondary);
  margin-bottom: 16px;
}

.tfa-setup {
  display: grid;
  gap: 16px;
}

.tfa-qr {
  display: flex;
  justify-content: center;
  padding: 16px;
  background: #fff;
  border-radius: var(--radius);
  border: 1px solid var(--border-primary);
  width: fit-content;
}

.tfa-secret-group {
  display: grid;
  gap: 6px;
}

.tfa-secret {
  display: block;
  padding: 10px 14px;
  background: var(--bg-secondary);
  border: 1px solid var(--border-primary);
  border-radius: var(--radius);
  font-size: 14px;
  font-family: 'JetBrains Mono', 'Fira Code', monospace;
  letter-spacing: 2px;
  word-break: break-all;
  user-select: all;
}

.totp-input {
  font-size: 20px;
  text-align: center;
  letter-spacing: 8px;
  font-family: 'JetBrains Mono', 'Fira Code', monospace;
  max-width: 220px;
}

.form-hint {
  font-size: 12px;
  color: var(--text-muted);
  margin-top: 4px;
}

.badge {
  display: inline-block;
  padding: 2px 8px;
  border-radius: 9999px;
  font-size: 11px;
  font-weight: 600;
  text-transform: uppercase;
  letter-spacing: 0.03em;
}

.badge-success {
  background: var(--success-50, #f0fdf4);
  color: var(--success-700, #15803d);
}

.badge-secondary {
  background: var(--bg-secondary);
  color: var(--text-muted);
}

.oauth-notice {
  padding: 10px 14px;
  margin-bottom: 16px;
  border-radius: var(--radius);
  background: var(--warning-50, #fffbeb);
  color: var(--warning-700, #a16207);
  font-size: 13px;
  border: 1px solid var(--warning-200, #fde68a);
}
</style>
