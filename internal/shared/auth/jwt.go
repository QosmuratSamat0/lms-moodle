package auth

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"time"

	"github.com/ap1-final-mini-moodle/internal/domain/auth"
	"github.com/golang-jwt/jwt/v5"
)

type JWTConfig struct {
	SecretKey            []byte
	AccessTokenDuration  time.Duration
	RefreshTokenDuration time.Duration
	Issuer               string
	Audience             string
}

func GenerateTokens(user *UserData, config *JWTConfig) (*auth.Token, error) {
	refreshToken, err := generateRandomToken()
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	accessToken, err := GenerateAccessToken(user, refreshToken, config)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	return &auth.Token{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    int64(config.AccessTokenDuration.Seconds()),
		IssuedAt:     now,
	}, nil
}

func GenerateAccessToken(user *UserData, refreshToken string, config *JWTConfig) (string, error) {
	now := time.Now()
	expiresAt := now.Add(config.AccessTokenDuration)

	claims := &auth.TokenClaims{
		UserID:       user.ID,
		Email:        user.Email,
		FirstName:    user.FirstName,
		LastName:     user.LastName,
		Role:         user.Role,
		RefreshToken: refreshToken,
		IssuedAt:     now.Unix(),
		ExpiresAt:    expiresAt.Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":       claims.UserID,
		"email":         claims.Email,
		"first_name":    claims.FirstName,
		"last_name":     claims.LastName,
		"role":          claims.Role,
		"refresh_token": claims.RefreshToken,
		"iat":           claims.IssuedAt,
		"exp":           claims.ExpiresAt,
		"iss":           config.Issuer,
		"aud":           config.Audience,
	})

	return token.SignedString(config.SecretKey)
}

func VerifyAccessToken(tokenString string, config *JWTConfig) (*auth.TokenClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return config.SecretKey, nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("invalid token claims")
	}

	if exp, ok := claims["exp"].(float64); ok {
		if time.Now().Unix() > int64(exp) {
			return nil, fmt.Errorf("token expired")
		}
	}

	return &auth.TokenClaims{
		UserID:       claims["user_id"].(string),
		Email:        claims["email"].(string),
		FirstName:    claims["first_name"].(string),
		LastName:     claims["last_name"].(string),
		Role:         claims["role"].(string),
		RefreshToken: claims["refresh_token"].(string),
		IssuedAt:     int64(claims["iat"].(float64)),
		ExpiresAt:    int64(claims["exp"].(float64)),
	}, nil
}

func generateRandomToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

type UserData struct {
	ID        string
	Email     string
	FirstName string
	LastName  string
	Role      string
}
