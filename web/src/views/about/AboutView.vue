<template>
  <div>
    <div class="page-header">
      <h1>About</h1>
    </div>

    <!-- Product hero -->
    <div class="card section-card">
      <div class="card-body about-hero">
        <img src="/logo.png" alt="Goma" class="about-logo" />
        <div class="about-hero-text">
          <h2>{{ info.name || 'Goma Admin' }}</h2>
          <p class="about-tagline">
            Control Plane for Goma Gateway — Manage, configure, and monitor distributed API gateways from a single, unified dashboard.
          </p>
          <div class="about-badges">
            <span class="badge badge-info">{{ info.version ? `v${info.version}` : 'dev' }}</span>
            <span v-if="info.commit_id && info.commit_id !== 'unknown'" class="badge badge-muted text-mono">
              {{ info.commit_id.slice(0, 7) }}
            </span>
          </div>
        </div>
      </div>
    </div>

    <!-- Build information -->
    <div class="card section-card">
      <div class="card-header">
        <h2>Build Information</h2>
      </div>
      <div class="card-body">
        <dl class="about-details">
          <div class="about-detail">
            <dt>Application</dt>
            <dd>{{ info.name || 'Goma Admin' }}</dd>
          </div>
          <div class="about-detail">
            <dt>Version</dt>
            <dd>{{ info.version || 'dev' }}</dd>
          </div>
          <div class="about-detail">
            <dt>Commit</dt>
            <dd class="text-mono">{{ info.commit_id || 'unknown' }}</dd>
          </div>
          <div class="about-detail">
            <dt>API Documentation</dt>
            <dd>
              <span class="badge" :class="info.openapi_docs ? 'badge-success' : 'badge-muted'">
                {{ info.openapi_docs ? 'Enabled' : 'Disabled' }}
              </span>
            </dd>
          </div>
        </dl>
      </div>
    </div>

    <!-- Resources -->
    <div class="card section-card">
      <div class="card-header">
        <h2>Resources</h2>
      </div>
      <div class="card-body">
        <div class="about-links">
          <a
            v-for="link in resources"
            :key="link.href"
            class="about-link"
            :href="link.href"
            target="_blank"
            rel="noopener noreferrer"
          >
            <span class="about-link-icon" v-html="link.icon" />
            <span class="about-link-text">
              <span class="about-link-title">{{ link.title }}</span>
              <span class="about-link-desc">{{ link.desc }}</span>
            </span>
            <svg class="about-link-arrow" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
              <path d="M7 17L17 7M17 7H7M17 7v10" />
            </svg>
          </a>
        </div>
      </div>
    </div>

    <p class="about-copyright">
      &copy; {{ new Date().getFullYear() }}
      <a href="https://jkaninda.dev" target="_blank" rel="noopener noreferrer">Jonas Kaninda</a>.
      Released under the MIT License.
    </p>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import axios from 'axios'

interface AppInfo {
  name: string
  version: string
  commit_id: string
  openapi_docs: boolean
}

const info = ref<AppInfo>({ name: '', version: '', commit_id: '', openapi_docs: false })

const resources = computed(() => {
  const items = [
    {
      title: 'Documentation',
      desc: 'Guides and reference for Goma',
      href: 'https://goma.jkaninda.dev',
      icon: '<svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M2 3h6a4 4 0 0 1 4 4v14a3 3 0 0 0-3-3H2z"/><path d="M22 3h-6a4 4 0 0 0-4 4v14a3 3 0 0 1 3-3h7z"/></svg>',
    },
    {
      title: 'GitHub',
      desc: 'Source code and releases',
      href: 'https://github.com/jkaninda/goma-admin',
      icon: '<svg width="20" height="20" viewBox="0 0 24 24" fill="currentColor"><path d="M12 0C5.37 0 0 5.37 0 12c0 5.31 3.435 9.795 8.205 11.385.6.105.825-.255.825-.57 0-.285-.015-1.23-.015-2.235-3.015.555-3.795-.735-4.035-1.41-.135-.345-.72-1.41-1.23-1.695-.42-.225-1.02-.78-.015-.795.945-.015 1.62.87 1.845 1.23 1.08 1.815 2.805 1.305 3.495.99.105-.78.42-1.305.765-1.605-2.67-.3-5.46-1.335-5.46-5.925 0-1.305.465-2.385 1.23-3.225-.12-.3-.54-1.53.12-3.18 0 0 1.005-.315 3.3 1.23.96-.27 1.98-.405 3-.405s2.04.135 3 .405c2.295-1.56 3.3-1.23 3.3-1.23.66 1.65.24 2.88.12 3.18.765.84 1.23 1.905 1.23 3.225 0 4.605-2.805 5.625-5.475 5.925.435.375.81 1.095.81 2.22 0 1.605-.015 2.895-.015 3.3 0 .315.225.69.825.57A12.02 12.02 0 0 0 24 12c0-6.63-5.37-12-12-12z"/></svg>',
    },
    {
      title: 'Goma Gateway',
      desc: 'The gateway managed by this UI',
      href: 'https://github.com/jkaninda/goma-gateway',
      icon: '<svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><circle cx="6" cy="6" r="3"/><circle cx="18" cy="18" r="3"/><path d="M6 9a9 9 0 0 0 9 9"/></svg>',
    },
  ]
  if (info.value.openapi_docs) {
    items.push({
      title: 'API Reference',
      desc: 'Interactive Swagger documentation',
      href: '/docs',
      icon: '<svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z"/><polyline points="14 2 14 8 20 8"/><line x1="16" y1="13" x2="8" y2="13"/><line x1="16" y1="17" x2="8" y2="17"/></svg>',
    })
  }
  return items
})

onMounted(() => {
  axios
    .get<AppInfo>('/info')
    .then((res) => {
      info.value = { ...info.value, ...res.data }
    })
    .catch(() => {})
})
</script>

<style scoped>
.section-card {
  margin-bottom: 24px;
}

/* ── Hero ── */
.about-hero {
  display: flex;
  align-items: center;
  gap: 20px;
}

.about-logo {
  width: 64px;
  height: 64px;
  border-radius: var(--radius-lg);
  flex-shrink: 0;
}

.about-hero-text h2 {
  margin: 0 0 4px;
}

.about-tagline {
  color: var(--text-secondary);
  font-size: 14px;
  line-height: 1.5;
  margin: 0 0 12px;
  max-width: 60ch;
}

.about-badges {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
}

.badge-muted {
  background: var(--bg-tertiary);
  color: var(--text-muted);
}

/* ── Details ── */
.about-details {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 16px 32px;
  margin: 0;
}

.about-detail dt {
  font-size: 12px;
  font-weight: 600;
  text-transform: uppercase;
  letter-spacing: 0.04em;
  color: var(--text-muted);
  margin-bottom: 4px;
}

.about-detail dd {
  margin: 0;
  font-size: 14px;
  color: var(--text-primary);
  word-break: break-all;
}

.text-mono {
  font-family: ui-monospace, SFMono-Regular, Menlo, monospace;
}

/* ── Resource links ── */
.about-links {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 12px;
}

.about-link {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 14px 16px;
  border: 1px solid var(--border-primary);
  border-radius: var(--radius);
  text-decoration: none;
  color: var(--text-primary);
  transition: border-color var(--transition), background var(--transition);
}

.about-link:hover {
  border-color: var(--primary-500);
  background: var(--bg-hover);
}

.about-link-icon {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 40px;
  height: 40px;
  border-radius: var(--radius);
  background: var(--bg-tertiary);
  color: var(--primary-600);
  flex-shrink: 0;
}

.about-link-text {
  display: flex;
  flex-direction: column;
  gap: 2px;
  min-width: 0;
}

.about-link-title {
  font-size: 14px;
  font-weight: 600;
}

.about-link-desc {
  font-size: 12px;
  color: var(--text-muted);
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.about-link-arrow {
  margin-left: auto;
  color: var(--text-muted);
  flex-shrink: 0;
}

/* ── Copyright ── */
.about-copyright {
  text-align: center;
  font-size: 13px;
  color: var(--text-muted);
  margin: 8px 0 0;
}

.about-copyright a {
  color: var(--primary-600);
  text-decoration: none;
}

.about-copyright a:hover {
  color: var(--primary-700);
}

@media (max-width: 640px) {
  .about-hero {
    flex-direction: column;
    text-align: center;
  }
  .about-badges {
    justify-content: center;
  }
  .about-details,
  .about-links {
    grid-template-columns: 1fr;
  }
}
</style>
