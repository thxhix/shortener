package middleware

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"log"
	"net/http"
	"strings"

	"github.com/google/uuid"
)

type ctxKey string

const (
	cookieName string = "token"
	secretKey  string = "yp_iter14"
	separator  string = "."
	UserIDKey  ctxKey = "user_id"
)

func CheckAuth() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log.Println("CheckAuth START: ", r.URL.Path)
			cookie, err := r.Cookie(cookieName)
			log.Println("CheckAuth cookie: ", cookie)

			if err == nil {
				// Если токен-кука есть – пилим ее на userID и secretKey
				parts := strings.Split(cookie.Value, separator)
				if len(parts) == 2 {
					userID := parts[0]
					signature := parts[1]
					// Проверяем secretKey
					if len(parts) == 2 && verifyToken(userID, signature) {
						log.Println("CheckAuth cookie valid: ", userID)
						ctx := context.WithValue(r.Context(), UserIDKey, userID)
						// secretKey верный, идем дальше
						next.ServeHTTP(w, r.WithContext(ctx))
						return
					}
				}
			}
			// Если нет куки или она невалидна — обрываем с ошибкой
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
		})
	}
}

func SetAuth() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log.Println("SetAuth START: ", r.URL.Path)

			// ИЛИ если secretKey "левый" — создаём новый и пихаем в ответ
			cookie, err := r.Cookie(cookieName)
			log.Println("SetAuth cookie: ", cookie)
			var userID string

			if err != nil || !isValidCookie(cookie) {
				log.Println("SetAuth cookie not valid: ", cookie)
				// Куки нет или она невалидна — создаем новую
				userID = uuid.NewString()
				signature := generateToken(userID)
				token := userID + separator + signature

				http.SetCookie(w, &http.Cookie{
					Name:     cookieName,
					Value:    token,
					Path:     "/",
					HttpOnly: true,
				})
			} else {
				// Кука валидна — извлекаем userID
				parts := strings.Split(cookie.Value, separator)
				userID = parts[0]
				log.Println("SetAuth cookie valid: ", userID)
			}

			log.Println("SetAuth userID:", userID)
			ctx := context.WithValue(r.Context(), UserIDKey, userID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func isValidCookie(cookie *http.Cookie) bool {
	if cookie == nil {
		return false
	}
	parts := strings.Split(cookie.Value, separator)
	if len(parts) != 2 {
		return false
	}
	return verifyToken(parts[0], parts[1])
}

func generateToken(userID string) string {
	h := hmac.New(sha256.New, []byte(secretKey))
	h.Write([]byte(userID))
	return hex.EncodeToString(h.Sum(nil))
}

func verifyToken(userID string, signature string) bool {
	expected := generateToken(userID)
	return hmac.Equal([]byte(expected), []byte(signature))
}

func GetUserID(ctx context.Context) string {
	val := ctx.Value(UserIDKey)
	if uid, ok := val.(string); ok {
		return uid
	}
	return ""
}
