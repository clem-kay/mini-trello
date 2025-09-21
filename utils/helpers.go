package utils

import (
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

var JWTSecret = []byte(GetEnv("JWT_SECRET", "your_super_secret_key_here_at_least_32_chars"))
var JWTExpiresIn = GetEnv("JWT_EXPIRES_IN", "24h")

func GetEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func GenerateJWT(userID uint, email string) (string, error) {
	// Parse duration
	duration, err := time.ParseDuration(JWTExpiresIn)
	if err != nil {
		return "", err
	}

	// Create the Claims
	claims := jwt.MapClaims{
		"user_id": userID,
		"email":   email,
		"exp":     time.Now().Add(duration).Unix(),
	}

	// Create token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Generate encoded token
	t, err := token.SignedString(JWTSecret)
	if err != nil {
		return "", err
	}

	return t, nil
}

func GetUserIDFromToken(tokenString string) (uint, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Ensure token is signed with HMAC
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return JWTSecret, nil
	})

	if err != nil {
		return 0, err // e.g., signature invalid, expired, etc.
	}

	// Check if token is valid and has expected claims
	if !token.Valid {
		return 0, jwt.ErrTokenInvalidId
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return 0, jwt.ErrTokenMalformed
	}

	userIDFloat, ok := claims["user_id"].(float64)
	if !ok {
		return 0, jwt.ErrTokenMalformed
	}

	return uint(userIDFloat), nil
}

func GetIDParam(c *fiber.Ctx) string {
	id := c.Params("id")
	if id == "" {
		c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "ID is required",
		})
		return ""
	}
	return id
}

func GetIDFromContext(c *fiber.Ctx) uint {
	userID, ok := c.Locals("user_id").(uint)

	if !ok {
		return 0
	}
	return userID
}
