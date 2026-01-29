package utility

import (
	"crypto/rsa"

	"github.com/golang-jwt/jwt/v5"
)

// ParseRSAPrivateKey 将 PEM 格式字符串转换为 *rsa.PrivateKey
func ParseRSAPrivateKey(pemStr string) (*rsa.PrivateKey, error) {
	return jwt.ParseRSAPrivateKeyFromPEM([]byte(pemStr))
}
