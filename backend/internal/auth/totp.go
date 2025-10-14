package auth

import (
	"github.com/pquerna/otp/totp"
)

func GenerateTOTPSecret(accountName, issuer string) (string, string, error) {
	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      issuer,
		AccountName: accountName,
	})
	if err != nil {
		return "", "", err
	}

	return key.Secret(), key.URL(), nil
}

func ValidateTOTP(secret, code string) bool {
	return totp.Validate(code, secret)
}

func GenerateQRCode(url string) (string, error) {
	return url, nil
}

