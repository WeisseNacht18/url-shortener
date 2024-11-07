package middlewares

import (
	"fmt"
	"net/http"
	"time"

	"github.com/WeisseNacht18/url-shortener/internal/generator"
	"github.com/WeisseNacht18/url-shortener/internal/logger"
	"github.com/WeisseNacht18/url-shortener/internal/storage"
	"github.com/golang-jwt/jwt/v4"
)

type Claims struct {
	jwt.RegisteredClaims
	UserID string
}

const TokenExp = time.Hour * 3
const SecretKey = "28ea60d2-3126-4c96-88bf-a3505d7b6ea0"

func GetUserID(tokenString string) string {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims,
		func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
			}
			return []byte(SecretKey), nil
		})
	if err != nil {
		return ""
	}

	if !token.Valid {
		return ""
	}

	return claims.UserID
}

func BuildJWTString(id string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(TokenExp)),
		},
		UserID: id,
	})

	tokenString, err := token.SignedString([]byte(SecretKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func WithAuthentification(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, err := r.Cookie("auth")

		userID := ""

		if err == nil {
			userID = GetUserID(user.Value)

			if storage.CheckUserIDWithToken(userID, user.Value) {
				r.Header.Set("x-user-id", userID)
				next.ServeHTTP(w, r)
				return
			}
		}

		if userID == "" {
			userID, err = generator.GenerateUserID()

			if err != nil {
				logger.Logger.Infoln(err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}

		jwtString, err := BuildJWTString(userID)

		if err != nil {
			logger.Logger.Infoln(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		ok := storage.AddUserIDWithToken(userID, jwtString)

		if !ok {
			http.Error(w, "user couldn't add", http.StatusInternalServerError)
			return
		}

		authCookie := &http.Cookie{
			Name:  "auth",
			Value: jwtString,
		}

		http.SetCookie(w, authCookie)

		if r.RequestURI != "/api/user/urls" {
			next.ServeHTTP(w, r)
			r.Header.Set("x-user-id", userID)
			return
		}

		w.WriteHeader(http.StatusUnauthorized)
	})
}
