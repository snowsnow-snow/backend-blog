package util

import (
	"backend-blog/logger"
	"backend-blog/result"
	"crypto/rand"
	"crypto/rsa"
	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
)

var (
	privateKey    *rsa.PrivateKey
	JWTMiddleware fiber.Handler
)

func init() {
	rng := rand.Reader
	var err error
	privateKey, err = rsa.GenerateKey(rng, 2048)
	if err != nil {
		logger.Error.Fatalf("rsa.GenerateKey: %v", err)
	}
	JWTMiddleware = jwtware.New(jwtware.Config{
		SigningKey: jwtware.SigningKey{
			JWTAlg: jwtware.RS256,
			Key:    privateKey.Public(),
		},
		ErrorHandler: JwtError,
	})
}

var JwtKey = []byte("your_secret_key")

func TokenAuth() fiber.Handler {
	return func(c *fiber.Ctx) error {
		//// 获取 Authorization 头部的 Token
		//authHeader := c.Get("Authorization")
		//tokenString := authHeader[7:] // 去除 "Bearer " 前缀
		//
		//// 解析和验证 Token
		//token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		//	// 验证签名的方法
		//	if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
		//		return nil, fiber.ErrUnauthorized
		//	}
		//	return JwtKey, nil
		//})
		//
		//// 处理验证结果
		//if err != nil || !token.Valid {
		//	return fiber.ErrUnauthorized
		//}
		//
		//// 验证通过，将用户标识信息存储在上下文中
		//claims, _ := token.Claims.(jwt.MapClaims)
		//c.Locals("user", claims)

		return c.Next()
	}
}
func JwtError(c *fiber.Ctx, err error) error {
	if err != nil {
		logger.Error.Println(err.Error())
		return result.ErrorResult(c, result.SignatureInvalid)
	}
	return err
}
