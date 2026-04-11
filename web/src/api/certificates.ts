import api from './client'

export interface CertificateSummary {
  domain: string
  domains: string[]
  issuer: string
  source: 'acme' | 'custom'
  expiresAt: string
  issuedAt: string
  isExpired: boolean
  daysUntilExpiry: number
  serialNumber: string
  fingerprint: string
}

export interface TLSOverview {
  totalCertificates: number
  acmeCertificates: number
  customCertificates: number
  expiredCount: number
  expiringSoonCount: number
}

export const certificatesApi = {
  list(page = 0, size = 20, search = '') {
    const params: Record<string, unknown> = { page, size }
    if (search) params.q = search
    return api.get('/tls/certificates', { params })
  },
  get(domain: string) {
    return api.get<CertificateSummary>(`/tls/certificates/${encodeURIComponent(domain)}`)
  },
  overview() {
    return api.get<TLSOverview>('/tls/overview')
  },
}
