<template>
  <div class="dashboard-layout" :class="{ 'sidebar-collapsed': sidebarCollapsed }">
    <!-- Sidebar -->
    <aside class="sidebar">
      <div class="sidebar-header">
        <img src="/logo.png" alt="Goma" class="sidebar-logo" />
        <span class="sidebar-brand-text">Goma Admin</span>

        <button class="sidebar-collapse-btn" @click="toggleSidebar" :title="sidebarCollapsed ? 'Expand sidebar' : 'Collapse sidebar'">
          <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
            <polyline :points="sidebarCollapsed ? '9 18 15 12 9 6' : '15 18 9 12 15 6'" />
          </svg>
        </button>
      </div>

      <!-- Instance Switcher -->
      <div v-if="instanceStore.instances.length > 0" class="instance-switcher">
        <div class="instance-switcher-toggle" @click="instanceDropdownOpen = !instanceDropdownOpen">
          <div style="display: flex; align-items: center; gap: 8px; overflow: hidden;">
            <span class="instance-avatar">{{ activeInstanceInitial }}</span>
            <span class="instance-name nav-label">{{ activeInstanceName }}</span>
          </div>
          <svg class="nav-label" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
            <polyline points="6 9 12 15 18 9" />
          </svg>
        </div>
        <Transition name="dropdown">
          <div v-if="instanceDropdownOpen" class="instance-dropdown">
            <div
              class="instance-option"
              :class="{ active: !instanceStore.currentInstanceId }"
              @click="switchInstance(null)"
            >
              <span class="instance-avatar">A</span>
              <span>All Instances</span>
            </div>
            <div
              v-for="inst in instanceStore.instances"
              :key="inst.id"
              class="instance-option"
              :class="{ active: instanceStore.currentInstanceId === inst.id }"
              @click="switchInstance(inst.id)"
            >
              <span class="instance-avatar">{{ inst.name?.charAt(0)?.toUpperCase() || '?' }}</span>
              <span>{{ inst.name }}</span>
            </div>
          </div>
        </Transition>
      </div>

      <nav class="sidebar-nav">
        <div class="nav-section">
          <router-link to="/" class="nav-item" :class="{ active: $route.path === '/' }">
            <svg class="nav-icon" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
              <rect x="3" y="3" width="7" height="7" /><rect x="14" y="3" width="7" height="7" /><rect x="3" y="14" width="7" height="7" /><rect x="14" y="14" width="7" height="7" />
            </svg>
            <span class="nav-label">Dashboard</span>
          </router-link>
          <router-link to="/instances" class="nav-item" :class="{ active: $route.path.startsWith('/instances') }">
            <svg class="nav-icon" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
              <rect x="2" y="2" width="20" height="8" rx="2" /><rect x="2" y="14" width="20" height="8" rx="2" /><circle cx="6" cy="6" r="1" /><circle cx="6" cy="18" r="1" />
            </svg>
            <span class="nav-label">Instances</span>
          </router-link>
          <router-link to="/routes" class="nav-item" :class="{ active: $route.path.startsWith('/routes') }">
            <svg class="nav-icon" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
              <circle cx="6" cy="6" r="3" /><circle cx="18" cy="18" r="3" /><path d="M6 9a9 9 0 0 0 9 9" />
            </svg>
            <span class="nav-label">Routes</span>
          </router-link>
          <router-link to="/middlewares" class="nav-item" :class="{ active: $route.path.startsWith('/middlewares') }">
            <svg class="nav-icon" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
              <line x1="8" y1="6" x2="21" y2="6" /><line x1="8" y1="12" x2="21" y2="12" /><line x1="8" y1="18" x2="21" y2="18" /><line x1="3" y1="6" x2="3.01" y2="6" /><line x1="3" y1="12" x2="3.01" y2="12" /><line x1="3" y1="18" x2="3.01" y2="18" />
            </svg>
            <span class="nav-label">Middlewares</span>
          </router-link>
          <router-link to="/audit" class="nav-item" :class="{ active: $route.path.startsWith('/audit') }">
            <svg class="nav-icon" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
              <path d="M9 5H7a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2m-3 7h3m-3 4h3m-6-4h.01M9 16h.01" />
            </svg>
            <span class="nav-label">Audit Log</span>
          </router-link>
          <router-link to="/api-keys" class="nav-item" :class="{ active: $route.path.startsWith('/api-keys') }">
            <svg class="nav-icon" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
              <path d="M21 2l-2 2m-7.61 7.61a5.5 5.5 0 1 1-7.778 7.778 5.5 5.5 0 0 1 7.777-7.777zm0 0L15.5 7.5m0 0l3 3L22 7l-3-3m-3.5 3.5L19 4" />
            </svg>
            <span class="nav-label">API Keys</span>
          </router-link>
          <router-link to="/profile" class="nav-item" :class="{ active: $route.path === '/profile' }">
            <svg class="nav-icon" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
              <path d="M20 21v-2a4 4 0 0 0-4-4H8a4 4 0 0 0-4 4v2" /><circle cx="12" cy="7" r="4" />
            </svg>
            <span class="nav-label">Profile</span>
          </router-link>
          <router-link to="/settings" class="nav-item" :class="{ active: $route.path.startsWith('/settings') }">
            <svg class="nav-icon" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
              <circle cx="12" cy="12" r="3" /><path d="M19.4 15a1.65 1.65 0 0 0 .33 1.82l.06.06a2 2 0 1 1-2.83 2.83l-.06-.06a1.65 1.65 0 0 0-1.82-.33 1.65 1.65 0 0 0-1 1.51V21a2 2 0 0 1-4 0v-.09A1.65 1.65 0 0 0 9 19.4a1.65 1.65 0 0 0-1.82.33l-.06.06a2 2 0 1 1-2.83-2.83l.06-.06A1.65 1.65 0 0 0 4.68 15a1.65 1.65 0 0 0-1.51-1H3a2 2 0 0 1 0-4h.09A1.65 1.65 0 0 0 4.6 9a1.65 1.65 0 0 0-.33-1.82l-.06-.06a2 2 0 1 1 2.83-2.83l.06.06A1.65 1.65 0 0 0 9 4.68a1.65 1.65 0 0 0 1-1.51V3a2 2 0 0 1 4 0v.09a1.65 1.65 0 0 0 1 1.51 1.65 1.65 0 0 0 1.82-.33l.06-.06a2 2 0 1 1 2.83 2.83l-.06.06A1.65 1.65 0 0 0 19.4 9a1.65 1.65 0 0 0 1.51 1H21a2 2 0 0 1 0 4h-.09a1.65 1.65 0 0 0-1.51 1z" />
            </svg>
            <span class="nav-label">Settings</span>
          </router-link>
        </div>
        <!-- Admin section -->
        <div v-if="authStore.isAdmin" class="nav-section">
          <div class="nav-section-title">Admin</div>
          <router-link to="/users" class="nav-item" :class="{ active: $route.path.startsWith('/users') }">
            <svg class="nav-icon" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
              <path d="M17 21v-2a4 4 0 0 0-4-4H5a4 4 0 0 0-4 4v2" /><circle cx="9" cy="7" r="4" /><path d="M23 21v-2a4 4 0 0 0-3-3.87" /><path d="M16 3.13a4 4 0 0 1 0 7.75" />
            </svg>
            <span class="nav-label">Users</span>
          </router-link>
          <router-link to="/oauth-settings" class="nav-item" :class="{ active: $route.path.startsWith('/oauth-settings') }">
            <svg class="nav-icon" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
              <path d="M15 7a2 2 0 0 1 2 2m4 0a6 6 0 0 1-7.743 5.743L11 17H9v2H7v2H4a1 1 0 0 1-1-1v-2.586a1 1 0 0 1 .293-.707l5.964-5.964A6 6 0 1 1 21 9z" />
            </svg>
            <span class="nav-label">OAuth</span>
          </router-link>
        </div>
      </nav>

      <!-- Sidebar bottom -->
      <div class="sidebar-footer">
        <div class="nav-section">
          <div class="nav-section-title">Docs</div>
          <a href="/docs" target="_blank" class="nav-item">
            <svg class="nav-icon" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
              <path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z" /><polyline points="14 2 14 8 20 8" /><line x1="16" y1="13" x2="8" y2="13" /><line x1="16" y1="17" x2="8" y2="17" /><polyline points="10 9 9 9 8 9" />
            </svg>
            <span class="nav-label">Swagger UI</span>
          </a>
          <a href="/redoc" target="_blank" class="nav-item">
            <svg class="nav-icon" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
              <path d="M4 19.5A2.5 2.5 0 0 1 6.5 17H20" /><path d="M6.5 2H20v20H6.5A2.5 2.5 0 0 1 4 19.5v-15A2.5 2.5 0 0 1 6.5 2z" />
            </svg>
            <span class="nav-label">Redoc</span>
          </a>
        </div>
        <div class="sidebar-version">
          <span class="nav-label">v1.0.0</span>
        </div>
      </div>
    </aside>

    <!-- Main content area -->
    <div class="main-wrapper">
      <!-- Top bar -->
      <header class="topbar">
        <div class="topbar-left">
          <!-- Mobile menu button -->
          <button class="mobile-menu-btn" @click="mobileOpen = true">
            <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
              <line x1="3" y1="12" x2="21" y2="12" /><line x1="3" y1="6" x2="21" y2="6" /><line x1="3" y1="18" x2="21" y2="18" />
            </svg>
          </button>
        </div>
        <div class="topbar-right">
          <!-- User menu -->
          <div class="user-menu" ref="userMenuRef">
            <button class="user-menu-trigger" @click="userMenuOpen = !userMenuOpen">
              <span class="user-avatar">{{ userInitial }}</span>
              <span class="user-name">{{ authStore.user?.name || authStore.user?.email || 'User' }}</span>
              <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                <polyline points="6 9 12 15 18 9" />
              </svg>
            </button>
            <Transition name="dropdown">
              <div v-if="userMenuOpen" class="user-dropdown">
                <div class="user-dropdown-header">
                  <span class="user-avatar user-avatar-lg">{{ userInitial }}</span>
                  <div class="user-dropdown-info">
                    <div class="user-dropdown-name">{{ authStore.user?.name || 'User' }}</div>
                    <div class="user-dropdown-email">{{ authStore.user?.email }}</div>
                  </div>
                </div>
                <div class="user-dropdown-divider" />
                <router-link to="/profile" class="user-dropdown-item" @click="userMenuOpen = false">
                  <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                    <path d="M20 21v-2a4 4 0 0 0-4-4H8a4 4 0 0 0-4 4v2" /><circle cx="12" cy="7" r="4" />
                  </svg>
                  Profile
                </router-link>
                <div class="user-dropdown-divider" />
                <div class="user-dropdown-theme">
                  <span class="user-dropdown-theme-label">Theme</span>
                  <div class="theme-switcher">
                    <button
                      class="theme-btn"
                      :class="{ active: currentTheme === 'light' }"
                      @click="setTheme('light')"
                    >Light</button>
                    <button
                      class="theme-btn"
                      :class="{ active: currentTheme === 'dark' }"
                      @click="setTheme('dark')"
                    >Dark</button>
                    <button
                      class="theme-btn"
                      :class="{ active: currentTheme === 'system' }"
                      @click="setTheme('system')"
                    >System</button>
                  </div>
                </div>
                <div class="user-dropdown-divider" />
                <button class="user-dropdown-item user-dropdown-logout" @click="handleLogout">
                  <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                    <path d="M9 21H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h4" /><polyline points="16 17 21 12 16 7" /><line x1="21" y1="12" x2="9" y2="12" />
                  </svg>
                  Logout
                </button>
              </div>
            </Transition>
          </div>
        </div>
      </header>

      <!-- Page content -->
      <main class="main-content">
        <router-view />
      </main>

      <!-- Footer -->
      <footer class="main-footer">
        <span>&copy; {{ new Date().getFullYear() }}
          <a href="https://github.com/jkaninda" target="_blank" rel="noopener noreferrer">Jonas Kaninda</a>.
        </span>
        <a href="https://jkaninda.dev" target="_blank" rel="noopener noreferrer" class="footer-github">
          <svg width="16" height="16" viewBox="0 0 24 24" fill="currentColor">
            <path d="M12 0C5.37 0 0 5.37 0 12c0 5.31 3.435 9.795 8.205 11.385.6.105.825-.255.825-.57 0-.285-.015-1.23-.015-2.235-3.015.555-3.795-.735-4.035-1.41-.135-.345-.72-1.41-1.23-1.695-.42-.225-1.02-.78-.015-.795.945-.015 1.62.87 1.845 1.23 1.08 1.815 2.805 1.305 3.495.99.105-.78.42-1.305.765-1.605-2.67-.3-5.46-1.335-5.46-5.925 0-1.305.465-2.385 1.23-3.225-.12-.3-.54-1.53.12-3.18 0 0 1.005-.315 3.3 1.23.96-.27 1.98-.405 3-.405s2.04.135 3 .405c2.295-1.56 3.3-1.23 3.3-1.23.66 1.65.24 2.88.12 3.18.765.84 1.23 1.905 1.23 3.225 0 4.605-2.805 5.625-5.475 5.925.435.375.81 1.095.81 2.22 0 1.605-.015 2.895-.015 3.3 0 .315.225.69.825.57A12.02 12.02 0 0 0 24 12c0-6.63-5.37-12-12-12z" />
          </svg>
          GitHub
        </a>
      </footer>
    </div>

    <!-- Mobile sidebar overlay -->
    <Transition name="overlay-fade">
      <div v-if="mobileOpen" class="sidebar-overlay" @click="mobileOpen = false" />
    </Transition>

    <!-- Mobile sidebar -->
    <Transition name="sidebar-slide">
      <aside v-if="mobileOpen" class="sidebar sidebar-mobile">
        <div class="sidebar-header">
          <img src="/logo.png" alt="Goma" class="sidebar-logo" />
          <span class="sidebar-brand-text">Goma Admin</span>
          <button class="sidebar-collapse-btn" @click="mobileOpen = false">
            <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
              <line x1="18" y1="6" x2="6" y2="18" /><line x1="6" y1="6" x2="18" y2="18" />
            </svg>
          </button>
        </div>

        <!-- Instance Switcher (mobile) -->
        <div v-if="instanceStore.instances.length > 0" class="instance-switcher">
          <div class="instance-switcher-toggle" @click="instanceDropdownOpen = !instanceDropdownOpen">
            <div style="display: flex; align-items: center; gap: 8px; overflow: hidden;">
              <span class="instance-avatar">{{ activeInstanceInitial }}</span>
              <span class="instance-name">{{ activeInstanceName }}</span>
            </div>
            <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
              <polyline points="6 9 12 15 18 9" />
            </svg>
          </div>
          <Transition name="dropdown">
            <div v-if="instanceDropdownOpen" class="instance-dropdown">
              <div
                class="instance-option"
                :class="{ active: !instanceStore.currentInstanceId }"
                @click="switchInstance(null); mobileOpen = false"
              >
                <span class="instance-avatar">A</span>
                <span>All Instances</span>
              </div>
              <div
                v-for="inst in instanceStore.instances"
                :key="inst.id"
                class="instance-option"
                :class="{ active: instanceStore.currentInstanceId === inst.id }"
                @click="switchInstance(inst.id); mobileOpen = false"
              >
                <span class="instance-avatar">{{ inst.name?.charAt(0)?.toUpperCase() || '?' }}</span>
                <span>{{ inst.name }}</span>
              </div>
            </div>
          </Transition>
        </div>

        <nav class="sidebar-nav">
          <div class="nav-section">
            <router-link to="/" class="nav-item" :class="{ active: $route.path === '/' }" @click="mobileOpen = false">
              <svg class="nav-icon" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                <rect x="3" y="3" width="7" height="7" /><rect x="14" y="3" width="7" height="7" /><rect x="3" y="14" width="7" height="7" /><rect x="14" y="14" width="7" height="7" />
              </svg>
              <span class="nav-label">Dashboard</span>
            </router-link>
            <router-link to="/routes" class="nav-item" :class="{ active: $route.path.startsWith('/routes') }" @click="mobileOpen = false">
              <svg class="nav-icon" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                <circle cx="6" cy="6" r="3" /><circle cx="18" cy="18" r="3" /><path d="M6 9a9 9 0 0 0 9 9" />
              </svg>
              <span class="nav-label">Routes</span>
            </router-link>
            <router-link to="/middlewares" class="nav-item" :class="{ active: $route.path.startsWith('/middlewares') }" @click="mobileOpen = false">
              <svg class="nav-icon" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                <line x1="8" y1="6" x2="21" y2="6" /><line x1="8" y1="12" x2="21" y2="12" /><line x1="8" y1="18" x2="21" y2="18" /><line x1="3" y1="6" x2="3.01" y2="6" /><line x1="3" y1="12" x2="3.01" y2="12" /><line x1="3" y1="18" x2="3.01" y2="18" />
              </svg>
              <span class="nav-label">Middlewares</span>
            </router-link>
            <router-link to="/audit" class="nav-item" :class="{ active: $route.path.startsWith('/audit') }" @click="mobileOpen = false">
              <svg class="nav-icon" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                <path d="M9 5H7a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2m-3 7h3m-3 4h3m-6-4h.01M9 16h.01" />
              </svg>
              <span class="nav-label">Audit Log</span>
            </router-link>
            <router-link to="/api-keys" class="nav-item" :class="{ active: $route.path.startsWith('/api-keys') }" @click="mobileOpen = false">
              <svg class="nav-icon" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                <path d="M21 2l-2 2m-7.61 7.61a5.5 5.5 0 1 1-7.778 7.778 5.5 5.5 0 0 1 7.777-7.777zm0 0L15.5 7.5m0 0l3 3L22 7l-3-3m-3.5 3.5L19 4" />
              </svg>
              <span class="nav-label">API Keys</span>
            </router-link>
            <router-link to="/instances" class="nav-item" :class="{ active: $route.path.startsWith('/instances') }" @click="mobileOpen = false">
              <svg class="nav-icon" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                <rect x="2" y="2" width="20" height="8" rx="2" /><rect x="2" y="14" width="20" height="8" rx="2" /><circle cx="6" cy="6" r="1" /><circle cx="6" cy="18" r="1" />
              </svg>
              <span class="nav-label">Instances</span>
            </router-link>
            <router-link to="/settings" class="nav-item" :class="{ active: $route.path.startsWith('/settings') }" @click="mobileOpen = false">
              <svg class="nav-icon" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                <circle cx="12" cy="12" r="3" /><path d="M19.4 15a1.65 1.65 0 0 0 .33 1.82l.06.06a2 2 0 1 1-2.83 2.83l-.06-.06a1.65 1.65 0 0 0-1.82-.33 1.65 1.65 0 0 0-1 1.51V21a2 2 0 0 1-4 0v-.09A1.65 1.65 0 0 0 9 19.4a1.65 1.65 0 0 0-1.82.33l-.06.06a2 2 0 1 1-2.83-2.83l.06-.06A1.65 1.65 0 0 0 4.68 15a1.65 1.65 0 0 0-1.51-1H3a2 2 0 0 1 0-4h.09A1.65 1.65 0 0 0 4.6 9a1.65 1.65 0 0 0-.33-1.82l-.06-.06a2 2 0 1 1 2.83-2.83l.06.06A1.65 1.65 0 0 0 9 4.68a1.65 1.65 0 0 0 1-1.51V3a2 2 0 0 1 4 0v.09a1.65 1.65 0 0 0 1 1.51 1.65 1.65 0 0 0 1.82-.33l.06-.06a2 2 0 1 1 2.83 2.83l-.06.06A1.65 1.65 0 0 0 19.4 9a1.65 1.65 0 0 0 1.51 1H21a2 2 0 0 1 0 4h-.09a1.65 1.65 0 0 0-1.51 1z" />
              </svg>
              <span class="nav-label">Settings</span>
            </router-link>
            <router-link to="/profile" class="nav-item" :class="{ active: $route.path === '/profile' }" @click="mobileOpen = false">
              <svg class="nav-icon" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                <path d="M20 21v-2a4 4 0 0 0-4-4H8a4 4 0 0 0-4 4v2" /><circle cx="12" cy="7" r="4" />
              </svg>
              <span class="nav-label">Profile</span>
            </router-link>
          </div>
          <!-- Admin section -->
          <div v-if="authStore.isAdmin" class="nav-section">
            <div class="nav-section-title">Admin</div>
            <router-link to="/users" class="nav-item" :class="{ active: $route.path.startsWith('/users') }" @click="mobileOpen = false">
              <svg class="nav-icon" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                <path d="M17 21v-2a4 4 0 0 0-4-4H5a4 4 0 0 0-4 4v2" /><circle cx="9" cy="7" r="4" /><path d="M23 21v-2a4 4 0 0 0-3-3.87" /><path d="M16 3.13a4 4 0 0 1 0 7.75" />
              </svg>
              <span class="nav-label">Users</span>
            </router-link>
            <router-link to="/oauth-settings" class="nav-item" :class="{ active: $route.path.startsWith('/oauth-settings') }" @click="mobileOpen = false">
              <svg class="nav-icon" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                <path d="M15 7a2 2 0 0 1 2 2m4 0a6 6 0 0 1-7.743 5.743L11 17H9v2H7v2H4a1 1 0 0 1-1-1v-2.586a1 1 0 0 1 .293-.707l5.964-5.964A6 6 0 1 1 21 9z" />
              </svg>
              <span class="nav-label">OAuth</span>
            </router-link>
          </div>
        </nav>
      </aside>
    </Transition>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { useInstanceStore } from '@/stores/instance'

const router = useRouter()
const authStore = useAuthStore()
const instanceStore = useInstanceStore()

// Sidebar state
const sidebarCollapsed = ref(localStorage.getItem('goma_sidebar_collapsed') === 'true')
const mobileOpen = ref(false)
const userMenuOpen = ref(false)
const userMenuRef = ref<HTMLElement | null>(null)
const instanceDropdownOpen = ref(false)

// Theme state
const currentTheme = ref(localStorage.getItem('goma_theme') || 'system')

function toggleSidebar() {
  sidebarCollapsed.value = !sidebarCollapsed.value
  localStorage.setItem('goma_sidebar_collapsed', String(sidebarCollapsed.value))
}

const userInitial = computed(() => {
  const name = authStore.user?.name || authStore.user?.email || '?'
  return name.charAt(0).toUpperCase()
})

const activeInstanceName = computed(() => {
  if (!instanceStore.currentInstanceId) return 'All Instances'
  const inst = instanceStore.instances.find((i: any) => i.id === instanceStore.currentInstanceId)
  return inst?.name || 'All Instances'
})

const activeInstanceInitial = computed(() => {
  if (!instanceStore.currentInstanceId) return 'A'
  const inst = instanceStore.instances.find((i: any) => i.id === instanceStore.currentInstanceId)
  return inst?.name?.charAt(0)?.toUpperCase() || 'A'
})

function switchInstance(id: number | null) {
  instanceStore.setInstance(id)
  instanceDropdownOpen.value = false
  // Force reload current route so views refetch data for the new instance
  router.go(0)
}

function setTheme(theme: string) {
  currentTheme.value = theme
  localStorage.setItem('goma_theme', theme)
  applyTheme(theme)
}

function applyTheme(theme: string) {
  const root = document.documentElement
  if (theme === 'system') {
    const prefersDark = window.matchMedia('(prefers-color-scheme: dark)').matches
    root.setAttribute('data-theme', prefersDark ? 'dark' : 'light')
  } else {
    root.setAttribute('data-theme', theme)
  }
}

function handleLogout() {
  userMenuOpen.value = false
  authStore.logout()
}

function handleClickOutside(e: MouseEvent) {
  if (userMenuRef.value && !userMenuRef.value.contains(e.target as Node)) {
    userMenuOpen.value = false
  }
}

function handleSystemThemeChange() {
  if (currentTheme.value === 'system') {
    applyTheme('system')
  }
}

onMounted(() => {
  applyTheme(currentTheme.value)
  document.addEventListener('click', handleClickOutside)
  window.matchMedia('(prefers-color-scheme: dark)').addEventListener('change', handleSystemThemeChange)
  instanceStore.fetchInstances()
})

onUnmounted(() => {
  document.removeEventListener('click', handleClickOutside)
  window.matchMedia('(prefers-color-scheme: dark)').removeEventListener('change', handleSystemThemeChange)
})
</script>

<style scoped>
.dashboard-layout {
  display: flex;
  min-height: 100vh;
  background: var(--bg-secondary);
}

/* ─── Sidebar ─── */
.sidebar {
  position: fixed;
  top: 0;
  left: 0;
  bottom: 0;
  width: 240px;
  background: var(--bg-sidebar);
  display: flex;
  flex-direction: column;
  z-index: 40;
  transition: width var(--transition-slow);
  overflow: hidden;
}

.sidebar-collapsed .sidebar:not(.sidebar-mobile) {
  width: 64px;
}

.sidebar-header {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 16px 14px 12px;
  flex-shrink: 0;
  position: relative;
}

.sidebar-logo {
  width: 28px;
  height: 28px;
  flex-shrink: 0;
  border-radius: 6px;
}

.sidebar-brand-text {
  font-size: 15px;
  font-weight: 700;
  color: #ffffff;
  white-space: nowrap;
  overflow: hidden;
  opacity: 1;
  transition: opacity var(--transition-slow);
}

.sidebar-collapsed .sidebar:not(.sidebar-mobile) .sidebar-brand-text {
  opacity: 0;
  width: 0;
}

.sidebar-collapse-btn {
  background: none;
  border: none;
  color: var(--sidebar-text);
  cursor: pointer;
  padding: 4px;
  border-radius: var(--radius-sm);
  display: flex;
  align-items: center;
  justify-content: center;
  transition: color var(--transition), background var(--transition);
  flex-shrink: 0;
  position: absolute;
  top: 18px;
  right: 10px;
}

.sidebar-collapse-btn:hover {
  color: #ffffff;
  background: var(--sidebar-hover);
}

.sidebar-collapsed .sidebar:not(.sidebar-mobile) .sidebar-collapse-btn {
  /* left: 50%; */
  transform: translateX(-50%);
  right: -15px;
  top: 20px;
  z-index: 999;
}

/* ─── Navigation ─── */
.sidebar-nav {
  flex: 1;
  overflow-y: auto;
  overflow-x: hidden;
  padding: 8px;
}

.nav-section {
  margin-bottom: 8px;
}

.nav-section-title {
  font-size: 11px;
  font-weight: 600;
  color: var(--sidebar-text);
  text-transform: uppercase;
  letter-spacing: 0.06em;
  padding: 12px 12px 6px;
  white-space: nowrap;
  overflow: hidden;
  opacity: 1;
  transition: opacity var(--transition-slow);
}

.sidebar-collapsed .sidebar:not(.sidebar-mobile) .nav-section-title {
  opacity: 0;
  height: 0;
  padding: 0;
  margin: 0;
}

.nav-item {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 9px 12px;
  border-radius: var(--radius);
  color: var(--sidebar-text);
  font-size: 14px;
  font-weight: 500;
  cursor: pointer;
  transition: background var(--transition), color var(--transition);
  text-decoration: none;
  white-space: nowrap;
  overflow: hidden;
}

.nav-item:hover {
  background: var(--sidebar-hover);
  color: var(--sidebar-text-active);
}

.nav-item.active {
  background: var(--sidebar-hover);
  color: var(--sidebar-text-active);
}

.nav-item.active::before {
  content: '';
  position: absolute;
  left: 0;
  top: 50%;
  transform: translateY(-50%);
  width: 3px;
  height: 20px;
  background: var(--primary-500);
  border-radius: 0 3px 3px 0;
}

.nav-icon {
  flex-shrink: 0;
  width: 18px;
  height: 18px;
}

.nav-label {
  overflow: hidden;
  opacity: 1;
  transition: opacity var(--transition-slow);
}

.sidebar-collapsed .sidebar:not(.sidebar-mobile) .nav-label {
  opacity: 0;
  width: 0;
}

.sidebar-collapsed .sidebar:not(.sidebar-mobile) .nav-item {
  justify-content: center;
  padding: 9px;
}

/* ─── Sidebar footer ─── */
.sidebar-footer {
  border-top: 1px solid var(--sidebar-border);
  padding: 8px;
  flex-shrink: 0;
}

.sidebar-version {
  padding: 8px 12px;
  font-size: 12px;
  color: var(--sidebar-text);
  opacity: 0.6;
  white-space: nowrap;
  overflow: hidden;
}

.sidebar-collapsed .sidebar:not(.sidebar-mobile) .sidebar-version {
  opacity: 0;
}

.sidebar-collapsed .sidebar:not(.sidebar-mobile) .sidebar-footer .nav-item {
  justify-content: center;
  padding: 9px;
}

/* ─── Main wrapper ─── */
.main-wrapper {
  flex: 1;
  margin-left: 240px;
  display: flex;
  flex-direction: column;
  min-height: 100vh;
  transition: margin-left var(--transition-slow);
}

.sidebar-collapsed .main-wrapper {
  margin-left: 64px;
}

/* ─── Top bar ─── */
.topbar {
  position: sticky;
  top: 0;
  z-index: 30;
  display: flex;
  align-items: center;
  justify-content: space-between;
  height: 56px;
  padding: 0 24px;
  background: var(--bg-primary);
  border-bottom: 1px solid var(--border-primary);
  transition: background var(--transition-slow), border-color var(--transition-slow);
}

.topbar-left {
  display: flex;
  align-items: center;
}

.mobile-menu-btn {
  display: none;
  align-items: center;
  justify-content: center;
  background: none;
  border: none;
  color: var(--text-tertiary);
  cursor: pointer;
  padding: 6px;
  border-radius: var(--radius-sm);
  transition: color var(--transition), background var(--transition);
}

.mobile-menu-btn:hover {
  color: var(--text-primary);
  background: var(--bg-hover);
}

.topbar-right {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-left: auto;
}

/* ─── User menu ─── */
.user-menu {
  position: relative;
}

.user-menu-trigger {
  display: flex;
  align-items: center;
  gap: 8px;
  background: none;
  border: 1px solid transparent;
  padding: 5px 10px 5px 5px;
  border-radius: var(--radius);
  cursor: pointer;
  color: var(--text-secondary);
  font-family: inherit;
  font-size: 14px;
  transition: background var(--transition), border-color var(--transition);
}

.user-menu-trigger:hover {
  background: var(--bg-hover);
  border-color: var(--border-primary);
}

.user-avatar {
  width: 30px;
  height: 30px;
  border-radius: 50%;
  background: var(--primary-600);
  color: var(--text-on-primary);
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 13px;
  font-weight: 600;
  flex-shrink: 0;
}

.user-avatar-lg {
  width: 36px;
  height: 36px;
  font-size: 15px;
}

.user-name {
  font-weight: 500;
  max-width: 140px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

/* ─── User dropdown ─── */
.user-dropdown {
  position: absolute;
  top: calc(100% + 6px);
  right: 0;
  width: 280px;
  background: var(--bg-primary);
  border: 1px solid var(--border-primary);
  border-radius: var(--radius-lg);
  box-shadow: var(--shadow-lg);
  z-index: 50;
  overflow: hidden;
}

.user-dropdown-header {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 14px 16px;
}

.user-dropdown-info {
  overflow: hidden;
}

.user-dropdown-name {
  font-size: 14px;
  font-weight: 600;
  color: var(--text-primary);
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.user-dropdown-email {
  font-size: 12px;
  color: var(--text-muted);
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.user-dropdown-divider {
  height: 1px;
  background: var(--border-primary);
}

.user-dropdown-item {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 10px 16px;
  font-size: 14px;
  color: var(--text-secondary);
  cursor: pointer;
  transition: background var(--transition), color var(--transition);
  text-decoration: none;
  border: none;
  background: none;
  width: 100%;
  font-family: inherit;
}

.user-dropdown-item:hover {
  background: var(--bg-hover);
  color: var(--text-primary);
}

.user-dropdown-logout {
  color: var(--danger-600);
}

.user-dropdown-logout:hover {
  background: var(--danger-50);
  color: var(--danger-700);
}

/* ─── Theme switcher ─── */
.user-dropdown-theme {
  padding: 10px 16px;
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.user-dropdown-theme-label {
  font-size: 13px;
  color: var(--text-muted);
  font-weight: 500;
}

.theme-switcher {
  display: flex;
  background: var(--bg-tertiary);
  border-radius: var(--radius-sm);
  padding: 2px;
  gap: 2px;
}

.theme-btn {
  padding: 4px 10px;
  border: none;
  border-radius: 4px;
  font-size: 12px;
  font-weight: 500;
  cursor: pointer;
  color: var(--text-tertiary);
  background: transparent;
  font-family: inherit;
  transition: all var(--transition);
}

.theme-btn:hover {
  color: var(--text-primary);
}

.theme-btn.active {
  background: var(--bg-primary);
  color: var(--text-primary);
  box-shadow: var(--shadow-sm);
}

/* ─── Main content ─── */
.main-content {
  flex: 1;
  padding: 28px;
}

/* ─── Footer ─── */
.main-footer {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 16px 28px;
  border-top: 1px solid var(--border-primary);
  font-size: 13px;
  color: var(--text-muted);
  background: var(--bg-primary);
  transition: background var(--transition-slow), border-color var(--transition-slow);
}

.main-footer a {
  color: var(--primary-600);
  text-decoration: none;
}

.main-footer a:hover {
  color: var(--primary-700);
}

.footer-github {
  display: flex;
  align-items: center;
  gap: 6px;
  color: var(--text-muted);
  transition: color var(--transition);
}

.footer-github:hover {
  color: var(--text-primary);
}

/* ─── Mobile overlay ─── */
.sidebar-overlay {
  position: fixed;
  inset: 0;
  background: var(--overlay);
  z-index: 35;
  backdrop-filter: blur(4px);
}

.sidebar-mobile {
  z-index: 45;
  width: 240px;
}

/* ─── Dropdown transitions ─── */
.dropdown-enter-active {
  transition: opacity 150ms ease, transform 150ms ease;
}
.dropdown-leave-active {
  transition: opacity 100ms ease, transform 100ms ease;
}
.dropdown-enter-from,
.dropdown-leave-to {
  opacity: 0;
  transform: translateY(-4px);
}

.overlay-fade-enter-active {
  transition: opacity 200ms ease;
}
.overlay-fade-leave-active {
  transition: opacity 150ms ease;
}
.overlay-fade-enter-from,
.overlay-fade-leave-to {
  opacity: 0;
}

.sidebar-slide-enter-active {
  transition: transform 200ms ease;
}
.sidebar-slide-leave-active {
  transition: transform 150ms ease;
}
.sidebar-slide-enter-from,
.sidebar-slide-leave-to {
  transform: translateX(-100%);
}

/* ─── Responsive ─── */
@media (max-width: 1024px) {
  .sidebar:not(.sidebar-mobile) {
    display: none;
  }

  .main-wrapper {
    margin-left: 0 !important;
  }

  .mobile-menu-btn {
    display: flex;
  }
}

@media (max-width: 640px) {
  .main-content {
    padding: 20px 16px;
  }

  .topbar {
    padding: 0 16px;
  }

  .main-footer {
    padding: 14px 16px;
    flex-direction: column;
    gap: 8px;
    text-align: center;
  }

  .user-name {
    display: none;
  }
}

/* ─── Instance Switcher ─── */
.instance-switcher {
  padding: 0 8px 8px;
  position: relative;
}

.instance-switcher-toggle {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 8px 10px;
  border-radius: var(--radius);
  cursor: pointer;
  color: var(--sidebar-text);
  border: 1px solid var(--sidebar-border);
  transition: background var(--transition), color var(--transition);
}

.instance-switcher-toggle:hover {
  background: var(--sidebar-hover);
  color: var(--sidebar-text-active);
}

.instance-avatar {
  width: 24px;
  height: 24px;
  border-radius: 6px;
  background: var(--primary-600);
  color: #fff;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 12px;
  font-weight: 600;
  flex-shrink: 0;
}

.instance-name {
  font-size: 13px;
  font-weight: 500;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.instance-dropdown {
  position: absolute;
  top: 100%;
  left: 8px;
  right: 8px;
  background: var(--bg-primary);
  border: 1px solid var(--border-primary);
  border-radius: var(--radius);
  box-shadow: var(--shadow-lg);
  padding: 4px;
  z-index: 200;
  margin-top: 4px;
}

.instance-option {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 8px 10px;
  border-radius: var(--radius-sm);
  cursor: pointer;
  font-size: 13px;
  color: var(--text-secondary);
  transition: background var(--transition), color var(--transition);
}

.instance-option:hover {
  background: var(--bg-hover);
  color: var(--text-primary);
}

.instance-option.active {
  background: var(--primary-50, rgba(99, 102, 241, 0.1));
  color: var(--primary-600);
  font-weight: 500;
}

.sidebar-collapsed .sidebar:not(.sidebar-mobile) .instance-switcher-toggle {
  justify-content: center;
  padding: 8px;
}

.sidebar-collapsed .sidebar:not(.sidebar-mobile) .instance-switcher-toggle .instance-name,
.sidebar-collapsed .sidebar:not(.sidebar-mobile) .instance-switcher-toggle svg.nav-label {
  opacity: 0;
  width: 0;
  overflow: hidden;
}
</style>
