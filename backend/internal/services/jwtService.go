package services

import (
	"FreeConnect/internal/models"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTService interface {
	GenerateToken(user *models.User) (string, error)
	ValidateToken(encodedToken string) (*jwt.Token, error)
}

type jwtService struct {
	secretKey string
	issuer    string
}

func NewJWTService() JWTService {
	// Usually you store your JWT secret in an environment variable:
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "ChangeThisSecretInProduction" // fallback for local dev
	}
	return &jwtService{
		secretKey: secret,
		issuer:    "FreeConnect",
	}
}

type CustomClaims struct {
	UserID   uint   `json:"user_id"`
	UserRole string `json:"user_role"`
	jwt.RegisteredClaims
}

func (j *jwtService) GenerateToken(user *models.User) (string, error) {
	claims := CustomClaims{
		UserID:   user.ID,
		UserRole: user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    j.issuer,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)), // 1-day expiry
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(j.secretKey))
}

// ValidateToken parses and validates a given JWT string
func (j *jwtService) ValidateToken(encodedToken string) (*jwt.Token, error) {
	return jwt.ParseWithClaims(encodedToken, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(j.secretKey), nil
	})
}
