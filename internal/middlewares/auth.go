package middlewares

import (
	"github.com/jkaninda/goma-admin/internal/config"
	"github.com/jkaninda/okapi"
)

type Auth struct {
	JWT     *okapi.JWTAuth
	SSEAuth *okapi.JWTAuth
}

func NewAuth(conf *config.Config) *Auth {
	jwtAuth := &okapi.JWTAuth{
		SigningSecret: []byte(conf.JWT.Secret),
		TokenLookup:   "header:Authorization",
		Issuer:        conf.JWT.Issuer,
		Audience:      "goma-admin",
		ForwardClaims: map[string]string{
			"user_id": "sub",
			"email":   "email",
			"role":    "role",
		},
	}
	// SSE endpoints use query param auth since EventSource cannot send headers.
	sseAuth := &okapi.JWTAuth{
		SigningSecret: []byte(conf.JWT.Secret),
		TokenLookup:   "query:token",
		Issuer:        conf.JWT.Issuer,
		Audience:      "goma-admin",
		ForwardClaims: map[string]string{
			"user_id": "sub",
			"email":   "email",
			"role":    "role",
		},
	}
	return &Auth{JWT: jwtAuth, SSEAuth: sseAuth}
}
