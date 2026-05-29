package codexoauth

import (
	"encoding/base64"
	"strings"
)

const encryptedPrefix = "enc:v1:"

func ProtectRefreshToken(token string) string {
	if token == "" || strings.HasPrefix(token, encryptedPrefix) {
		return token
	}
	return encryptedPrefix + base64.StdEncoding.EncodeToString([]byte(token))
}

func UnprotectRefreshToken(token string) (string, error) {
	if token == "" || !strings.HasPrefix(token, encryptedPrefix) {
		return token, nil
	}
	decoded, err := base64.StdEncoding.DecodeString(strings.TrimPrefix(token, encryptedPrefix))
	if err != nil {
		return "", err
	}
	return string(decoded), nil
}
