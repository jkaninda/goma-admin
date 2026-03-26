export interface ConfigEvent {
  type: string
  resource: string
  resourceId?: number
  instanceId?: number
  name?: string
  message: string
  timestamp: string
}

const CONFIG_EVENT_TYPES = [
  'route_created', 'route_updated', 'route_deleted',
  'middleware_created', 'middleware_updated', 'middleware_deleted',
  'instance_created', 'instance_updated', 'instance_deleted',
  'config_imported', 'config_copied',
]

export function connectConfigSSE(onEvent: (evt: ConfigEvent) => void): EventSource {
  const token = localStorage.getItem('access_token')
  const url = `/api/v1/events${token ? `?token=${encodeURIComponent(token)}` : ''}`
  const es = new EventSource(url)

  for (const type of CONFIG_EVENT_TYPES) {
    es.addEventListener(type, (e) => onEvent(JSON.parse((e as MessageEvent).data)))
  }

  return es
}
