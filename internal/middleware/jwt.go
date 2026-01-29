package middleware

import (
	"backend-blog/internal/constant"
	"backend-blog/internal/model/dto"
	"backend-blog/internal/pkg/response"
	"context"
	"crypto/rsa"
	"log/slog"

	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

var (
	PrivateKey     *rsa.PrivateKey
	AuthMiddleware fiber.Handler
)

// InitJWT 初始化 JWT 配置。建议在 main.go 加载完配置后显式调用，而不是用 init()
func InitJWT(privateKey *rsa.PrivateKey) {
	PrivateKey = privateKey

	AuthMiddleware = jwtware.New(jwtware.Config{
		SigningKey: jwtware.SigningKey{
			JWTAlg: jwtware.RS256,
			Key:    PrivateKey.Public(),
		},
		SuccessHandler: func(c *fiber.Ctx) error {
			token := c.Locals("user").(*jwt.Token)
			claims := token.Claims.(jwt.MapClaims)
			uid := uint64(claims["uid"].(float64))
			username, _ := claims["username"].(string)
			userInfo := dto.UserContextInfo{
				ID:       uid,
				Username: username,
			}
			ctx := context.WithValue(c.UserContext(), constant.ContextUserKey, userInfo)
			c.SetUserContext(ctx)

			return c.Next()
		},
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			slog.Warn("JWT Auth Error", "error", err)
			return response.Unauthorized(c, "无效的身份令牌")
		},
	})
}

// GetUserId 从 Context 中获取当前登录用户的 ID
func GetUserId(ctx context.Context) uint64 {
	if info, ok := ctx.Value(constant.ContextUserKey).(dto.UserContextInfo); ok {
		return info.ID
	}
	return 0
}

// GetUsername 从 Context 中获取当前登录用户的用户名
func GetUsername(ctx context.Context) string {
	if info, ok := ctx.Value(constant.ContextUserKey).(dto.UserContextInfo); ok {
		return info.Username
	}
	return ""
}
