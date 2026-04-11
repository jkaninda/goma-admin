<template>
  <div class="flex-container">
    <div class="page-header">
      <h1>TLS Certificates</h1>
    </div>

    <div v-if="loading" class="card card-body loading-page">
      <div class="spinner"></div>
    </div>

    <template v-else>
      <!-- Overview cards -->
      <div class="overview-cards">
        <div class="overview-card">
          <div class="overview-value">{{ overview.totalCertificates }}</div>
          <div class="overview-label">Total</div>
        </div>
        <div class="overview-card">
          <div class="overview-value">{{ overview.acmeCertificates }}</div>
          <div class="overview-label">ACME</div>
        </div>
        <div class="overview-card">
          <div class="overview-value">{{ overview.customCertificates }}</div>
          <div class="overview-label">Custom</div>
        </div>
        <div class="overview-card" :class="{ 'overview-card-danger': overview.expiredCount > 0 }">
          <div class="overview-value">{{ overview.expiredCount }}</div>
          <div class="overview-label">Expired</div>
        </div>
        <div class="overview-card" :class="{ 'overview-card-warning': overview.expiringSoonCount > 0 }">
          <div class="overview-value">{{ overview.expiringSoonCount }}</div>
          <div class="overview-label">Expiring Soon</div>
        </div>
      </div>

      <EmptyState
        v-if="certificates.length === 0 && !search"
        title="No certificates"
        description="No TLS certificates found. Certificates are managed by Goma Gateway via ACME or custom certificate files."
      />

      <template v-else>
        <!-- Search & filter -->
        <div class="filter-row">
          <div class="search-bar">
            <svg class="search-icon" width="16" height="16" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <circle cx="11" cy="11" r="8" stroke-width="2" /><path stroke-linecap="round" stroke-width="2" d="m21 21-4.35-4.35" />
            </svg>
            <input v-model="search" class="form-input search-input" placeholder="Search certificates..." />
          </div>
          <div class="filter-buttons">
            <button
              v-for="f in sourceFilters"
              :key="f.value"
              :class="['btn btn-sm', sourceFilter === f.value ? 'btn-primary' : 'btn-secondary']"
              @click="sourceFilter = f.value"
            >
              {{ f.label }}
            </button>
          </div>
        </div>

        <div class="card">
          <div class="table-wrapper">
            <table>
              <thead>
                <tr>
                  <th>Domain</th>
                  <th>Source</th>
                  <th>Issuer</th>
                  <th>Expires</th>
                  <th>Status</th>
                </tr>
              </thead>
              <tbody>
                <tr
                  v-for="cert in filteredCertificates"
                  :key="cert.domain"
                  class="cursor-pointer"
                  @click="$router.push(`/certificates/${encodeURIComponent(cert.domain)}`)"
                >
                  <td>
                    <router-link :to="`/certificates/${encodeURIComponent(cert.domain)}`" class="cell-name-link">
                      {{ cert.domain }}
                    </router-link>
                    <div v-if="cert.domains.length > 1" class="cell-san-count">
                      +{{ cert.domains.length - 1 }} SAN{{ cert.domains.length - 1 > 1 ? 's' : '' }}
                    </div>
                  </td>
                  <td>
                    <span :class="['badge', cert.source === 'acme' ? 'badge-info' : 'badge-neutral']">
                      {{ cert.source === 'acme' ? 'ACME' : 'Custom' }}
                    </span>
                  </td>
                  <td class="cell-issuer">{{ cert.issuer || '-' }}</td>
                  <td class="cell-expires">{{ formatDate(cert.expiresAt) }}</td>
                  <td>
                    <span :class="['badge', statusBadge(cert)]">{{ statusLabel(cert) }}</span>
                  </td>
                </tr>
                <tr v-if="filteredCertificates.length === 0">
                  <td colspan="5" class="text-center text-muted" style="padding: 32px">No matching certificates</td>
                </tr>
              </tbody>
            </table>
          </div>

          <Pagination :pageable="pageable" @page="goToPage" />
        </div>
      </template>
    </template>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, watch, onMounted, onUnmounted } from 'vue'
import { certificatesApi, type CertificateSummary, type TLSOverview } from '@/api/certificates'
import EmptyState from '@/components/EmptyState.vue'
import Pagination from '@/components/Pagination.vue'

const loading = ref(true)
const certificates = ref<CertificateSummary[]>([])
const overview = ref<TLSOverview>({ totalCertificates: 0, acmeCertificates: 0, customCertificates: 0, expiredCount: 0, expiringSoonCount: 0 })
const page = ref(0)
const pageable = ref({ current_page: 0, total_pages: 1, total_elements: 0, size: 20, empty: true })
const search = ref('')
const sourceFilter = ref<'' | 'acme' | 'custom'>('')
let searchDebounceTimer: ReturnType<typeof setTimeout> | null = null

const sourceFilters = [
  { label: 'All', value: '' as const },
  { label: 'ACME', value: 'acme' as const },
  { label: 'Custom', value: 'custom' as const },
]

watch(search, () => {
  if (searchDebounceTimer) clearTimeout(searchDebounceTimer)
  searchDebounceTimer = setTimeout(() => {
    page.value = 0
    fetchCertificates()
  }, 300)
})

onUnmounted(() => {
  if (searchDebounceTimer) clearTimeout(searchDebounceTimer)
})

const filteredCertificates = computed(() => {
  if (!sourceFilter.value) return certificates.value
  return certificates.value.filter(c => c.source === sourceFilter.value)
})

function statusBadge(cert: CertificateSummary): string {
  if (cert.isExpired) return 'badge-danger'
  if (cert.daysUntilExpiry <= 30) return 'badge-warning'
  return 'badge-success'
}

function statusLabel(cert: CertificateSummary): string {
  if (cert.isExpired) return 'Expired'
  if (cert.daysUntilExpiry <= 30) return `${cert.daysUntilExpiry}d left`
  return 'Valid'
}

function formatDate(dateStr: string): string {
  if (!dateStr) return '-'
  const d = new Date(dateStr)
  return d.toLocaleDateString(undefined, { year: 'numeric', month: 'short', day: 'numeric' })
}

function goToPage(p: number) {
  page.value = p
  fetchCertificates()
}

async function fetchCertificates() {
  loading.value = true
  try {
    const [listRes, overviewRes] = await Promise.all([
      certificatesApi.list(page.value, 20, search.value),
      certificatesApi.overview(),
    ])
    certificates.value = listRes.data.data || []
    pageable.value = listRes.data.pageable
    overview.value = overviewRes.data
  } catch {
    // Error handled by API interceptor
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  fetchCertificates()
})
</script>

<style scoped>
.overview-cards {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(140px, 1fr));
  gap: 12px;
  margin-bottom: 20px;
}
.overview-card {
  background: var(--bg-primary);
  border: 1px solid var(--border-primary);
  border-radius: var(--radius);
  padding: 16px;
  text-align: center;
}
.overview-card-danger {
  border-color: var(--danger-300, #fca5a5);
  background: var(--danger-50, rgba(239, 68, 68, 0.05));
}
.overview-card-warning {
  border-color: var(--warning-300, #fcd34d);
  background: var(--warning-50, rgba(245, 158, 11, 0.05));
}
.overview-value {
  font-size: 28px;
  font-weight: 700;
  color: var(--text-primary);
  line-height: 1;
  margin-bottom: 4px;
}
.overview-label {
  font-size: 12px;
  font-weight: 500;
  color: var(--text-muted);
  text-transform: uppercase;
  letter-spacing: 0.05em;
}

.filter-row {
  display: flex;
  gap: 12px;
  margin-bottom: 16px;
  align-items: center;
}
.search-bar {
  position: relative;
  flex: 1;
}
.search-icon {
  position: absolute;
  left: 14px;
  top: 50%;
  transform: translateY(-50%);
  color: var(--text-muted);
  pointer-events: none;
}
.search-input {
  padding-left: 38px;
}
.filter-buttons {
  display: flex;
  gap: 4px;
}

.cell-name-link {
  font-weight: 600;
  color: var(--primary-600);
  text-decoration: none;
}
.cell-name-link:hover {
  text-decoration: underline;
}
.cell-san-count {
  font-size: 11px;
  color: var(--text-muted);
  margin-top: 2px;
}
.cell-issuer {
  color: var(--text-secondary);
  font-size: 13px;
}
.cell-expires {
  font-family: 'JetBrains Mono', 'Fira Code', monospace;
  font-size: 12px;
  color: var(--text-tertiary);
}
</style>
