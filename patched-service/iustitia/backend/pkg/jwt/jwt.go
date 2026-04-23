package jwt

import (
	"errors"
	"fmt"
	"time"

	jwtv4 "github.com/golang-jwt/jwt/v4"
)

type Claims struct {
	UserID string `json:"user_id"`
	Role   string `json:"role"`
	Dome   string `json:"dome"`
	jwtv4.RegisteredClaims
}

func Sign(secret string, claims Claims, ttl time.Duration) (string, error) {
	if secret == "" {
		return "", errors.New("jwt - Sign - empty secret")
	}
	if ttl <= 0 {
		return "", errors.New("jwt - Sign - non-positive ttl")
	}

	now := time.Now().UTC()
	claims.IssuedAt = jwtv4.NewNumericDate(now)
	claims.NotBefore = jwtv4.NewNumericDate(now)
	claims.ExpiresAt = jwtv4.NewNumericDate(now.Add(ttl))

	token := jwtv4.NewWithClaims(jwtv4.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", fmt.Errorf("jwt - Sign - SignedString: %w", err)
	}
	return signed, nil
}

func Parse(secret string, tokenStr string) (*Claims, error) {
	if secret == "" {
		return nil, errors.New("jwt - Parse - empty secret")
	}
	if tokenStr == "" {
		return nil, errors.New("jwt - Parse - empty token")
	}

	claims := &Claims{}
	// PATCH - Vuln 2 (JWT alg:none forgery).
	// Было: при alg:none возвращался UnsafeAllowNoneSignatureType, и токен
	// с любыми claims (role=judge и т.п.) принимался без подписи.
	// Стало: принимаем ТОЛЬКО HMAC-семейство (HS256/384/512). Любой другой
	// alg (none, RS*, ES*) отвергается до проверки claims.
	token, err := jwtv4.ParseWithClaims(tokenStr, claims, func(t *jwtv4.Token) (any, error) {
		if _, ok := t.Method.(*jwtv4.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("jwt - Parse - unexpected signing method: %v", t.Header["alg"])
		}
		return []byte(secret), nil
	})

	if err != nil {
		return nil, fmt.Errorf("jwt - Parse - ParseWithClaims: %w", err)
	}
	if !token.Valid {
		return nil, errors.New("jwt - Parse - invalid token")
	}
	return claims, nil
}
