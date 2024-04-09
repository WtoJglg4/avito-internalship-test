package service

import (
	"crypto/sha256"
	"fmt"
	"github/avito/entities"
	"github/avito/pkg/repository"
	"time"

	"github.com/dgrijalva/jwt-go"
)

const (
	sault      = "762t38cveteiy7fte9c2t18r723ef26gt86fvt2e9bc"
	signingKey = "scue5i6v7q$%^&*!(@)_byc3[nfe9udaf098ay-]"
	tokenTTL   = 12 * time.Hour
)

type tokenClaims struct {
	jwt.StandardClaims
	UserId   int    `json:"user_id"`
	UserRole string `json:"user_role"`
}

type AuthService struct {
	repo *repository.Repository
}

func NewAuthService(repo *repository.Repository) *AuthService {
	return &AuthService{repo: repo}
}

func (s *AuthService) CreateUser(user entities.User) (int, error) {
	user.Password = generatePasswordHash(user.Password)
	return s.repo.CreateUser(user)
}

func (s *AuthService) GenerateToken(login, password string) (string, error) {
	user, err := s.repo.GetUser(login, generatePasswordHash(password))
	if err != nil {
		return "", err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &tokenClaims{
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(tokenTTL).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		user.Id,
		user.Role,
	})
	return token.SignedString([]byte(signingKey))
}

func generatePasswordHash(password string) string {
	hash := sha256.New()
	hash.Sum([]byte(password))
	return fmt.Sprintf("%x", string(hash.Sum([]byte(sault))))
}
