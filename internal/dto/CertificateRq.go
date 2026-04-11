package dto

type CertificateByDomainRq struct {
	Domain string `param:"domain" required:"true" description:"Certificate domain name"`
}
