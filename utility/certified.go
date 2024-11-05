package utility

import (
	"backend-blog/internal/logger"
	"backend-blog/result"
	"crypto/rand"
	"crypto/rsa"
	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

var (
	PrivateKey    *rsa.PrivateKey
	JWTMiddleware fiber.Handler
)

func init() {
	rng := rand.Reader
	var err error
	PrivateKey, err = rsa.GenerateKey(rng, 2048)
	if err != nil {
		logger.Error.Fatalf("rsa.GenerateKey: %v", err)
	}
	JWTMiddleware = jwtware.New(jwtware.Config{
		SigningKey: jwtware.SigningKey{
			JWTAlg: jwtware.RS256,
			Key:    PrivateKey.Public(),
		},
		ErrorHandler: JwtError,
	})
}

func TokenAuth() fiber.Handler {
	return func(c *fiber.Ctx) error {
		user := c.Locals("user").(*jwt.Token)
		claims := user.Claims.(jwt.MapClaims)
		name := claims["username"].(string)
		c.Locals("user", name)
		return c.Next()
	}
}
func JwtError(c *fiber.Ctx, err error) error {
	if err != nil {
		logger.Error.Println(err.Error())
		return result.SignatureInvalidResult(c, result.SignatureInvalid)
	}
	return err
}
