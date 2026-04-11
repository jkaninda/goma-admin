<template>
  <div>
    <div v-if="loading" class="loading-page"><div class="spinner"></div></div>

    <div v-else-if="error" class="detail-error">
      <div class="card card-body" style="text-align: center; padding: 48px 24px">
        <svg width="48" height="48" fill="none" stroke="var(--text-muted)" viewBox="0 0 24 24" style="margin: 0 auto 16px">
          <circle cx="12" cy="12" r="10" stroke-width="1.5" />
          <path stroke-linecap="round" stroke-width="1.5" d="M12 8v4m0 4h.01" />
        </svg>
        <h2 style="margin-bottom: 8px">Route not found</h2>
        <p class="text-muted" style="margin-bottom: 20px">The route you're looking for doesn't exist or has been deleted.</p>
        <router-link to="/routes" class="btn btn-primary">Back to Routes</router-link>
      </div>
    </div>

    <div v-else-if="route">
      <!-- Header -->
      <div class="detail-header">
        <router-link to="/routes" class="back-link">
          <svg width="20" height="20" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7" />
          </svg>
        </router-link>
        <div class="detail-header-info">
          <h1>{{ route.name }}</h1>
          <span :class="['status-indicator', config.enabled !== false ? 'status-online' : 'status-offline']">
            <span class="status-dot"></span>
            {{ config.enabled !== false ? 'Online' : 'Offline' }}
          </span>
        </div>
        <div class="detail-header-actions">
          <button class="btn btn-secondary btn-sm" @click="handleEdit">Edit</button>
          <button class="btn btn-danger btn-sm" @click="handleDelete">Delete</button>
        </div>
      </div>

      <!-- Quick stats -->
      <div class="quick-stats">
        <div class="quick-stat" v-if="methodsList.length">
          <span v-for="m in methodsList" :key="m" :class="['badge', 'badge-sm', methodBadge(m)]">{{ m }}</span>
        </div>
        <div class="quick-stat" v-else>
          <span class="badge badge-sm badge-neutral">All Methods</span>
        </div>
        <div class="quick-stat-divider"></div>
        <div class="quick-stat">
          <svg width="14" height="14" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 12a9 9 0 01-9 9m9-9a9 9 0 00-9-9m9 9H3m9 9a9 9 0 01-9-9m9 9c1.66 0 3-4.03 3-9s-1.34-9-3-9m0 18c-1.66 0-3-4.03-3-9s1.34-9 3-9" />
          </svg>
          <span>{{ hostsList.length ? hostsList.join(', ') : 'Any host' }}</span>
        </div>
        <template v-if="middlewaresList.length">
          <div class="quick-stat-divider"></div>
          <div class="quick-stat">
            <svg width="14" height="14" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 6V4m0 2a2 2 0 100 4m0-4a2 2 0 110 4m-6 8a2 2 0 100-4m0 4a2 2 0 110-4m0 4v2m0-6V4m6 6v10m6-2a2 2 0 100-4m0 4a2 2 0 110-4m0 4v2m0-6V4"/>
            </svg>
            <span>{{ middlewaresList.length }} middleware{{ middlewaresList.length !== 1 ? 's' : '' }}</span>
          </div>
        </template>
        <template v-if="backendsList.length">
          <div class="quick-stat-divider"></div>
          <div class="quick-stat">
            <svg width="14" height="14" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <rect x="2" y="2" width="20" height="8" rx="2" stroke-width="1.5"/><rect x="2" y="14" width="20" height="8" rx="2" stroke-width="1.5"/><circle cx="6" cy="6" r="1" fill="currentColor"/><circle cx="6" cy="18" r="1" fill="currentColor"/>
            </svg>
            <span>{{ backendsList.length }} backend{{ backendsList.length !== 1 ? 's' : '' }}</span>
            <span v-if="isWeighted" class="badge badge-sm badge-warning">Canary</span>
          </div>
        </template>
      </div>

      <!-- Animated Traffic Flow -->
      <div class="traffic-flow card">
        <div class="traffic-flow-inner">
          <!-- Client -->
          <div class="tf-node tf-node--client">
            <div class="tf-node-icon">
              <svg width="22" height="22" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <rect x="2" y="3" width="20" height="14" rx="2" stroke-width="1.5"/>
                <path stroke-linecap="round" stroke-width="1.5" d="M8 21h8m-4-4v4"/>
              </svg>
            </div>
            <div class="tf-node-body">
              <div class="tf-node-label">Client</div>
              <div class="tf-node-detail">
                <template v-if="hostsList.length">
                  <a v-for="h in hostsList" :key="h" :href="hostUrl(h)" target="_blank" rel="noopener" class="badge badge-info flow-badge flow-host-link">{{ h }}</a>
                </template>
                <span v-else class="text-muted">Any host</span>
              </div>
            </div>
          </div>

          <!-- Connector: Client -> Path -->
          <div class="tf-connector">
            <div class="tf-connector-line">
              <div class="tf-packet tf-packet--1"></div>
              <div class="tf-packet tf-packet--2"></div>
              <div class="tf-packet tf-packet--3"></div>
            </div>
          </div>

          <!-- Path -->
          <div class="tf-node tf-node--path">
            <div class="tf-node-icon">
              <svg width="22" height="22" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M13 10V3L4 14h7v7l9-11h-7z"/>
              </svg>
            </div>
            <div class="tf-node-body">
              <div class="tf-node-label">Path</div>
              <div class="tf-node-detail tf-node-detail--copyable" @click="copyToClipboard(config.path || '/')">
                <span class="text-mono">{{ config.path || '/' }}</span>
                <svg class="copy-icon" width="12" height="12" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <rect x="9" y="9" width="13" height="13" rx="2" stroke-width="2"/><path stroke-width="2" d="M5 15H4a2 2 0 01-2-2V4a2 2 0 012-2h9a2 2 0 012 2v1"/>
                </svg>
              </div>
              <div v-if="config.rewrite" class="tf-rewrite">
                <svg width="12" height="12" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h5M20 20v-5h-5"/>
                  <path stroke-linecap="round" stroke-width="2" d="M20.49 9A9 9 0 005.64 5.64L4 4m16 16l-1.64-1.64A9 9 0 013.51 15"/>
                </svg>
                <span>{{ config.rewrite }}</span>
              </div>
            </div>
          </div>

          <!-- Middlewares (if any) -->
          <template v-if="middlewaresList.length">
            <div class="tf-connector">
              <div class="tf-connector-line">
                <div class="tf-packet tf-packet--1"></div>
                <div class="tf-packet tf-packet--2"></div>
                <div class="tf-packet tf-packet--3"></div>
              </div>
            </div>

            <div class="tf-middlewares-group">
              <div class="tf-middlewares-label">
                <svg width="14" height="14" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 6V4m0 2a2 2 0 100 4m0-4a2 2 0 110 4m-6 8a2 2 0 100-4m0 4a2 2 0 110-4m0 4v2m0-6V4m6 6v10m6-2a2 2 0 100-4m0 4a2 2 0 110-4m0 4v2m0-6V4"/>
                </svg>
                <span>Middlewares</span>
              </div>
              <div class="tf-middlewares-pipeline">
                <div v-for="(mw, i) in middlewaresList" :key="mw" class="tf-mw-step">
                  <div class="tf-node tf-node--middleware">
                    <span class="tf-mw-order">{{ i + 1 }}</span>
                    <span class="tf-mw-name text-mono">{{ mw }}</span>
                  </div>
                  <div v-if="i < middlewaresList.length - 1" class="tf-mw-connector">
                    <div class="tf-mw-connector-line">
                      <div class="tf-packet-sm tf-packet--1"></div>
                    </div>
                  </div>
                </div>
              </div>
            </div>
          </template>

          <!-- Connector: (Middlewares|Path) -> Target -->
          <div class="tf-connector">
            <div class="tf-connector-line">
              <div class="tf-packet tf-packet--1"></div>
              <div class="tf-packet tf-packet--2"></div>
              <div class="tf-packet tf-packet--3"></div>
            </div>
          </div>

          <!-- Target / Backends -->
          <div v-if="!backendsList.length" class="tf-node tf-node--target">
            <div class="tf-node-icon">
              <svg width="22" height="22" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <rect x="2" y="2" width="20" height="8" rx="2" stroke-width="1.5"/><rect x="2" y="14" width="20" height="8" rx="2" stroke-width="1.5"/><circle cx="6" cy="6" r="1" fill="currentColor"/><circle cx="6" cy="18" r="1" fill="currentColor"/>
              </svg>
            </div>
            <div class="tf-node-body">
              <div class="tf-node-label">Target</div>
              <div class="tf-node-detail tf-node-detail--copyable" @click="copyToClipboard(config.target || '')">
                <span class="text-mono">{{ config.target || '-' }}</span>
                <svg v-if="config.target" class="copy-icon" width="12" height="12" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <rect x="9" y="9" width="13" height="13" rx="2" stroke-width="2"/><path stroke-width="2" d="M5 15H4a2 2 0 01-2-2V4a2 2 0 012-2h9a2 2 0 012 2v1"/>
                </svg>
              </div>
            </div>
          </div>
          <div v-else class="tf-backends-group">
            <div class="tf-backends-label">
              <svg width="14" height="14" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <rect x="2" y="2" width="20" height="8" rx="2" stroke-width="1.5"/><rect x="2" y="14" width="20" height="8" rx="2" stroke-width="1.5"/><circle cx="6" cy="6" r="1" fill="currentColor"/><circle cx="6" cy="18" r="1" fill="currentColor"/>
              </svg>
              <span>Backends</span>
              <span v-if="isWeighted" class="badge badge-warning" style="font-size: 10px; padding: 1px 6px;">Canary</span>
              <span v-else class="badge badge-info" style="font-size: 10px; padding: 1px 6px;">Round Robin</span>
            </div>
            <div class="tf-backends-list">
              <div v-for="(b, i) in backendsList" :key="i" class="tf-node tf-node--backend">
                <div class="tf-backend-header">
                  <svg width="16" height="16" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <rect x="2" y="2" width="20" height="8" rx="2" stroke-width="1.5"/><rect x="2" y="14" width="20" height="8" rx="2" stroke-width="1.5"/><circle cx="6" cy="6" r="1" fill="currentColor"/><circle cx="6" cy="18" r="1" fill="currentColor"/>
                  </svg>
                  <span class="text-mono">{{ b.endpoint }}</span>
                </div>
                <div v-if="isWeighted" class="tf-backend-weight">
                  <div class="tf-backend-weight-track">
                    <div class="tf-backend-weight-fill" :style="{ width: b.weight + '%' }"></div>
                  </div>
                  <span class="tf-backend-weight-label">{{ b.weight }}%</span>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>

      <div class="detail-cards">
        <!-- Core Configuration -->
        <div class="card">
          <div class="card-header">
            <div class="card-header-title">
              <svg width="16" height="16" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10.325 4.317c.426-1.756 2.924-1.756 3.35 0a1.724 1.724 0 002.573 1.066c1.543-.94 3.31.826 2.37 2.37a1.724 1.724 0 001.066 2.573c1.756.426 1.756 2.924 0 3.35a1.724 1.724 0 00-1.066 2.573c.94 1.543-.826 3.31-2.37 2.37a1.724 1.724 0 00-2.573 1.066c-.426 1.756-2.924 1.756-3.35 0a1.724 1.724 0 00-2.573-1.066c-1.543.94-3.31-.826-2.37-2.37a1.724 1.724 0 00-1.066-2.573c-1.756-.426-1.756-2.924 0-3.35a1.724 1.724 0 001.066-2.573c-.94-1.543.826-3.31 2.37-2.37.996.608 2.296.07 2.572-1.065z" />
                <circle cx="12" cy="12" r="3" stroke-width="2" />
              </svg>
              <h2>Configuration</h2>
            </div>
          </div>
          <div class="card-body">
            <dl class="detail-grid">
              <div class="detail-row">
                <dt>Path</dt>
                <dd>
                  <div class="value-with-copy">
                    <span class="text-mono">{{ config.path || '-' }}</span>
                    <button v-if="config.path" class="copy-btn" @click="copyToClipboard(config.path)" title="Copy path">
                      <svg width="14" height="14" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <rect x="9" y="9" width="13" height="13" rx="2" stroke-width="2"/><path stroke-width="2" d="M5 15H4a2 2 0 01-2-2V4a2 2 0 012-2h9a2 2 0 012 2v1"/>
                      </svg>
                    </button>
                  </div>
                </dd>
              </div>
              <div v-if="!backendsList.length" class="detail-row">
                <dt>Target</dt>
                <dd>
                  <div class="value-with-copy">
                    <span class="text-mono">{{ config.target || '-' }}</span>
                    <button v-if="config.target" class="copy-btn" @click="copyToClipboard(config.target)" title="Copy target">
                      <svg width="14" height="14" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <rect x="9" y="9" width="13" height="13" rx="2" stroke-width="2"/><path stroke-width="2" d="M5 15H4a2 2 0 01-2-2V4a2 2 0 012-2h9a2 2 0 012 2v1"/>
                      </svg>
                    </button>
                  </div>
                </dd>
              </div>
              <div v-else class="detail-row">
                <dt>Routing</dt>
                <dd>
                  <span v-if="isWeighted" class="badge badge-warning">Canary</span>
                  <span v-else class="badge badge-info">Round Robin</span>
                </dd>
              </div>
              <div class="detail-row">
                <dt>Rewrite</dt>
                <dd class="text-mono">{{ config.rewrite || '-' }}</dd>
              </div>
              <div class="detail-row">
                <dt>Hosts</dt>
                <dd>
                  <template v-if="hostsList.length">
                    <span v-for="h in hostsList" :key="h" class="badge badge-info detail-badge">{{ h }}</span>
                  </template>
                  <span v-else class="text-muted">Any</span>
                </dd>
              </div>
              <div class="detail-row">
                <dt>Methods</dt>
                <dd>
                  <template v-if="methodsList.length">
                    <span v-for="m in methodsList" :key="m" :class="['badge', 'detail-badge', methodBadge(m)]">{{ m }}</span>
                  </template>
                  <span v-else class="text-muted">All</span>
                </dd>
              </div>
              <div v-if="config.priority" class="detail-row">
                <dt>Priority</dt>
                <dd>{{ config.priority }}</dd>
              </div>
              <div class="detail-row">
                <dt>Metrics</dt>
                <dd>
                  <span :class="['badge', config.disableMetrics ? 'badge-neutral' : 'badge-success']">
                    {{ config.disableMetrics ? 'Disabled' : 'Enabled' }}
                  </span>
                </dd>
              </div>
            </dl>
          </div>
        </div>

        <!-- Backends (Load Balancing) -->
        <div v-if="backendsList.length" class="card">
          <div class="card-header">
            <div class="card-header-title">
              <svg width="16" height="16" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <rect x="2" y="2" width="20" height="8" rx="2" stroke-width="1.5"/><rect x="2" y="14" width="20" height="8" rx="2" stroke-width="1.5"/><circle cx="6" cy="6" r="1" fill="currentColor"/><circle cx="6" cy="18" r="1" fill="currentColor"/>
              </svg>
              <h2>Backends</h2>
            </div>
            <div style="display: flex; gap: 4px;">
              <span v-if="isWeighted" class="badge badge-warning">Canary</span>
              <span v-else class="badge badge-info">Round Robin</span>
              <span class="badge badge-neutral">{{ backendsList.length }} endpoint{{ backendsList.length !== 1 ? 's' : '' }}</span>
            </div>
          </div>
          <div class="card-body">
            <div class="backends-list">
              <div v-for="(b, i) in backendsList" :key="i" class="backend-item">
                <div class="backend-info">
                  <div class="backend-endpoint-row">
                    <span class="backend-index">{{ i + 1 }}</span>
                    <span class="backend-endpoint text-mono">{{ b.endpoint }}</span>
                    <button class="copy-btn" @click="copyToClipboard(b.endpoint)" title="Copy endpoint">
                      <svg width="13" height="13" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <rect x="9" y="9" width="13" height="13" rx="2" stroke-width="2"/><path stroke-width="2" d="M5 15H4a2 2 0 01-2-2V4a2 2 0 012-2h9a2 2 0 012 2v1"/>
                      </svg>
                    </button>
                  </div>
                  <span v-if="isWeighted" class="backend-weight-label">{{ b.weight }}%</span>
                </div>
                <div v-if="isWeighted" class="backend-bar-track">
                  <div class="backend-bar-fill" :style="{ width: b.weight + '%' }"></div>
                </div>
              </div>
            </div>
          </div>
        </div>

        <!-- Middlewares -->
        <div v-if="middlewaresList.length" class="card">
          <div class="card-header">
            <div class="card-header-title">
              <svg width="16" height="16" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 6V4m0 2a2 2 0 100 4m0-4a2 2 0 110 4m-6 8a2 2 0 100-4m0 4a2 2 0 110-4m0 4v2m0-6V4m6 6v10m6-2a2 2 0 100-4m0 4a2 2 0 110-4m0 4v2m0-6V4"/>
              </svg>
              <h2>Middlewares</h2>
            </div>
            <span class="badge badge-info">{{ middlewaresList.length }}</span>
          </div>
          <div class="card-body">
            <div class="middleware-list">
              <router-link
                v-for="(mw, i) in middlewaresList"
                :key="mw"
                :to="`/middlewares?q=${encodeURIComponent(mw)}`"
                class="middleware-item middleware-item--link"
              >
                <div class="middleware-left">
                  <span class="middleware-order">{{ i + 1 }}</span>
                  <span class="middleware-name text-mono">{{ mw }}</span>
                </div>
                <svg width="16" height="16" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7" />
                </svg>
              </router-link>
            </div>
          </div>
        </div>

        <!-- Health Check -->
        <div v-if="healthCheck" class="card">
          <div class="card-header">
            <div class="card-header-title">
              <svg width="16" height="16" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4.318 6.318a4.5 4.5 0 000 6.364L12 20.364l7.682-7.682a4.5 4.5 0 00-6.364-6.364L12 7.636l-1.318-1.318a4.5 4.5 0 00-6.364 0z" />
              </svg>
              <h2>Health Check</h2>
            </div>
            <span class="badge badge-success">Active</span>
          </div>
          <div class="card-body">
            <dl class="detail-grid">
              <div class="detail-row">
                <dt>Path</dt>
                <dd class="text-mono">{{ healthCheck.path || '-' }}</dd>
              </div>
              <div class="detail-row">
                <dt>Interval</dt>
                <dd>
                  <span class="health-value">{{ healthCheck.interval || '-' }}</span>
                </dd>
              </div>
              <div class="detail-row">
                <dt>Timeout</dt>
                <dd>
                  <span class="health-value">{{ healthCheck.timeout || '-' }}</span>
                </dd>
              </div>
              <div v-if="healthCheck.healthyStatuses?.length" class="detail-row">
                <dt>Healthy Statuses</dt>
                <dd>
                  <span v-for="s in healthCheck.healthyStatuses" :key="s" class="badge badge-success detail-badge">{{ s }}</span>
                </dd>
              </div>
            </dl>
          </div>
        </div>

        <!-- Security -->
        <div v-if="security" class="card">
          <div class="card-header">
            <div class="card-header-title">
              <svg width="16" height="16" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m5.618-4.016A11.955 11.955 0 0112 2.944a11.955 11.955 0 01-8.618 3.04A12.02 12.02 0 003 9c0 5.591 3.824 10.29 9 11.622 5.176-1.332 9-6.03 9-11.622 0-1.042-.133-2.052-.382-3.016z" />
              </svg>
              <h2>Security</h2>
            </div>
          </div>
          <div class="card-body">
            <dl class="detail-grid">
              <div class="detail-row">
                <dt>Forward Host Headers</dt>
                <dd>
                  <span :class="['badge', security.forwardHostHeaders !== false ? 'badge-success' : 'badge-neutral']">
                    {{ security.forwardHostHeaders !== false ? 'Yes' : 'No' }}
                  </span>
                </dd>
              </div>
              <div class="detail-row">
                <dt>Exploit Protection</dt>
                <dd>
                  <span :class="['badge', security.enableExploitProtection ? 'badge-success' : 'badge-neutral']">
                    {{ security.enableExploitProtection ? 'Enabled' : 'Disabled' }}
                  </span>
                </dd>
              </div>
              <div v-if="security.tls" class="detail-row">
                <dt>TLS Skip Verify</dt>
                <dd>
                  <span :class="['badge', security.tls.insecureSkipVerify ? 'badge-warning' : 'badge-success']">
                    {{ security.tls.insecureSkipVerify ? 'Yes (insecure)' : 'No' }}
                  </span>
                </dd>
              </div>
            </dl>
          </div>
        </div>

        <!-- Extra fields (catch-all for anything unknown) -->
        <div v-if="hasExtraFields" class="card">
          <div class="card-header">
            <div class="card-header-title">
              <svg width="16" height="16" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 6h16M4 12h16M4 18h7" />
              </svg>
              <h2>Additional Configuration</h2>
            </div>
          </div>
          <div class="card-body">
            <pre class="code-block">{{ formatExtraFields() }}</pre>
          </div>
        </div>

        <!-- Raw Configuration -->
        <div class="card">
          <div class="card-header">
            <div class="card-header-title">
              <svg width="16" height="16" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10 20l4-16m4 4l4 4-4 4M6 16l-4-4 4-4" />
              </svg>
              <h2>Raw Configuration</h2>
            </div>
            <button class="btn btn-secondary btn-sm" @click="showRaw = !showRaw">
              {{ showRaw ? 'Hide' : 'Show' }}
            </button>
          </div>
          <div v-if="showRaw" class="card-body">
            <pre class="code-block">{{ JSON.stringify(config, null, 2) }}</pre>
          </div>
        </div>

        <!-- Metadata -->
        <div class="card">
          <div class="card-header">
            <div class="card-header-title">
              <svg width="16" height="16" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
              </svg>
              <h2>Metadata</h2>
            </div>
          </div>
          <div class="card-body">
            <dl class="detail-grid">
              <div class="detail-row">
                <dt>Route ID</dt>
                <dd class="text-mono">{{ route.id }}</dd>
              </div>
              <div class="detail-row">
                <dt>Created</dt>
                <dd>
                  <span>{{ formatDate(route.createdAt) }}</span>
                  <span class="relative-time">{{ relativeTime(route.createdAt) }}</span>
                </dd>
              </div>
              <div class="detail-row">
                <dt>Last Updated</dt>
                <dd>
                  <span>{{ formatDate(route.updatedAt) }}</span>
                  <span class="relative-time">{{ relativeTime(route.updatedAt) }}</span>
                </dd>
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
import { useRouter } from 'vue-router'
import { routesApi, type Route } from '@/api/routes'
import { useConfirm } from '@/composables/useConfirm'
import { useNotificationStore } from '@/stores/notification'

const props = defineProps<{ id: string }>()
const router = useRouter()
const { confirm } = useConfirm()
const notify = useNotificationStore()

const loading = ref(true)
const error = ref(false)
const route = ref<Route | null>(null)
const showRaw = ref(false)

const knownFieldKeys = new Set([
  'path', 'target', 'hosts', 'methods', 'rewrite', 'enabled',
  'priority', 'disableMetrics', 'middlewares', 'healthCheck', 'security',
  'backends',
])

interface Backend {
  endpoint: string
  weight?: number
}

const config = computed(() => route.value?.config || {})

const hostsList = computed(() => {
  const hosts = config.value.hosts
  return Array.isArray(hosts) ? hosts.map(String) : []
})

const methodsList = computed(() => {
  const methods = config.value.methods
  return Array.isArray(methods) ? methods.map(String) : []
})

const backendsList = computed((): Backend[] => {
  const backends = config.value.backends
  if (!Array.isArray(backends)) return []
  return backends.map((b: unknown) => {
    if (typeof b === 'string') return { endpoint: b }
    if (b && typeof b === 'object') {
      const obj = b as Record<string, unknown>
      return {
        endpoint: String(obj.endpoint || obj.url || obj.target || ''),
        weight: typeof obj.weight === 'number' ? obj.weight : undefined,
      }
    }
    return { endpoint: String(b) }
  })
})

const isWeighted = computed(() =>
  backendsList.value.some(b => b.weight != null && b.weight > 0)
)

const middlewaresList = computed(() => {
  const mw = config.value.middlewares
  return Array.isArray(mw) ? mw.map(String) : []
})

interface HealthCheckConfig {
  path?: string
  interval?: string
  timeout?: string
  healthyStatuses?: number[]
}

interface SecurityConfig {
  forwardHostHeaders?: boolean
  enableExploitProtection?: boolean
  tls?: { insecureSkipVerify?: boolean }
}

const healthCheck = computed((): HealthCheckConfig | null => {
  const hc = config.value.healthCheck
  return hc && typeof hc === 'object' ? hc as HealthCheckConfig : null
})

const security = computed((): SecurityConfig | null => {
  const sec = config.value.security
  return sec && typeof sec === 'object' ? sec as SecurityConfig : null
})

const extraFields = computed(() => {
  const result: Record<string, unknown> = {}
  for (const [key, value] of Object.entries(config.value)) {
    if (!knownFieldKeys.has(key)) {
      result[key] = value
    }
  }
  return result
})

const hasExtraFields = computed(() => Object.keys(extraFields.value).length > 0)

function hostUrl(host: string): string {
  const base = (host.startsWith('http://') || host.startsWith('https://')) ? host : `http://${host}`
  const path = (config.value.path as string) || ''
  const cleanPath = path.replace(/\/?\*$/, '')
  return base.replace(/\/$/, '') + cleanPath
}

function methodBadge(method: string): string {
  switch (method) {
    case 'GET': return 'badge-success'
    case 'POST': return 'badge-primary'
    case 'PUT': case 'PATCH': return 'badge-warning'
    case 'DELETE': return 'badge-danger'
    default: return 'badge-info'
  }
}

function formatExtraFields(): string {
  return JSON.stringify(extraFields.value, null, 2)
}

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

function relativeTime(dateStr: string): string {
  if (!dateStr) return ''
  const now = Date.now()
  const then = new Date(dateStr).getTime()
  const diffMs = now - then
  const diffSec = Math.floor(diffMs / 1000)
  const diffMin = Math.floor(diffSec / 60)
  const diffHr = Math.floor(diffMin / 60)
  const diffDays = Math.floor(diffHr / 24)

  if (diffSec < 60) return 'just now'
  if (diffMin < 60) return `${diffMin}m ago`
  if (diffHr < 24) return `${diffHr}h ago`
  if (diffDays < 30) return `${diffDays}d ago`
  const diffMonths = Math.floor(diffDays / 30)
  if (diffMonths < 12) return `${diffMonths}mo ago`
  const diffYears = Math.floor(diffDays / 365)
  return `${diffYears}y ago`
}

async function copyToClipboard(text: unknown) {
  if (!text) return
  try {
    await navigator.clipboard.writeText(String(text))
    notify.success('Copied to clipboard')
  } catch {
    // fallback: silently fail
  }
}

function handleEdit() {
  router.push(`/routes?edit=${route.value?.id}`)
}

async function handleDelete() {
  if (!route.value) return
  const confirmed = await confirm({
    title: 'Delete Route',
    message: `Are you sure you want to delete "${route.value.name}"? This action cannot be undone.`,
    confirmText: 'Delete',
    variant: 'danger',
  })
  if (!confirmed) return
  try {
    await routesApi.delete(route.value.id)
    notify.success(`Route "${route.value.name}" deleted`)
    router.push('/routes')
  } catch {
    notify.error('Failed to delete route')
  }
}

onMounted(async () => {
  try {
    const res = await routesApi.get(Number(props.id))
    route.value = res.data
  } catch {
    error.value = true
  } finally {
    loading.value = false
  }
})
</script>

<style scoped>
/* Header */
.detail-header {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-bottom: 12px;
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
.detail-header-actions {
  display: flex;
  gap: 8px;
  flex-shrink: 0;
}
.back-link {
  display: flex;
  align-items: center;
  color: var(--text-muted);
  transition: color var(--transition);
}
.back-link:hover { color: var(--text-secondary); }

/* Status indicator */
.status-indicator {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  font-size: 13px;
  font-weight: 500;
  padding: 4px 10px;
  border-radius: 9999px;
  flex-shrink: 0;
}
.status-dot {
  width: 8px;
  height: 8px;
  border-radius: 50%;
  flex-shrink: 0;
}
.status-online {
  color: var(--success-600);
  background: var(--success-50);
}
.status-online .status-dot {
  background: var(--success-500);
  animation: pulse-green 2s infinite;
}
.status-offline {
  color: var(--danger-500);
  background: var(--danger-50);
}
.status-offline .status-dot {
  background: var(--danger-500);
  animation: pulse-red 2s infinite;
}
@keyframes pulse-green {
  0% { box-shadow: 0 0 0 0 rgba(34, 197, 94, 0.6); }
  70% { box-shadow: 0 0 0 6px rgba(34, 197, 94, 0); }
  100% { box-shadow: 0 0 0 0 rgba(34, 197, 94, 0); }
}
@keyframes pulse-red {
  0% { box-shadow: 0 0 0 0 rgba(239, 68, 68, 0.6); }
  70% { box-shadow: 0 0 0 6px rgba(239, 68, 68, 0); }
  100% { box-shadow: 0 0 0 0 rgba(239, 68, 68, 0); }
}

/* Quick Stats */
.quick-stats {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 10px 16px;
  margin-bottom: 20px;
  background: var(--bg-secondary);
  border: 1px solid var(--border-primary);
  border-radius: var(--radius);
  flex-wrap: wrap;
  font-size: 13px;
  color: var(--text-secondary);
}
.quick-stat {
  display: flex;
  align-items: center;
  gap: 6px;
}
.quick-stat svg {
  color: var(--text-muted);
  flex-shrink: 0;
}
.quick-stat-divider {
  width: 1px;
  height: 18px;
  background: var(--border-primary);
  flex-shrink: 0;
}
.badge-sm {
  font-size: 11px;
  padding: 1px 7px;
}

/* Animated Traffic Flow */
.traffic-flow {
  margin-bottom: 20px;
  padding: 24px;
  overflow-x: auto;
}
.traffic-flow-inner {
  display: flex;
  align-items: center;
  gap: 0;
  min-width: min-content;
}

/* Flow nodes */
.tf-node {
  display: flex;
  flex-direction: column;
  gap: 6px;
  padding: 14px 18px;
  border: 1.5px solid var(--border-primary);
  border-radius: 10px;
  background: var(--bg-secondary);
  min-width: 130px;
  position: relative;
  transition: border-color 0.3s, box-shadow 0.3s;
}
.tf-node:hover {
  border-color: var(--primary-300);
  box-shadow: 0 0 0 3px var(--primary-50, rgba(147, 51, 234, 0.06));
}
.tf-node-icon {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 36px;
  height: 36px;
  border-radius: 8px;
  flex-shrink: 0;
  margin-bottom: 2px;
}
.tf-node-body {
  display: flex;
  flex-direction: column;
  gap: 3px;
}
.tf-node-label {
  font-size: 10px;
  font-weight: 700;
  text-transform: uppercase;
  letter-spacing: 0.08em;
  color: var(--text-muted);
}
.tf-node-detail {
  font-size: 13px;
  color: var(--text-primary);
  word-break: break-all;
}
.tf-node-detail--copyable {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  cursor: pointer;
  border-radius: 4px;
  padding: 1px 4px;
  margin: -1px -4px;
  transition: background 0.15s;
}
.tf-node-detail--copyable:hover {
  background: var(--bg-tertiary, rgba(0,0,0,0.04));
}
.tf-node-detail--copyable .copy-icon {
  opacity: 0;
  transition: opacity 0.15s;
  color: var(--text-muted);
  flex-shrink: 0;
}
.tf-node-detail--copyable:hover .copy-icon {
  opacity: 1;
}

/* Node variants */
.tf-node--client {
  border-color: var(--primary-300);
  background: var(--primary-50, rgba(147, 51, 234, 0.04));
}
.tf-node--client .tf-node-icon {
  background: var(--primary-100, rgba(147, 51, 234, 0.1));
  color: var(--primary-600);
}

.tf-node--path {
  border-color: var(--primary-300);
  background: var(--primary-50, rgba(147, 51, 234, 0.04));
}
.tf-node--path .tf-node-icon {
  background: var(--primary-100, rgba(147, 51, 234, 0.1));
  color: var(--primary-600);
}

.tf-node--target {
  border-color: var(--success-500, #22c55e);
  background: var(--success-50, rgba(34, 197, 94, 0.05));
}
.tf-node--target .tf-node-icon {
  background: var(--success-100, rgba(34, 197, 94, 0.12));
  color: var(--success-600, #16a34a);
}

/* Rewrite badge */
.tf-rewrite {
  display: inline-flex;
  align-items: center;
  gap: 4px;
  margin-top: 4px;
  padding: 2px 8px;
  border-radius: 4px;
  font-size: 11px;
  font-family: 'SF Mono', 'Fira Code', Menlo, Consolas, monospace;
  color: var(--primary-600);
  background: var(--primary-100, rgba(147, 51, 234, 0.08));
}

/* Connectors with animated packets */
.tf-connector {
  display: flex;
  align-items: center;
  padding: 0 4px;
  flex-shrink: 0;
}
.tf-connector-line {
  position: relative;
  width: 48px;
  height: 2px;
  background: var(--border-primary);
  border-radius: 1px;
  overflow: hidden;
}

/* Animated traffic packets */
.tf-packet, .tf-packet-sm {
  position: absolute;
  top: -2px;
  width: 6px;
  height: 6px;
  border-radius: 50%;
  background: var(--primary-500);
  opacity: 0;
  animation: packet-flow 2.4s ease-in-out infinite;
  box-shadow: 0 0 6px var(--primary-400, rgba(147, 51, 234, 0.5));
}
.tf-packet--1 { animation-delay: 0s; }
.tf-packet--2 { animation-delay: 0.8s; }
.tf-packet--3 { animation-delay: 1.6s; }

.tf-packet-sm {
  width: 4px;
  height: 4px;
  top: -1px;
  animation-duration: 1.6s;
}

@keyframes packet-flow {
  0% {
    left: -6px;
    opacity: 0;
  }
  10% {
    opacity: 1;
  }
  90% {
    opacity: 1;
  }
  100% {
    left: calc(100% + 2px);
    opacity: 0;
  }
}

/* Middlewares group */
.tf-middlewares-group {
  display: flex;
  flex-direction: column;
  gap: 8px;
}
.tf-middlewares-label {
  display: flex;
  align-items: center;
  gap: 5px;
  font-size: 10px;
  font-weight: 700;
  text-transform: uppercase;
  letter-spacing: 0.08em;
  color: var(--text-muted);
  padding-left: 4px;
}
.tf-middlewares-pipeline {
  display: flex;
  align-items: center;
  padding: 10px 14px;
  border: 1.5px dashed var(--warning-300, #fcd34d);
  border-radius: 10px;
  background: var(--warning-50, rgba(245, 158, 11, 0.04));
  gap: 0;
}

/* Middleware steps */
.tf-mw-step {
  display: flex;
  align-items: center;
}
.tf-node--middleware {
  display: flex;
  flex-direction: row;
  align-items: center;
  gap: 8px;
  padding: 8px 12px;
  border: 1px solid var(--warning-300, #fcd34d);
  border-radius: 8px;
  background: var(--bg-primary);
  min-width: auto;
  white-space: nowrap;
}
.tf-node--middleware:hover {
  border-color: var(--warning-400, #fbbf24);
  box-shadow: 0 0 0 3px var(--warning-50, rgba(245, 158, 11, 0.08));
}
.tf-mw-order {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 20px;
  height: 20px;
  border-radius: 50%;
  background: var(--warning-100, rgba(245, 158, 11, 0.15));
  color: var(--warning-700, #b45309);
  font-size: 11px;
  font-weight: 700;
  flex-shrink: 0;
}
.tf-mw-name {
  font-size: 12px;
  color: var(--text-primary);
}
.tf-mw-connector {
  display: flex;
  align-items: center;
  padding: 0 2px;
}
.tf-mw-connector-line {
  position: relative;
  width: 24px;
  height: 2px;
  background: var(--warning-200, #fde68a);
  border-radius: 1px;
  overflow: hidden;
}
.tf-mw-connector-line .tf-packet-sm {
  background: var(--warning-500, #f59e0b);
  box-shadow: 0 0 6px var(--warning-400, rgba(245, 158, 11, 0.5));
}

/* Backends group */
.tf-backends-group {
  display: flex;
  flex-direction: column;
  gap: 8px;
}
.tf-backends-label {
  display: flex;
  align-items: center;
  gap: 5px;
  font-size: 10px;
  font-weight: 700;
  text-transform: uppercase;
  letter-spacing: 0.08em;
  color: var(--text-muted);
  padding-left: 4px;
}
.tf-backends-list {
  display: flex;
  flex-direction: column;
  gap: 6px;
}
.tf-node--backend {
  display: flex;
  flex-direction: column;
  gap: 6px;
  padding: 10px 14px;
  border-color: var(--success-500, #22c55e);
  background: var(--success-50, rgba(34, 197, 94, 0.05));
  min-width: 180px;
}
.tf-node--backend:hover {
  border-color: var(--success-400, #4ade80);
  box-shadow: 0 0 0 3px var(--success-50, rgba(34, 197, 94, 0.1));
}
.tf-backend-header {
  display: flex;
  align-items: center;
  gap: 8px;
  color: var(--success-600, #16a34a);
}
.tf-backend-header .text-mono {
  font-size: 13px;
  color: var(--text-primary);
}
.tf-backend-weight {
  display: flex;
  align-items: center;
  gap: 8px;
}
.tf-backend-weight-track {
  flex: 1;
  height: 4px;
  border-radius: 2px;
  background: var(--border-primary);
  overflow: hidden;
}
.tf-backend-weight-fill {
  height: 100%;
  border-radius: 2px;
  background: var(--success-500, #22c55e);
  min-width: 4px;
  transition: width 0.3s ease;
}
.tf-backend-weight-label {
  font-size: 11px;
  font-weight: 600;
  color: var(--text-muted);
  white-space: nowrap;
}

/* Host links */
.flow-badge {
  margin-right: 4px;
  margin-bottom: 2px;
}
.flow-host-link {
  text-decoration: none;
  cursor: pointer;
  transition: opacity 0.15s;
}
.flow-host-link:hover {
  opacity: 0.8;
  text-decoration: underline;
}

/* Cards layout */
.detail-cards {
  display: flex;
  flex-direction: column;
  gap: 20px;
  max-width: 800px;
}

/* Card header with icon */
.card-header-title {
  display: flex;
  align-items: center;
  gap: 8px;
}
.card-header-title svg {
  color: var(--text-muted);
  flex-shrink: 0;
}

/* Detail grid */
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
  width: 160px;
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
.detail-badge {
  margin-right: 4px;
  margin-bottom: 2px;
}
.text-mono {
  font-family: 'SF Mono', 'Fira Code', 'Fira Mono', Menlo, Consolas, monospace;
  font-size: 13px;
}
.text-muted {
  color: var(--text-muted);
  font-style: italic;
}

/* Copy button */
.value-with-copy {
  display: inline-flex;
  align-items: center;
  gap: 6px;
}
.copy-btn {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 24px;
  height: 24px;
  border: none;
  background: none;
  border-radius: 4px;
  color: var(--text-muted);
  cursor: pointer;
  opacity: 0;
  transition: opacity 0.15s, background 0.15s, color 0.15s;
}
.copy-btn:hover {
  background: var(--bg-tertiary, rgba(0,0,0,0.06));
  color: var(--text-secondary);
}
.detail-row:hover .copy-btn,
.backend-item:hover .copy-btn {
  opacity: 1;
}

/* Relative time */
.relative-time {
  margin-left: 8px;
  font-size: 12px;
  color: var(--text-muted);
  font-style: italic;
}

/* Health check values */
.health-value {
  font-family: 'SF Mono', 'Fira Code', 'Fira Mono', Menlo, Consolas, monospace;
  font-size: 13px;
  color: var(--primary-600);
  font-weight: 500;
}

/* Middlewares */
.middleware-list {
  display: flex;
  flex-direction: column;
  gap: 6px;
}
.middleware-item {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
  padding: 10px 14px;
  border: 1px solid var(--border-primary);
  border-radius: var(--radius);
  background: var(--bg-secondary);
}
.middleware-item--link {
  text-decoration: none;
  color: inherit;
  cursor: pointer;
  transition: border-color 0.15s, background 0.15s;
}
.middleware-item--link:hover {
  border-color: var(--primary-300);
  background: var(--primary-50, rgba(147, 51, 234, 0.04));
}
.middleware-item--link svg {
  color: var(--text-muted);
  opacity: 0;
  transition: opacity 0.15s;
  flex-shrink: 0;
}
.middleware-item--link:hover svg {
  opacity: 1;
}
.middleware-left {
  display: flex;
  align-items: center;
  gap: 12px;
}
.middleware-order {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 24px;
  height: 24px;
  border-radius: 50%;
  background: var(--primary-100, rgba(147, 51, 234, 0.1));
  color: var(--primary-700);
  font-size: 12px;
  font-weight: 700;
  flex-shrink: 0;
}
.middleware-name {
  font-size: 14px;
  color: var(--text-primary);
}

/* Backends */
.backends-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
}
.backend-item {
  display: flex;
  flex-direction: column;
  gap: 6px;
  padding: 10px 14px;
  border: 1px solid var(--border-primary);
  border-radius: var(--radius);
  background: var(--bg-secondary);
}
.backend-info {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
}
.backend-endpoint-row {
  display: flex;
  align-items: center;
  gap: 10px;
}
.backend-index {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 22px;
  height: 22px;
  border-radius: 50%;
  background: var(--bg-tertiary);
  color: var(--text-muted);
  font-size: 11px;
  font-weight: 700;
  flex-shrink: 0;
}
.backend-endpoint {
  font-size: 14px;
  color: var(--text-primary);
}
.backend-weight-label {
  font-size: 13px;
  font-weight: 600;
  color: var(--text-secondary);
  flex-shrink: 0;
}
.backend-bar-track {
  height: 6px;
  border-radius: 3px;
  background: var(--border-primary);
  overflow: hidden;
}
.backend-bar-fill {
  height: 100%;
  border-radius: 3px;
  background: var(--primary-500);
  transition: width 0.3s ease;
  min-width: 4px;
}
</style>
