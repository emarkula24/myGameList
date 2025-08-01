package utils

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"example.com/mygamelist/errorutils"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct{}

func (AuthService) HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(hashed, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashed), []byte(password))
	return err == nil
}

func VerifyToken(tokenString string) error {
	k := os.Getenv("JWT_SECRET_KEY")
	var secretKey = []byte(k)
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		return secretKey, nil
	}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))
	if err != nil {
		return err
	}

	if !token.Valid {
		return errorutils.ErrInvalidToken
	}
	switch {
	case token.Valid:
		return nil
	case errors.Is(err, jwt.ErrTokenExpired):
		return jwt.ErrTokenExpired
	}
	return nil
}

func GenerateRefreshToken(username string) (string, string, error) {
	k := os.Getenv("REFRESH_SECRET_KEY")
	var secretKey = []byte(k)
	expirationTime := time.Now().Add(24 * time.Hour).Unix()
	jti := fmt.Sprintf("%s-%d", username, time.Now().UnixNano())

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username":  username,
		"ExpiresAt": expirationTime,
		"Issuer":    "mygamelist-back",
		"jti":       jti,
	})

	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", "", fmt.Errorf("failed to sign jwt token: %w", err)
	}
	return tokenString, jti, nil
}

func CreateFingerPrintCookie(jwtToken string) (*http.Cookie, error) {

	cookie := &http.Cookie{
		Name:     "access_token",
		Value:    jwtToken,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
	}
	return cookie, nil
}

func CreateRefreshTokenCookie(refreshToken string) (*http.Cookie, error) {

	secure := false
	if os.Getenv("MODE") == "production" {
		log.Println("using production")
		secure = true
	}

	cookie := &http.Cookie{
		Name:     "refreshToken",
		Value:    refreshToken,
		HttpOnly: true,
		Secure:   secure,
		SameSite: http.SameSiteStrictMode,
		MaxAge:   int(24 * time.Hour.Seconds()),
	}
	return cookie, nil
}

func IsRefreshTokenValid(db *sql.DB, userID int, token string) (bool, error) {
	var exists bool
	query := `
			SELECT EXISTS (
				SELECT 1 FROM refresh_tokens 
				WHERE user_id = $1 AND token = $2 
				AND revoked = FALSE AND expires_at > NOW()
			)`
	err := db.QueryRow(query, userID, token).Scan(&exists)
	return exists, err
}
