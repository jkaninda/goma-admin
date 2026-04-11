<template>
  <div>
    <div v-if="loading" class="loading-page"><div class="spinner"></div></div>

    <div v-else-if="error" class="detail-error">
      <div class="card card-body" style="text-align: center; padding: 48px 24px">
        <svg width="48" height="48" fill="none" stroke="var(--text-muted)" viewBox="0 0 24 24" style="margin: 0 auto 16px">
          <circle cx="12" cy="12" r="10" stroke-width="1.5" />
          <path stroke-linecap="round" stroke-width="1.5" d="M12 8v4m0 4h.01" />
        </svg>
        <h2 style="margin-bottom: 8px">Certificate not found</h2>
        <p class="text-muted" style="margin-bottom: 20px">The certificate you're looking for doesn't exist.</p>
        <router-link to="/certificates" class="btn btn-primary">Back to Certificates</router-link>
      </div>
    </div>

    <div v-else-if="cert">
      <!-- Header -->
      <div class="detail-header">
        <router-link to="/certificates" class="back-link">
          <svg width="20" height="20" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7" />
          </svg>
        </router-link>
        <div class="detail-header-info">
          <h1>{{ cert.domain }}</h1>
          <span :class="['badge', cert.source === 'acme' ? 'badge-info' : 'badge-neutral']">
            {{ cert.source === 'acme' ? 'ACME' : 'Custom' }}
          </span>
          <span :class="['badge', statusBadge]">{{ statusLabel }}</span>
        </div>
      </div>

      <div class="detail-cards">
        <!-- Certificate Info -->
        <div class="card">
          <div class="card-header"><h2>Certificate Info</h2></div>
          <div class="card-body">
            <dl class="detail-grid">
              <div class="detail-row">
                <dt>Domain</dt>
                <dd class="text-mono">{{ cert.domain }}</dd>
              </div>
              <div class="detail-row">
                <dt>Subject Alternative Names</dt>
                <dd>
                  <div class="san-tags">
                    <span v-for="d in cert.domains" :key="d" class="san-tag">{{ d }}</span>
                  </div>
                </dd>
              </div>
              <div class="detail-row">
                <dt>Issuer</dt>
                <dd>{{ cert.issuer || '-' }}</dd>
              </div>
              <div class="detail-row">
                <dt>Serial Number</dt>
                <dd class="text-mono">{{ cert.serialNumber || '-' }}</dd>
              </div>
              <div class="detail-row">
                <dt>Fingerprint (SHA-256)</dt>
                <dd class="text-mono fingerprint">{{ cert.fingerprint || '-' }}</dd>
              </div>
            </dl>
          </div>
        </div>

        <!-- Validity -->
        <div class="card">
          <div class="card-header"><h2>Validity</h2></div>
          <div class="card-body">
            <dl class="detail-grid">
              <div class="detail-row">
                <dt>Status</dt>
                <dd><span :class="['badge', statusBadge]">{{ statusLabel }}</span></dd>
              </div>
              <div class="detail-row">
                <dt>Issued</dt>
                <dd>{{ formatDate(cert.issuedAt) }}</dd>
              </div>
              <div class="detail-row">
                <dt>Expires</dt>
                <dd>{{ formatDate(cert.expiresAt) }}</dd>
              </div>
              <div v-if="!cert.isExpired" class="detail-row">
                <dt>Days Remaining</dt>
                <dd>
                  <span :class="{ 'text-warning': cert.daysUntilExpiry <= 30, 'text-danger': cert.daysUntilExpiry <= 7 }">
                    {{ cert.daysUntilExpiry }} day{{ cert.daysUntilExpiry !== 1 ? 's' : '' }}
                  </span>
                </dd>
              </div>
            </dl>
          </div>
        </div>

        <!-- Source -->
        <div class="card">
          <div class="card-header"><h2>Source</h2></div>
          <div class="card-body">
            <dl class="detail-grid">
              <div class="detail-row">
                <dt>Source</dt>
                <dd>
                  <span :class="['badge', cert.source === 'acme' ? 'badge-info' : 'badge-neutral']">
                    {{ cert.source === 'acme' ? 'ACME (Let\'s Encrypt)' : 'Custom Certificate' }}
                  </span>
                </dd>
              </div>
              <div class="detail-row">
                <dt>Managed By</dt>
                <dd>{{ cert.source === 'acme' ? 'Goma Gateway (automatic renewal)' : 'Manual' }}</dd>
              </div>
            </dl>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { certificatesApi, type CertificateSummary } from '@/api/certificates'

const props = defineProps<{ domain: string }>()

const loading = ref(true)
const error = ref(false)
const cert = ref<CertificateSummary | null>(null)

const statusBadge = computed(() => {
  if (!cert.value) return 'badge-neutral'
  if (cert.value.isExpired) return 'badge-danger'
  if (cert.value.daysUntilExpiry <= 30) return 'badge-warning'
  return 'badge-success'
})

const statusLabel = computed(() => {
  if (!cert.value) return ''
  if (cert.value.isExpired) return 'Expired'
  if (cert.value.daysUntilExpiry <= 30) return `Expiring in ${cert.value.daysUntilExpiry}d`
  return 'Valid'
})

function formatDate(dateStr: string): string {
  if (!dateStr) return '-'
  const d = new Date(dateStr)
  return d.toLocaleDateString(undefined, {
    year: 'numeric',
    month: 'short',
    day: 'numeric',
    hour: '2-digit',
    minute: '2-digit',
  })
}

onMounted(async () => {
  try {
    const res = await certificatesApi.get(props.domain)
    cert.value = res.data
  } catch {
    error.value = true
  } finally {
    loading.value = false
  }
})
</script>

<style scoped>
.detail-header {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-bottom: 16px;
}
.detail-header-info {
  display: flex;
  align-items: center;
  gap: 12px;
  flex: 1;
  min-width: 0;
}
.detail-header h1 {
  font-size: 22px;
  font-weight: 700;
  color: var(--text-primary);
  letter-spacing: -0.02em;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}
.back-link {
  display: flex;
  align-items: center;
  color: var(--text-muted);
  transition: color var(--transition);
}
.back-link:hover { color: var(--text-secondary); }

.detail-cards {
  display: flex;
  flex-direction: column;
  gap: 20px;
  max-width: 800px;
}

.detail-grid {
  display: grid;
  grid-template-columns: 1fr;
  gap: 0;
  margin: 0;
}
.detail-row {
  display: flex;
  align-items: baseline;
  padding: 10px 0;
  border-bottom: 1px solid var(--border-secondary);
}
.detail-row:last-child {
  border-bottom: none;
}
.detail-row dt {
  width: 200px;
  flex-shrink: 0;
  font-size: 13px;
  font-weight: 600;
  color: var(--text-tertiary);
  text-transform: uppercase;
  letter-spacing: 0.03em;
}
.detail-row dd {
  margin: 0;
  font-size: 14px;
  color: var(--text-primary);
}

.text-mono {
  font-family: 'SF Mono', 'Fira Code', 'Fira Mono', Menlo, Consolas, monospace;
  font-size: 13px;
}
.text-muted {
  color: var(--text-muted);
  font-style: italic;
}
.text-warning {
  color: var(--warning-600, #d97706);
  font-weight: 600;
}
.text-danger {
  color: var(--danger-600, #dc2626);
  font-weight: 600;
}

.fingerprint {
  font-size: 11px;
  word-break: break-all;
  line-height: 1.6;
}

.san-tags {
  display: flex;
  flex-wrap: wrap;
  gap: 4px;
}
.san-tag {
  display: inline-block;
  padding: 2px 10px;
  background: var(--bg-tertiary);
  border: 1px solid var(--border-primary);
  border-radius: 4px;
  font-family: 'SF Mono', 'Fira Code', 'Fira Mono', Menlo, Consolas, monospace;
  font-size: 12px;
  color: var(--text-secondary);
}
</style>
