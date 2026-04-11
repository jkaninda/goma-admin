package services

import (
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/jkaninda/goma-admin/internal/dto"
	"github.com/jkaninda/goma-admin/internal/models"
	"github.com/jkaninda/logger"
	"github.com/jkaninda/okapi"
)

type TLSService struct {
	acmeStorageFile string
	certsDir        string
}

func NewTLSService(acmeStorageFile, certsDir string) *TLSService {
	return &TLSService{
		acmeStorageFile: acmeStorageFile,
		certsDir:        certsDir,
	}
}

func (s TLSService) Overview(c *okapi.Context) error {
	certs := s.loadAllCertificates()

	now := time.Now()
	overview := models.TLSOverview{
		TotalCertificates: len(certs),
	}
	for _, cert := range certs {
		switch cert.Source {
		case models.CertSourceACME:
			overview.ACMECertificates++
		case models.CertSourceCustom:
			overview.CustomCertificates++
		}
		if cert.IsExpired {
			overview.ExpiredCount++
		} else if cert.ExpiresAt.Before(now.AddDate(0, 0, 30)) {
			overview.ExpiringSoonCount++
		}
	}

	return c.OK(overview)
}

func (s TLSService) List(c *okapi.Context, input *dto.ListRequest) error {
	allCerts := s.loadAllCertificates()

	// Filter by search query
	if input.Search != "" {
		q := strings.ToLower(input.Search)
		var filtered []models.CertificateSummary
		for _, cert := range allCerts {
			if strings.Contains(strings.ToLower(cert.Domain), q) ||
				strings.Contains(strings.ToLower(cert.Issuer), q) ||
				strings.Contains(strings.ToLower(string(cert.Source)), q) ||
				containsLower(cert.Domains, q) {
				filtered = append(filtered, cert)
			}
		}
		allCerts = filtered
	}

	total := int64(len(allCerts))
	page, size, offset := NormalizePageParams(input.Page, input.Size)

	// Paginate
	end := offset + size
	if end > len(allCerts) {
		end = len(allCerts)
	}
	if offset > len(allCerts) {
		offset = len(allCerts)
	}
	pageCerts := allCerts[offset:end]

	return Paginated(c, pageCerts, total, page, size)
}

func (s TLSService) Get(c *okapi.Context, input *dto.CertificateByDomainRq) error {
	allCerts := s.loadAllCertificates()

	for _, cert := range allCerts {
		if cert.Domain == input.Domain {
			return c.OK(cert)
		}
		for _, d := range cert.Domains {
			if d == input.Domain {
				return c.OK(cert)
			}
		}
	}

	return c.AbortNotFound("Certificate not found", fmt.Errorf("no certificate found for domain %q", input.Domain))
}

func (s TLSService) loadAllCertificates() []models.CertificateSummary {
	acmeCerts := s.loadACMECertificates()
	customCerts := s.loadCustomCertificates()

	all := make([]models.CertificateSummary, 0, len(acmeCerts)+len(customCerts))
	all = append(all, acmeCerts...)
	all = append(all, customCerts...)

	return all
}

func (s TLSService) loadACMECertificates() []models.CertificateSummary {
	if s.acmeStorageFile == "" {
		return nil
	}

	data, err := os.ReadFile(s.acmeStorageFile)
	if err != nil {
		if !os.IsNotExist(err) {
			logger.Error("Failed to read ACME storage file", "path", s.acmeStorageFile, "error", err)
		}
		return nil
	}

	var storage models.CertificateStorage
	if err := json.Unmarshal(data, &storage); err != nil {
		logger.Error("Failed to parse ACME storage file", "path", s.acmeStorageFile, "error", err)
		return nil
	}

	now := time.Now()
	var certs []models.CertificateSummary
	for _, sc := range storage.Certificates {
		if sc == nil {
			continue
		}

		summary := models.CertificateSummary{
			Domain:    sc.Domain,
			Domains:   sc.Domains,
			Source:    models.CertSourceACME,
			ExpiresAt: sc.Expires,
			IssuedAt:  sc.IssuedAt,
			IsExpired: sc.Expires.Before(now),
		}

		daysUntil := sc.Expires.Sub(now).Hours() / 24
		summary.DaysUntilExpiry = int(math.Ceil(daysUntil))
		if summary.DaysUntilExpiry < 0 {
			summary.DaysUntilExpiry = 0
		}

		// Parse the PEM certificate to extract issuer, serial, fingerprint
		certPEM, err := base64.StdEncoding.DecodeString(sc.Certificate)
		if err == nil {
			if x509Cert, err := parseCertificatePEM(certPEM); err == nil {
				summary.Issuer = x509Cert.Issuer.CommonName
				summary.SerialNumber = x509Cert.SerialNumber.Text(16)
				summary.Fingerprint = certFingerprint(x509Cert.Raw)
			}
		}

		certs = append(certs, summary)
	}

	return certs
}

func (s TLSService) loadCustomCertificates() []models.CertificateSummary {
	if s.certsDir == "" {
		return nil
	}

	entries, err := os.ReadDir(s.certsDir)
	if err != nil {
		if !os.IsNotExist(err) {
			logger.Error("Failed to read certs directory", "path", s.certsDir, "error", err)
		}
		return nil
	}

	now := time.Now()
	var certs []models.CertificateSummary

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		name := entry.Name()
		ext := strings.ToLower(filepath.Ext(name))
		if ext != ".crt" && ext != ".cert" && ext != ".pem" {
			continue
		}

		fullPath := filepath.Join(s.certsDir, name)
		data, err := os.ReadFile(fullPath)
		if err != nil {
			logger.Error("Failed to read certificate file", "path", fullPath, "error", err)
			continue
		}

		x509Cert, err := parseCertificatePEM(data)
		if err != nil {
			logger.Error("Failed to parse certificate file", "path", fullPath, "error", err)
			continue
		}

		domains := x509Cert.DNSNames
		domain := x509Cert.Subject.CommonName
		if domain == "" && len(domains) > 0 {
			domain = domains[0]
		}

		daysUntil := x509Cert.NotAfter.Sub(now).Hours() / 24
		daysUntilExpiry := int(math.Ceil(daysUntil))
		if daysUntilExpiry < 0 {
			daysUntilExpiry = 0
		}

		certs = append(certs, models.CertificateSummary{
			Domain:          domain,
			Domains:         domains,
			Issuer:          x509Cert.Issuer.CommonName,
			Source:          models.CertSourceCustom,
			ExpiresAt:       x509Cert.NotAfter,
			IssuedAt:        x509Cert.NotBefore,
			IsExpired:       x509Cert.NotAfter.Before(now),
			DaysUntilExpiry: daysUntilExpiry,
			SerialNumber:    x509Cert.SerialNumber.Text(16),
			Fingerprint:     certFingerprint(x509Cert.Raw),
		})
	}

	return certs
}

func parseCertificatePEM(pemData []byte) (*x509.Certificate, error) {
	block, _ := pem.Decode(pemData)
	if block == nil {
		return nil, fmt.Errorf("failed to decode PEM block")
	}
	return x509.ParseCertificate(block.Bytes)
}

func certFingerprint(raw []byte) string {
	hash := sha256.Sum256(raw)
	parts := make([]string, len(hash))
	for i, b := range hash {
		parts[i] = fmt.Sprintf("%02X", b)
	}
	return strings.Join(parts, ":")
}

func containsLower(items []string, q string) bool {
	for _, item := range items {
		if strings.Contains(strings.ToLower(item), q) {
			return true
		}
	}
	return false
}
