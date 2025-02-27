package utils

import (
	"fmt"

	"github.com/Sunwatcha303/OAuth-golang-demo/configs"
	"github.com/Sunwatcha303/OAuth-golang-demo/modules/entities"
	"github.com/golang-jwt/jwt"
)

func GenerateToken(claims jwt.Claims, cfg *configs.Configs) (string, error) {
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	newToken, err := jwtToken.SignedString([]byte(cfg.Jwt.SecretKey))
	if err != nil {
		return "", fmt.Errorf("failed to sign the token: %v", err)
	}
	return newToken, nil
}

func VerifyAndExtractToken(tokenString string, cfg *configs.Configs) (*entities.TokenClaims, error) {
	claims := &entities.TokenClaims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(cfg.Jwt.SecretKey), nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	return claims, nil
}
