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
	// cookieName is the name of the authentication cookie
	// used to store the signed user token.
	cookieName string = "token"

	// separator is the delimiter used to split userID and its signature
	// inside the authentication cookie value.
	separator string = "."

	// UserIDKey is the context key used to store the authenticated user ID.
	UserIDKey ctxKey = "user_id"
)

// Auth checks authorize cookie from request, or generate new if not exists.
//
// If authorized cookie is present, the user ID is extracted and placed into the request context.
// If the cookie is missing or invalid, a new user ID is generated, signed, and set as a cookie.
//
// The user ID can be retrieved later from the request context using GetUserID.
func Auth(secretKey string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// ИЛИ если secretKey "левый" — создаём новый и пихаем в ответ
			cookie, err := r.Cookie(cookieName)
			var userID string

			if err != nil || !isValidCookie(cookie, secretKey) {
				// Куки нет или она невалидна — создаем новую
				userID = uuid.NewString()
				signature := generateToken(userID, secretKey)
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

func isValidCookie(cookie *http.Cookie, secretKey string) bool {
	if cookie == nil {
		return false
	}
	parts := strings.Split(cookie.Value, separator)
	if len(parts) != 2 {
		return false
	}
	return verifyToken(parts[0], parts[1], secretKey)
}

func generateToken(userID string, secretKey string) string {
	h := hmac.New(sha256.New, []byte(secretKey))
	h.Write([]byte(userID))
	return hex.EncodeToString(h.Sum(nil))
}

func verifyToken(userID string, signature string, secretKey string) bool {
	expected := generateToken(userID, secretKey)
	return hmac.Equal([]byte(expected), []byte(signature))
}

// GetUserID extracts the user ID from the given context.
// If no user ID is found, it returns an empty string.
func GetUserID(ctx context.Context) string {
	val := ctx.Value(UserIDKey)
	if uid, ok := val.(string); ok {
		return uid
	}
	return ""
}
