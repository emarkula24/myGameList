package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"example.com/mygamelist/errorutils"
	"example.com/mygamelist/service"
	"example.com/mygamelist/utils"
	"github.com/golang-jwt/jwt/v5"
)

// Defines dependencies for UserHandler struct
type UserHandler struct {
	UserService *service.UserService
}

// Creates a new instance of UserHandler
func NewUserHandler(us *service.UserService) *UserHandler {
	return &UserHandler{UserService: us}
}

func (h *UserHandler) Register(w http.ResponseWriter, req *http.Request) {

	if req.Header.Get("Content-Type") != "application/json" {
		http.Error(w, "Unsupported Media Type", http.StatusUnsupportedMediaType)
		return
	}

	type RegisterRequest struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	var regReq RegisterRequest
	decoder := json.NewDecoder(req.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&regReq); err != nil {
		errorutils.WriteJSONError(w, "Invalid JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	userId, err := h.UserService.RegisterUser(regReq.Username, regReq.Email, regReq.Password)
	if err != nil {
		switch {
		case errors.Is(err, errorutils.ErrUserExists):
			errorutils.WriteJSONError(w, "User already exists", http.StatusBadRequest)
			return
		default:
			log.Printf("Failed to register user: %s", err)
			errorutils.WriteJSONError(w, "Error adding user: "+err.Error(), http.StatusInternalServerError)
			return
		}

	}
	type RegisterResponse struct {
		UserID int64 `json:"user_id"`
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(RegisterResponse{UserID: userId})
	if err != nil {
		log.Printf("Failed to register user: %s", err)
		errorutils.WriteJSONError(w, "Error adding user: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *UserHandler) Login(w http.ResponseWriter, req *http.Request) {

	type LoginRequest struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	var loginReq LoginRequest
	decoder := json.NewDecoder(req.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&loginReq); err != nil {
		errorutils.WriteJSONError(w, "invalid JSON", http.StatusBadRequest)
		return
	}
	jwtToken, userId, err := h.UserService.LoginUser(loginReq.Username, loginReq.Password)
	if err != nil {
		switch {
		case errors.Is(err, errorutils.ErrPasswordMatch):
			log.Printf("Failed to login user: %s", err)
			errorutils.WriteJSONError(w, "incorrect username or password", http.StatusUnauthorized)
			return
		default:
			log.Printf("Failed to login user: %s", err)
			errorutils.WriteJSONError(w, "authentication failed", http.StatusUnauthorized)
			return
		}
	}

	type LoginResponse struct {
		AccessToken string `json:"accessToken"`
		UserId      int    `json:"userId"`
	}

	refreshToken, jti, err := utils.GenerateRefreshToken(loginReq.Username)
	if err != nil {
		log.Printf("Failed to login user: %s", err)
		errorutils.WriteJSONError(w, "authentication failed", http.StatusUnauthorized)
		return
	}
	refreshTokenCookie, err := utils.CreateRefreshTokenCookie(refreshToken)
	if err != nil {
		log.Printf("Failed to login user: %s", err)
		errorutils.WriteJSONError(w, "authentication failed", http.StatusUnauthorized)
		return
	}

	err = h.UserService.StoreRefreshToken(loginReq.Username, refreshToken, jti)
	if err != nil {
		log.Printf("Failed to login user: %s", err)
		errorutils.WriteJSONError(w, "authentication failed", http.StatusUnauthorized)
		return
	}
	http.SetCookie(w, refreshTokenCookie)
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(LoginResponse{AccessToken: jwtToken, UserId: userId})
	if err != nil {
		log.Printf("Failed to login user: %s", err)
		errorutils.WriteJSONError(w, "authentication failed", http.StatusUnauthorized)
		return
	}
}

func (h *UserHandler) Refresh(w http.ResponseWriter, req *http.Request) {

	type RefreshRequest struct {
		Username string `json:"username"`
	}
	var refreshReq RefreshRequest
	decoder := json.NewDecoder(req.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&refreshReq); err != nil {
		log.Printf("Failed to refresh, invalid JSON: %s", err)
		errorutils.WriteJSONError(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	cookie, err := req.Cookie("refreshToken")
	if err != nil {
		log.Printf("Failed to refresh, missing refresh token: %s", err)
		errorutils.WriteJSONError(w, "missing refresh token", http.StatusUnauthorized)
		return
	}

	tokenStr := cookie.Value
	k := os.Getenv("REFRESH_SECRET_KEY")
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return []byte(k), nil
	}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))
	if err != nil {
		log.Printf("Failed to refresh, invalid refresh token: %s", err)
		errorutils.WriteJSONError(w, "invalid refresh token", http.StatusUnauthorized)
		return
	}

	jtiFromDb, err := h.UserService.FetchRefreshToken(refreshReq.Username)
	if err != nil {
		log.Printf("Failed to refresh, invalid user: %s", err)
		errorutils.WriteJSONError(w, "invalid user", http.StatusUnauthorized)
		return
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		log.Printf("Failed to refresh, invalid token claims: %s", err)
		errorutils.WriteJSONError(w, "invalid token claims", http.StatusUnauthorized)
		return
	}

	jti, ok := claims["jti"].(string)

	if !ok {
		log.Printf("Failed to refresh, invalid claims: %s", err)
		errorutils.WriteJSONError(w, "jti claim missing or invalid", http.StatusUnauthorized)
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
				log.Printf("Failed to sing refresh token: %s", err)
				errorutils.WriteJSONError(w, "invalid refresh token", http.StatusUnauthorized)
			}
			w.Header().Set("Content-Type", "application/json")
			err = json.NewEncoder(w).Encode(RefreshResponse{AccessToken: tokenString})
			if err != nil {
				log.Printf("Failed to sing refresh token: %s", err)
				errorutils.WriteJSONError(w, "invalid refresh token", http.StatusUnauthorized)
			}
			fmt.Println("You look nice today")
		case errors.Is(err, jwt.ErrTokenMalformed):
			fmt.Println("That's not even a token")
		case errors.Is(err, jwt.ErrTokenSignatureInvalid):
			// Invalid signature
			fmt.Println("Invalid signature")
		case errors.Is(err, jwt.ErrTokenExpired) || errors.Is(err, jwt.ErrTokenNotValidYet):
			// Token is either expired or not active yet
			fmt.Println("Timing is everything")
		default:
			fmt.Println("Couldn't handle this token:", err)
		}
	} else {
		log.Printf("refreshtoken does not exist in database: %s", err)
		errorutils.WriteJSONError(w, "authentication failed", http.StatusUnauthorized)
	}
}
