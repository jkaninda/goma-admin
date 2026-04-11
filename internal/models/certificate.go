package models

import "time"

type CertificateSource string

const (
	CertSourceACME   CertificateSource = "acme"
	CertSourceCustom CertificateSource = "custom"
)

type CertificateSummary struct {
	Domain          string            `json:"domain"`
	Domains         []string          `json:"domains"`
	Issuer          string            `json:"issuer"`
	Source          CertificateSource `json:"source"`
	ExpiresAt       time.Time         `json:"expiresAt"`
	IssuedAt        time.Time         `json:"issuedAt"`
	IsExpired       bool              `json:"isExpired"`
	DaysUntilExpiry int               `json:"daysUntilExpiry"`
	SerialNumber    string            `json:"serialNumber"`
	Fingerprint     string            `json:"fingerprint"`
}

type TLSOverview struct {
	TotalCertificates  int `json:"totalCertificates"`
	ACMECertificates   int `json:"acmeCertificates"`
	CustomCertificates int `json:"customCertificates"`
	ExpiredCount       int `json:"expiredCount"`
	ExpiringSoonCount  int `json:"expiringSoonCount"`
}

type CertificateStorage struct {
	UserAccount  *StoredUserAccount   `json:"user_account"`
	Certificates []*StoredCertificate `json:"certificates"`
	Version      string               `json:"version"`
	UpdatedAt    time.Time            `json:"updated_at"`
}

type StoredUserAccount struct {
	Email        string `json:"email"`
	PrivateKey   string `json:"private_key"`
	Registration string `json:"registration"`
}

// StoredCertificate holds a stored certificate from acme.json.
type StoredCertificate struct {
	Domain      string    `json:"domain"`
	Certificate string    `json:"certificate"`
	PrivateKey  string    `json:"private_key"`
	Domains     []string  `json:"domains"`
	Expires     time.Time `json:"expires"`
	IssuedAt    time.Time `json:"issued_at"`
}
