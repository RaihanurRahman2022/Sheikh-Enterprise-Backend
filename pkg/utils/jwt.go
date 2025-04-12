package utils

import (
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type JWTClaims struct {
	UserID   uuid.UUID  `json:"user_id"`
	Username string     `json:"username"`
	Role     string     `json:"role"`
	ShopID   *uuid.UUID `json:"shop_id,omitempty"`
	jwt.RegisteredClaims
}

func GenerateJWT(userID uuid.UUID, username, role string, shopID *uuid.UUID) (string, error) {
	expHours := 24 // default to 24 hours
	if val := os.Getenv("JWT_EXPIRATION_HOURS"); val != "" {
		fmt.Sscanf(val, "%d", &expHours)
	}

	claims := JWTClaims{
		UserID:   userID,
		Username: username,
		Role:     role,
		ShopID:   shopID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * time.Duration(expHours))),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(os.Getenv("JWT_SECRET")))
}

func ValidateJWT(tokenString string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token claims")
}
