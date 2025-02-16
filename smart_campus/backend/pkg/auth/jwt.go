package auth

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	ErrInvalidToken         = errors.New("invalid token")
	ErrExpiredToken         = errors.New("token has expired")
	ErrInvalidClaims        = errors.New("invalid token claims")
	ErrInvalidSigningMethod = errors.New("invalid signing method")
	ErrTokenNotValidYet     = errors.New("token not valid yet")
	ErrInvalidTokenType     = errors.New("invalid token type")
)

// TokenType represents the type of JWT token
type TokenType string

const (
	AccessToken  TokenType = "access"
	RefreshToken TokenType = "refresh"
)

// Claims represents the JWT claims
type Claims struct {
	jwt.RegisteredClaims
	UserID    string    `json:"user_id"`
	Role      string    `json:"role"`
	DeviceID  string    `json:"device_id"`
	TokenType TokenType `json:"token_type"`
}

// JWTService handles JWT token operations
type JWTService struct {
	accessSecret  string
	refreshSecret string
	accessExp     time.Duration
	refreshExp    time.Duration
	issuer        string
}

// NewJWTService creates a new JWT service
func NewJWTService(accessSecret, refreshSecret string, accessExp, refreshExp time.Duration) *JWTService {
	return &JWTService{
		accessSecret:  accessSecret,
		refreshSecret: refreshSecret,
		accessExp:     accessExp,
		refreshExp:    refreshExp,
		issuer:        "smart_campus",
	}
}

// GenerateTokenPair generates an access and refresh token pair
func (s *JWTService) GenerateTokenPair(userID, role, deviceID string) (accessToken, refreshToken string, err error) {
	// Generate access token
	accessToken, err = s.generateToken(userID, role, deviceID, AccessToken, s.accessSecret, s.accessExp)
	if err != nil {
		return "", "", fmt.Errorf("failed to generate access token: %w", err)
	}

	// Generate refresh token
	refreshToken, err = s.generateToken(userID, role, deviceID, RefreshToken, s.refreshSecret, s.refreshExp)
	if err != nil {
		return "", "", fmt.Errorf("failed to generate refresh token: %w", err)
	}

	return accessToken, refreshToken, nil
}

// generateToken generates a new JWT token
func (s *JWTService) generateToken(userID, role, deviceID string, tokenType TokenType, secret string, expiry time.Duration) (string, error) {
	now := time.Now()
	claims := &Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(expiry)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			Issuer:    s.issuer,
			Subject:   userID,
		},
		UserID:    userID,
		Role:      role,
		DeviceID:  deviceID,
		TokenType: tokenType,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", fmt.Errorf("error signing token: %w", err)
	}

	return signedToken, nil
}

// ValidateAccessToken validates an access token
func (s *JWTService) ValidateAccessToken(tokenString string) (*Claims, error) {
	return s.validateToken(tokenString, s.accessSecret, AccessToken)
}

// ValidateRefreshToken validates a refresh token
func (s *JWTService) ValidateRefreshToken(tokenString string) (*Claims, error) {
	return s.validateToken(tokenString, s.refreshSecret, RefreshToken)
}

// validateToken validates a JWT token
func (s *JWTService) validateToken(tokenString, secret string, tokenType TokenType) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidSigningMethod
		}
		return []byte(secret), nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrExpiredToken
		}
		if errors.Is(err, jwt.ErrTokenNotValidYet) {
			return nil, ErrTokenNotValidYet
		}
		return nil, ErrInvalidToken
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, ErrInvalidClaims
	}

	if claims.TokenType != tokenType {
		return nil, ErrInvalidTokenType
	}

	if claims.Issuer != s.issuer {
		return nil, ErrInvalidToken
	}

	return claims, nil
}

// RefreshTokens refreshes an access token using a valid refresh token
func (s *JWTService) RefreshTokens(refreshToken string) (newAccessToken, newRefreshToken string, err error) {
	// Validate refresh token
	claims, err := s.ValidateRefreshToken(refreshToken)
	if err != nil {
		return "", "", fmt.Errorf("invalid refresh token: %w", err)
	}

	// Generate new token pair
	return s.GenerateTokenPair(claims.UserID, claims.Role, claims.DeviceID)
}

// ExtractTokenMetadata extracts metadata from a validated token
func (s *JWTService) ExtractTokenMetadata(claims *Claims) map[string]interface{} {
	return map[string]interface{}{
		"user_id":    claims.UserID,
		"role":       claims.Role,
		"device_id":  claims.DeviceID,
		"token_type": claims.TokenType,
		"exp":        claims.ExpiresAt,
		"iat":        claims.IssuedAt,
		"nbf":        claims.NotBefore,
		"iss":        claims.Issuer,
		"sub":        claims.Subject,
	}
}

// ParseTokenFromString parses a token string without validation
func (s *JWTService) ParseTokenFromString(tokenString string) (*Claims, error) {
	parser := jwt.NewParser()
	token, _, err := parser.ParseUnverified(tokenString, &Claims{})
	if err != nil {
		return nil, fmt.Errorf("error parsing token: %w", err)
	}

	claims, ok := token.Claims.(*Claims)
	if !ok {
		return nil, ErrInvalidClaims
	}

	return claims, nil
}
