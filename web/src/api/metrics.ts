import api from './client'

export interface RouteMetric {
  routeName: string
  totalRequests: number
  errorCount: number
  errorRate: number
  avgLatencyMs: number
}

export interface InstanceMetrics {
  totalRequests: number
  totalErrors: number
  errorRate: number
  avgLatencyMs: number
  realtimeVisitors: number
  routesCount: number
  middlewaresCount: number
  uptimeSeconds: number
  routeMetrics: RouteMetric[]
}

export const metricsApi = {
  getMetrics(instanceId: number) {
    return api.get<InstanceMetrics>(`/instances/${instanceId}/metrics`)
  },
  getRawMetrics(instanceId: number) {
    return api.get<string>(`/instances/${instanceId}/metrics/raw`, {
      responseType: 'text',
    })
  },
}
