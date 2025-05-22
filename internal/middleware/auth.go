package middleware

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
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
			cookie, err := r.Cookie(cookieName)

			if err == nil {
				// Если токен-кука есть – пилим ее на userID и secretKey
				parts := strings.Split(cookie.Value, separator)
				userID := parts[0]
				signature := parts[1]
				// Проверяем secretKey
				if len(parts) == 2 && verifyToken(userID, signature) {
					ctx := context.WithValue(r.Context(), UserIDKey, userID)
					// secretKey верный, идем дальше
					next.ServeHTTP(w, r.WithContext(ctx))
					return
				}
			}
			// Если нет куки или она невалидна — обрываем с ошибкой
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			return
		})
	}
}

func SetAuth() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// ИЛИ если secretKey "левый" — создаём новый и пихаем в ответ
			if GetUserID(r.Context()) == "" {
				userID := uuid.NewString()
				signature := generateToken(userID)
				token := userID + separator + signature

				http.SetCookie(w, &http.Cookie{
					Name:  cookieName,
					Value: token,
					Path:  "/",
				})

				ctx := context.WithValue(r.Context(), UserIDKey, userID)
				next.ServeHTTP(w, r.WithContext(ctx))
				return
			}
			next.ServeHTTP(w, r)
		})
	}
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
