package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"example.com/mygamelist/errorutils"
	"example.com/mygamelist/repository"
	"example.com/mygamelist/service"
	"example.com/mygamelist/utils"
	"github.com/golang-jwt/jwt/v5"
)

type RegisterRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func Register(env *utils.Env) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {

		if req.Header.Get("Content-Type") != "application/json" {
			http.Error(w, "Unsupported Media Type", http.StatusUnsupportedMediaType)
			return
		}

		var regReq RegisterRequest
		decoder := json.NewDecoder(req.Body)
		decoder.DisallowUnknownFields()
		if err := decoder.Decode(&regReq); err != nil {
			http.Error(w, "Invalid JSON: "+err.Error(), http.StatusBadRequest)
			return
		}

		userId, err := service.RegisterUser(env.DB, regReq.Username, regReq.Email, regReq.Password)
		if err != nil {
			switch {
			case errors.Is(err, errorutils.ErrUserExists):
				errorutils.WriteJSONError(w, "User already exists", http.StatusBadRequest)
			default:
				log.Printf("Failed to register user: %s", err)
				http.Error(w, "Error adding user: "+err.Error(), http.StatusInternalServerError)
			}
			return
		}
		type RegisterResponse struct {
			UserID int64 `json:"user_id"`
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(RegisterResponse{UserID: userId})
	}

}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func Login(env *utils.Env) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {

		var loginReq LoginRequest
		decoder := json.NewDecoder(req.Body)
		decoder.DisallowUnknownFields()

		if err := decoder.Decode(&loginReq); err != nil {
			errorutils.WriteJSONError(w, "invalid JSON", http.StatusBadRequest)
			return
		}
		jwtToken, err := service.LoginUser(env.DB, loginReq.Username, loginReq.Password)
		if err != nil {
			switch {
			case errors.Is(err, errorutils.ErrPasswordMatch):
				log.Printf("Failed to login user: %s", err)
				errorutils.WriteJSONError(w, "incorrect username or password", http.StatusUnauthorized)
			default:
				log.Printf("Failed to login user: %s", err)
				errorutils.WriteJSONError(w, "authentication failed", http.StatusUnauthorized)
			}
			return
		}

		type LoginResponse struct {
			AccessToken string `json:"accessToken"`
		}
		refreshToken, jti, err := utils.GenerateRefreshToken(loginReq.Username)
		if err != nil {
			log.Printf("Failed to login user: %s", err)
			errorutils.WriteJSONError(w, "authentication failed", http.StatusUnauthorized)
		}
		refreshTokenCookie, err := utils.CreateRefreshTokenCookie(refreshToken)
		if err != nil {
			log.Printf("Failed to login user: %s", err)
			errorutils.WriteJSONError(w, "authentication failed", http.StatusUnauthorized)
		}

		service.StoreRefreshToken(env.DB, loginReq.Username, refreshToken, jti)
		http.SetCookie(w, refreshTokenCookie)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(LoginResponse{AccessToken: jwtToken})
	}
}

type RefreshRequest struct {
	Username string `json:"username"`
}

func Refresh(env *utils.Env) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {

		var refreshReq RefreshRequest
		decoder := json.NewDecoder(req.Body)
		decoder.DisallowUnknownFields()

		if err := decoder.Decode(&refreshReq); err != nil {
			errorutils.WriteJSONError(w, "invalid JSON", http.StatusBadRequest)
			return
		}

		cookie, err := req.Cookie("refreshToken")
		if err != nil {
			http.Error(w, "missing refresh token", http.StatusUnauthorized)
			return
		}

		tokenStr := cookie.Value

		token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
			return []byte("secret-key"), nil
		}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))
		if err != nil {
			http.Error(w, "invalid refresh token", http.StatusUnauthorized)
			return
		}
		userId, err := repository.SelectUserIdByUsername(env.DB, refreshReq.Username)
		if err != nil {
			http.Error(w, "failed to retrieve userId", http.StatusUnauthorized)
		}
		_, jtifromDb, err := repository.RefreshTokenById(env.DB, userId)
		if err != nil {
			http.Error(w, "failed to retrieve refresh token", http.StatusUnauthorized)
		}
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			http.Error(w, "invalid token claims", http.StatusUnauthorized)
			return
		}

		jti, ok := claims["jti"].(string)

		if !ok {
			http.Error(w, "jti claim missing or invalid", http.StatusUnauthorized)
			return
		}
		type RefreshResponse struct {
			AccessToken string `json:"accessToken"`
		}

		if jti == jtifromDb {
			switch {
			case token.Valid:
				var secretKey = []byte("secret-key")
				expirationTime := time.Now().Add(5 * time.Minute).Unix()
				token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
					"username": refreshReq.Username,
					"exp":      expirationTime,
				})
				tokenString, err := token.SignedString(secretKey)
				if err != nil {
					http.Error(w, "failed to sign key in refresh", http.StatusBadGateway)
				}
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				json.NewEncoder(w).Encode(RefreshResponse{AccessToken: tokenString})
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
}
