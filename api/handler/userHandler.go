package handler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"example.com/mygamelist/errorutils"
	"example.com/mygamelist/service"
	"example.com/mygamelist/utils"
	"github.com/golang-jwt/jwt/v5"
)

// UserHandler defines a user HTTP handler.
type UserHandler struct {
	UserService *service.UserService
}

// NewUserHandler creates a new user HTTP handler.
func NewUserHandler(us *service.UserService) *UserHandler {
	return &UserHandler{UserService: us}
}

// Register handles POST /user/register requests.
func (h *UserHandler) Register(w http.ResponseWriter, req *http.Request) {

	type RegisterRequest struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	var regReq RegisterRequest
	decoder := json.NewDecoder(req.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&regReq); err != nil {
		log.Printf("invalid json %s", err)
		errorutils.Write(w, "", http.StatusInternalServerError)
		return
	}
	if len(regReq.Password) <= 6 {
		errorutils.Write(w, "password too short", http.StatusBadRequest)
		return
	}
	userId, err := h.UserService.RegisterUser(regReq.Username, regReq.Email, regReq.Password)
	if err != nil {
		switch {
		case errors.Is(err, errorutils.ErrUserExists):
			log.Printf("user already exists %s", err)
			errorutils.Write(w, "user already exists", http.StatusConflict)
			return
		default:
			log.Printf("failed to register user: %s", err)
			errorutils.Write(w, "", http.StatusInternalServerError)
			return
		}

	}
	type RegisterResponse struct {
		UserID int64 `json:"user_id"`
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(RegisterResponse{UserID: userId})
	if err != nil {
		log.Printf("failed to register user: %s", err)
		errorutils.Write(w, " ", http.StatusInternalServerError)
		return
	}
}

// Login handles POST /user/login requests.
func (h *UserHandler) Login(w http.ResponseWriter, req *http.Request) {

	type LoginRequest struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	var loginReq LoginRequest
	decoder := json.NewDecoder(req.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&loginReq); err != nil {
		log.Printf("invalid json")
		errorutils.Write(w, "", http.StatusInternalServerError)
		return
	}
	jwtToken, userId, err := h.UserService.LoginUser(loginReq.Username, loginReq.Password)
	userIDStr := strconv.Itoa(userId)
	if err != nil {
		switch {
		case errors.Is(err, errorutils.ErrPasswordMatch):
			log.Printf("failed to login user: %s", err)
			errorutils.Write(w, "incorrect username or password", http.StatusUnauthorized)
			return
		default:
			log.Printf("failed to login user: %s", err)
			errorutils.Write(w, "authentication failed", http.StatusUnauthorized)
			return
		}
	}

	type LoginResponse struct {
		AccessToken string `json:"accessToken"`
		UserId      string `json:"userId"`
		Username    string `json:"username"`
	}
	username := loginReq.Username
	refreshToken, jti, err := utils.GenerateRefreshToken(username)
	if err != nil {
		log.Printf("failed to login user: %s", err)
		errorutils.Write(w, "authentication failed", http.StatusUnauthorized)
		return
	}
	refreshTokenCookie, err := utils.CreateRefreshTokenCookie(refreshToken)
	if err != nil {
		log.Printf("failed to login user: %s", err)
		errorutils.Write(w, "authentication failed", http.StatusUnauthorized)
		return
	}

	err = h.UserService.StoreRefreshToken(loginReq.Username, refreshToken, jti)
	if err != nil {
		log.Printf("failed to login user: %s", err)
		errorutils.Write(w, "authentication failed", http.StatusUnauthorized)
		return
	}
	http.SetCookie(w, refreshTokenCookie)
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(LoginResponse{AccessToken: jwtToken, UserId: userIDStr, Username: username})
	if err != nil {
		log.Printf("failed to login user: %s", err)
		errorutils.Write(w, "authentication failed", http.StatusUnauthorized)
		return
	}
}

// Logout handles POST user/logout requests.
func (h *UserHandler) Logout(w http.ResponseWriter, r *http.Request) {
	type LogoutRequest struct {
		Username string `json:"username"`
		UserId   string `json:"userId"`
	}
	var logoutReq LogoutRequest
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&logoutReq); err != nil {
		log.Printf("failed to decode logoutreq: %s", err)
		errorutils.Write(w, "", http.StatusInternalServerError)
		return
	}
	if logoutReq.UserId == "" || logoutReq.Username == "" {
		errorutils.Write(w, "missing values from body", http.StatusBadRequest)
		return
	}
	userIdInt, err := strconv.Atoi(logoutReq.UserId)
	if err != nil {
		log.Printf("failed to convert userId")
		errorutils.Write(w, "", http.StatusInternalServerError)
		return

	}
	cookie, err := r.Cookie("refreshToken")
	log.Printf("cookie: %s", cookie)
	if err != nil {
		log.Printf("Failed to logout, missing refresh token: %s", err)
		errorutils.Write(w, "missing refresh token", http.StatusUnauthorized)
		return
	}
	tokenStr := cookie.Value
	log.Println(tokenStr)
	k := os.Getenv("REFRESH_SECRET_KEY")
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (any, error) {
		return []byte(k), nil
	}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))
	if err != nil {
		log.Printf("Failed to logout, invalid refresh token: %s", err)
		errorutils.Write(w, "invalid refresh token", http.StatusUnauthorized)
		return
	}

	jtiFromDb, err := h.UserService.FetchRefreshToken(logoutReq.Username, userIdInt)
	if err != nil {
		log.Printf("Failed to logout, invalid user: %s", err)
		errorutils.Write(w, "invalid user", http.StatusUnauthorized)
		return
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		log.Printf("Failed to logout, invalid token claims: %s", err)
		errorutils.Write(w, "", http.StatusUnauthorized)
		return
	}

	jti, ok := claims["jti"].(string)

	if !ok {
		log.Printf("Failed to logout, invalid claims: %s", err)
		errorutils.Write(w, "", http.StatusUnauthorized)
		return
	}
	if jti == jtiFromDb {
		switch {
		case token.Valid:
			err := h.UserService.UserRepository.DeleteRefreshToken(userIdInt, jti)
			if err != nil {
				log.Printf("failed to delete refresh token from database: %s", err)
				errorutils.Write(w, "", http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)

		}
	}
}

// Refresh handles POST /user/refresh requests.
func (h *UserHandler) Refresh(w http.ResponseWriter, req *http.Request) {

	type RefreshRequest struct {
		Username string `json:"username"`
		UserId   string `json:"userId"`
	}
	var refreshReq RefreshRequest
	decoder := json.NewDecoder(req.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&refreshReq); err != nil {
		log.Printf("failed to refresh, invalid JSON: %s", err)
		errorutils.Write(w, "", http.StatusInternalServerError)
		return
	}
	if refreshReq.UserId == "" || refreshReq.Username == "" {
		errorutils.Write(w, "missing values from body", http.StatusBadRequest)
		return
	}
	userIdInt, err := strconv.Atoi(refreshReq.UserId)
	if err != nil {
		errorutils.Write(w, "", http.StatusInternalServerError)
		return

	}
	cookie, err := req.Cookie("refreshToken")
	if err != nil {
		log.Printf("failed to refresh, missing refresh token: %s", err)
		errorutils.Write(w, "missing refresh token", http.StatusUnauthorized)
		return
	}

	tokenStr := cookie.Value
	log.Println(tokenStr)
	k := os.Getenv("REFRESH_SECRET_KEY")
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (any, error) {
		return []byte(k), nil
	}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))
	if err != nil {
		log.Printf("failed to refresh, invalid refresh token: %s", err)
		errorutils.Write(w, "invalid refresh token", http.StatusUnauthorized)
		return
	}

	jtiFromDb, err := h.UserService.FetchRefreshToken(refreshReq.Username, userIdInt)
	if err != nil {
		log.Printf("failed to refresh, invalid user: %s", err)
		errorutils.Write(w, "invalid user", http.StatusUnauthorized)
		return
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		log.Printf("failed to refresh, invalid token claims: %s", err)
		errorutils.Write(w, "", http.StatusUnauthorized)
		return
	}

	jti, ok := claims["jti"].(string)

	if !ok {
		log.Printf("failed to refresh, invalid claims: %s", err)
		errorutils.Write(w, "", http.StatusUnauthorized)
		return
	}
	type RefreshResponse struct {
		AccessToken string `json:"accessToken"`
	}

	if jti == jtiFromDb {
		switch {
		case token.Valid:
			k := os.Getenv("JWT_SECRET_KEY")
			var secretKey = []byte(k)
			expirationTime := time.Now().Add(5 * time.Minute).Unix()
			token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
				"username": refreshReq.Username,
				"exp":      expirationTime,
			})
			tokenString, err := token.SignedString(secretKey)
			if err != nil {
				log.Printf("failed to sign refresh token: %s", err)
				errorutils.Write(w, "", http.StatusUnauthorized)
			}
			w.Header().Set("Content-Type", "application/json")
			err = json.NewEncoder(w).Encode(RefreshResponse{AccessToken: tokenString})
			if err != nil {
				log.Printf("failed to sign refresh token: %s", err)
				errorutils.Write(w, "", http.StatusUnauthorized)
			}
			fmt.Println("refresh successfull")
		case errors.Is(err, jwt.ErrTokenMalformed):
			log.Printf("malformed token: %s", err)
			errorutils.Write(w, "", http.StatusUnauthorized)
		case errors.Is(err, jwt.ErrTokenSignatureInvalid):
			// Invalid signature
			log.Printf("invalid token signature: %s", err)
			errorutils.Write(w, "", http.StatusUnauthorized)
		case errors.Is(err, jwt.ErrTokenExpired) || errors.Is(err, jwt.ErrTokenNotValidYet):
			err := h.UserService.UserRepository.DeleteRefreshToken(userIdInt, jti)
			if err != nil {
				log.Printf("Token expired, deleted from database: %s", err)
				errorutils.Write(w, "invalid refresh token", http.StatusUnauthorized)
			}
		}
	} else {
		log.Printf("unexpected error: %s", err)
		errorutils.Write(w, "authentication failed", http.StatusUnauthorized)
	}
}

// GetUsers handles GET users requests.
func (h *UserHandler) GetUsers(w http.ResponseWriter, req *http.Request) {
	ctx := context.Background()

	users, err := h.UserService.FetchUsers(ctx)
	if err != nil {
		log.Printf("failed to get users: %s", err)
		errorutils.Write(w, "userfetch  failed", http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)

}
