package service

import (
	"fmt"
	"os"
	"time"

	"example.com/mygamelist/errorutils"
	"example.com/mygamelist/interfaces"
	"example.com/mygamelist/utils"
	"github.com/golang-jwt/jwt/v5"
)

// UserService defines a user service controller.
type UserService struct {
	UserRepository interfaces.UserRepository
	AuthService    interfaces.AuthService
}

// NewUserService creates a new user service controller.
func NewUserService(repo interfaces.UserRepository, auth interfaces.AuthService) *UserService {
	return &UserService{
		UserRepository: repo,
		AuthService:    auth,
	}
}

// RegisterUser authenticates user.
func (s *UserService) RegisterUser(username, email, password string) (int64, error) {

	isUser, err := s.UserRepository.SelectUserByUsername(username)
	if err != nil {
		return 0, fmt.Errorf("failed to select user: %w", err)
	}

	if isUser {
		return 0, errorutils.ErrUserExists
	}

	hashedPassword, err := s.AuthService.HashPassword(password)
	if err != nil {
		return 0, fmt.Errorf("failed to hash password: %w", err)
	}

	userId, err := s.UserRepository.InsertUser(username, email, hashedPassword)
	if err != nil {
		return 0, fmt.Errorf("failed to insert user: %w", err)
	}

	return userId, nil
}

// LoginUser authorizes user.
func (s *UserService) LoginUser(username, password string) (string, int, error) {
	k := os.Getenv("JWT_SECRET_KEY")
	var secretKey = []byte(k)
	expirationTime := time.Now().Add(5 * time.Minute).Unix()
	hashedPassword, err := s.UserRepository.PasswordByUsername(username)
	if err != nil {
		return "", 0, fmt.Errorf("failed to retrieve password: %w", err)
	}
	if !utils.CheckPasswordHash(hashedPassword, password) {
		return "", 0, errorutils.ErrPasswordMatch
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": username,
		"exp":      expirationTime,
	})

	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", 0, fmt.Errorf("failed to sign token: %w", err)
	}

	userId, err := s.UserRepository.SelectUserIdByUsername(username)
	if err != nil {
		return "", 0, fmt.Errorf("failed to retrieve user ID: %w", err)
	}

	return tokenString, userId, err

}

// StoreRefreshToken adds refreshtoken for a given user.
func (s *UserService) StoreRefreshToken(username, refreshToken, jti string) error {
	userId, err := s.UserRepository.SelectUserIdByUsername(username)
	if err != nil {
		return err
	}

	err = s.UserRepository.InsertRefreshToken(userId, refreshToken, jti)
	if err != nil {
		return fmt.Errorf("failed to insert refreshtoken: %w", err)
	}

	return nil
}

// FetchRefreshToken retrieves refreshtoken for a given user.
func (s *UserService) FetchRefreshToken(username string, userId int) (string, error) {
	_, jtiFromDb, err := s.UserRepository.RefreshTokenById(userId)
	if err != nil {
		return "", fmt.Errorf("failed to fetch refreshtoken: %w", err)
	}
	return jtiFromDb, nil
}
