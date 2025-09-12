package auth

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	UserID    string `json:"user_id"`
	Role      string `json:"role"`
	TokenType string `json:"token_type"` // "access" or "refresh"
	jwt.RegisteredClaims
}

func GenerateTokens(secret string, userID, role string, accessTTLMinutes int, refreshTTLDays int) (string, string, error) {
	now := time.Now()

	// access token (JWT)
	accessClaims := Claims{
		UserID:    userID,
		Role:      role,
		TokenType: "access",
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(time.Duration(accessTTLMinutes) * time.Minute)),
		},
	}
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	at, err := accessToken.SignedString([]byte(secret))
	if err != nil {
		return "", "", err
	}

	// refresh token — генерируем крипто-стойкий plain (hex)
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", "", err
	}
	plainRefresh := hex.EncodeToString(b)

	return at, plainRefresh, nil
}

func ParseAccessToken(secret, tokenStr string) (*Claims, error) {
	return parseTokenOfType(secret, tokenStr, "access")
}

// ParseRefreshToken validates refresh token and returns claims.
func ParseRefreshToken(secret, tokenStr string) (*Claims, error) {
	return parseTokenOfType(secret, tokenStr, "refresh")
}

func parseTokenOfType(secret, tokenStr, wantType string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}
	if claims.TokenType != wantType {
		return nil, errors.New("invalid token type")
	}
	return claims, nil
}
func HashToken(token string) string {
	h := sha256.Sum256([]byte(token))
	return hex.EncodeToString(h[:])
}
