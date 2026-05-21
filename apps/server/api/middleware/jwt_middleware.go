package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v3"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"

	"omni-pixel/domain"
	"omni-pixel/internal/response"
)

func JWTAuthMiddleware(secret string) fiber.Handler {
	return func(c fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return response.Write(c, response.ErrUnauthorized)
		}

		tokenString := strings.TrimSpace(strings.TrimPrefix(authHeader, "Bearer "))
		if tokenString == "" || tokenString == authHeader {
			return response.Write(c, response.ErrUnauthorized)
		}

		claims := &domain.JwtCustomClaims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(secret), nil
		}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))
		if err != nil || !token.Valid {
			return response.Write(c, response.ErrUnauthorized)
		}

		fiber.Locals[uuid.UUID](c, "user_id", claims.UserID)
		fiber.Locals[string](c, "email", claims.Email)
		return c.Next()
	}
}
