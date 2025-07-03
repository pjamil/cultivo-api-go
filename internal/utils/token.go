package utils

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// TODO: Mover para configuração / variável de ambiente
var jwtSecret = []byte("supersecretkey")

// GenerateToken cria um novo token JWT para um usuário.
func GenerateToken(userID uint) (string, error) {
	// O token expira em 24 horas
	expirationTime := time.Now().Add(24 * time.Hour)

	claims := &jwt.RegisteredClaims{
		Subject:   fmt.Sprint(userID),
		ExpiresAt: jwt.NewNumericDate(expirationTime),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(jwtSecret)
}

// ValidateToken verifica se um token é válido e retorna o ID do usuário.
func ValidateToken(tokenString string) (string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return jwtSecret, nil
	})

	if err != nil {
		return "", err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims["sub"].(string), nil
	} else {
		return "", fmt.Errorf("invalid token")
	}
}
