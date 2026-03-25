import api from './client'

export interface DashboardStats {
  users: number
  instances: number
  middlewares: number
  routes: number
}

export const dashboardApi = {
  getStats() {
    return api.get<DashboardStats>('/dashboard')
  },
}
