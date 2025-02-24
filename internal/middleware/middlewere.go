package middleware

import (
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/golang-jwt/jwt/v5"
)

type Middleware struct {
	secretKey []byte
	logger    *slog.Logger
}

func NewMiddleware(secretKey string, logger *slog.Logger) *Middleware {
	return &Middleware{
		secretKey: []byte(secretKey),
		logger:    logger,
	}
}

func (m *Middleware) AuthMiddleware() fiber.Handler {
	return func(c fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "missing or invalid authorization header"})
		}

		token := strings.TrimPrefix(authHeader, "Bearer ")

		parsedToken, err := jwt.Parse(token, func(j *jwt.Token) (interface{}, error) {
			if _, ok := j.Method.(*jwt.SigningMethodHMAC); !ok {
				err := fmt.Errorf("unexpected signing method: %v", j.Header["alg"])
				m.logger.Error("Invalid token signing method", "token", token, "error", err)
				return nil, err
			}
			return m.secretKey, nil
		})
		if err != nil || !parsedToken.Valid {
			m.logger.Error("Invalid token", "token", token, "error", err)
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid token"})
		}

		m.logger.Info("Token validated successfully")
		return c.Next()
	}
}

func (m *Middleware) LoggerMiddleware() fiber.Handler {
	return func(c fiber.Ctx) error {
		start := time.Now()

		m.logger.Info("request started",
			slog.String("method", c.Method()),
			slog.String("path", c.Path()),
			slog.String("ip", c.IP()),
		)

		err := c.Next()

		m.logger.Info("request completed",
			slog.String("method", c.Method()),
			slog.String("path", c.Path()),
			slog.Int("status", c.Response().StatusCode()),
			slog.Duration("duration", time.Since(start)),
		)

		return err
	}
}
