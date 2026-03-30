package twofactor

import (
	"github.com/pquerna/otp/totp"
)

func GenerateSecret(email string) (string, string, error) {
	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      "Goma Admin",
		AccountName: email,
	})
	if err != nil {
		return "", "", err
	}
	return key.Secret(), key.URL(), nil
}

func ValidateCode(secret, code string) bool {
	return totp.Validate(code, secret)
}
