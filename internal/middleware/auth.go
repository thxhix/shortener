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

func Auth() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// ИЛИ если secretKey "левый" — создаём новый и пихаем в ответ
			cookie, err := r.Cookie(cookieName)
			var userID string

			if err != nil || !isValidCookie(cookie) {
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
			}

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
