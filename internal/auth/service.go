package auth

import (
	"ai-resume-analizyzer/internal/user"
	"fmt"
	"log"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	Register(req *RegisterRequest) (*AuthResponse, error)
	Login(req *LoginRequest) (*AuthResponse, error)
}

type authService struct {
	userRepo  user.UserRepository
	jwtSecret string
}

func NewAuthService(userRepo user.UserRepository, jwtSecret string) AuthService {
	return &authService{userRepo, jwtSecret}
}

func (s *authService) Register(req *RegisterRequest) (*AuthResponse, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password")
	}

	user := &user.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: string(hashed),
	}

	if err := s.userRepo.Create(user); err != nil {
		log.Printf("Error creating user: %v", err)
		return nil, fmt.Errorf("email already registered")
	}

	token, err := s.generateToken(user)
	if err != nil {
		return nil, err
	}

	return &AuthResponse{Token: token, User: *user}, nil
}

func (s *authService) Login(req *LoginRequest) (*AuthResponse, error) {
	user, err := s.userRepo.FindByEmail(req.Email)
	if err != nil {
		return nil, fmt.Errorf("invalid credentials")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, fmt.Errorf("invalid credentials")
	}

	token, err := s.generateToken(user)
	if err != nil {
		return nil, err
	}

	return &AuthResponse{Token: token, User: *user}, nil
}

func (s *authService) generateToken(user *user.User) (string, error) {
	claims := jwt.MapClaims{
		"user_id": user.ID.String(),
		"email":   user.Email,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.jwtSecret))
}
