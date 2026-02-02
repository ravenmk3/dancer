package auth

import (
	"time"

	"dancer/internal/config"
	"dancer/internal/errors"
	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	UserType string `json:"user_type"`
	jwt.RegisteredClaims
}

func GenerateToken(userID, username, userType string) (string, error) {
	cfg := config.GetConfig()

	claims := Claims{
		UserID:   userID,
		Username: username,
		UserType: userType,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(cfg.JWT.Expiry) * time.Second)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(cfg.JWT.Secret))
}

func ValidateToken(tokenString string) (*Claims, error) {
	cfg := config.GetConfig()

	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.ErrInvalidToken
		}
		return []byte(cfg.JWT.Secret), nil
	})

	if err != nil {
		if err == jwt.ErrTokenExpired {
			return nil, errors.ErrTokenExpired
		}
		return nil, errors.ErrInvalidToken
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.ErrInvalidToken
}
