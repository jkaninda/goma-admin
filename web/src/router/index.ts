import { createRouter, createWebHistory } from 'vue-router'

const router = createRouter({
  history: createWebHistory(),
  routes: [
    {
      path: '/auth',
      component: () => import('@/layouts/GuestLayout.vue'),
      children: [
        {
          path: 'login',
          name: 'login',
          component: () => import('@/views/auth/LoginView.vue'),
        },
        {
          path: 'oauth/callback',
          name: 'oauth-callback',
          component: () => import('@/views/auth/OAuthCallbackView.vue'),
        },
      ],
    },
    {
      path: '/',
      component: () => import('@/layouts/DashboardLayout.vue'),
      meta: { requiresAuth: true },
      children: [
        {
          path: '',
          name: 'dashboard',
          component: () => import('@/views/dashboard/DashboardView.vue'),
        },
        {
          path: 'instances',
          name: 'instances',
          component: () => import('@/views/instances/InstancesView.vue'),
        },
        {
          path: 'instances/:id',
          name: 'instance-detail',
          component: () => import('@/views/instances/InstanceDetailView.vue'),
          props: true,
        },
        {
          path: 'repositories',
          name: 'repositories',
          component: () => import('@/views/repositories/RepositoriesView.vue'),
        },
        {
          path: 'repositories/:id',
          name: 'repository-detail',
          component: () => import('@/views/repositories/RepositoryDetailView.vue'),
          props: true,
        },
        {
          path: 'routes',
          name: 'routes',
          component: () => import('@/views/routes/RoutesView.vue'),
        },
        {
          path: 'routes/:id',
          name: 'route-detail',
          component: () => import('@/views/routes/RouteDetailView.vue'),
          props: true,
        },
        {
          path: 'middlewares',
          name: 'middlewares',
          component: () => import('@/views/middlewares/MiddlewaresView.vue'),
        },
        {
          path: 'middlewares/:id',
          name: 'middleware-detail',
          component: () => import('@/views/middlewares/MiddlewareDetailView.vue'),
          props: true,
        },
        {
          path: 'certificates',
          name: 'certificates',
          component: () => import('@/views/certificates/CertificatesView.vue'),
        },
        {
          path: 'certificates/:domain',
          name: 'certificate-detail',
          component: () => import('@/views/certificates/CertificateDetailView.vue'),
          props: true,
        },
        {
          path: 'audit',
          name: 'audit',
          component: () => import('@/views/audit/AuditLogView.vue'),
        },
        {
          path: 'api-keys',
          name: 'api-keys',
          component: () => import('@/views/apikeys/ApiKeysView.vue'),
        },
        {
          path: 'profile',
          name: 'profile',
          component: () => import('@/views/profile/ProfileView.vue'),
        },
        {
          path: 'settings',
          name: 'settings',
          component: () => import('@/views/settings/SettingsView.vue'),
        },
        {
          path: 'users',
          name: 'users',
          component: () => import('@/views/users/UsersView.vue'),
        },
        {
          path: 'oauth-settings',
          name: 'oauth-settings',
          component: () => import('@/views/settings/OAuthSettingsView.vue'),
        },
      ],
    },
  ],
})

router.beforeEach((to, _from, next) => {
  const token = localStorage.getItem('access_token')
  if (to.matched.some((r) => r.meta.requiresAuth) && !token) {
    next({ name: 'login' })
  } else if (to.name === 'login' && token) {
    next({ name: 'dashboard' })
  } else {
    next()
  }
})

export default router
