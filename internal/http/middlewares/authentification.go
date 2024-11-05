package middlewares

import (
	"fmt"
	"net/http"
	"time"

	"github.com/WeisseNacht18/url-shortener/internal/generator"
	"github.com/WeisseNacht18/url-shortener/internal/storage"
	"github.com/golang-jwt/jwt/v4"
)

type Claims struct {
	jwt.RegisteredClaims
	UserID string
}

const TOKEN_EXP = time.Hour * 3
const SECRET_KEY = "28ea60d2-3126-4c96-88bf-a3505d7b6ea0"

func GetUserID(tokenString string) string {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims,
		func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
			}
			return []byte(SECRET_KEY), nil
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
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(TOKEN_EXP)),
		},
		UserID: id,
	})

	tokenString, err := token.SignedString([]byte(SECRET_KEY))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func WithAuthentification(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, err := r.Cookie("auth")

		if err == nil {
			userID := GetUserID(user.Value)

			if err == nil && storage.CheckUserIDWithToken(userID, user.Value) {
				r.Header.Set("x-user-id", userID)
				next.ServeHTTP(w, r)
				return
			}
		}

		userID, err := generator.GenerateUserID()

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}

		jwtString, err := BuildJWTString(userID)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}

		ok := storage.AddUserIDWithToken(userID, jwtString)

		if !ok {
			w.WriteHeader(http.StatusInternalServerError)
		}

		authCookie := &http.Cookie{
			Name:  "auth",
			Value: jwtString,
		}

		http.SetCookie(w, authCookie)
		w.WriteHeader(http.StatusUnauthorized)
	})
}
