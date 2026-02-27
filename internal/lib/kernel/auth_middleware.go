package kernel

import (
	"context"
	"md/internal/lib/auth"
	"net/http"
	"os"

	"github.com/golang-jwt/jwt/v5"
)

func NewAuthAttrHttpMiddleware(attr string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if _, hasToken := r.Context().Value(auth.ContextKeyAuthToken).(auth.Token); !hasToken {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), auth.ContextKeyAuthAttr, attr)))
		})
	}
}

var JwtHttpMiddleware = func(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookieName := os.Getenv("APP_JWT_COOKIE_NAME")

		if cookieName == "" {
			cookieName = "i"
		}

		if cookie, err := r.Cookie(cookieName); err == nil {
			token, errParse := jwt.Parse(cookie.Value, func(token *jwt.Token) (any, error) {
				return os.Getenv("APP_JWT_KEY"), nil
			})

			if errParse == nil && token.Valid {
				if claims, ok := token.Claims.(jwt.MapClaims); ok {
					r = r.WithContext(context.WithValue(r.Context(), auth.ContextKeyAuthToken, &jwtToken{claims: claims}))
				}
			}
		}

		next.ServeHTTP(w, r)
	})
}

type jwtToken struct {
	claims jwt.MapClaims
}

func (t *jwtToken) Data() map[string]any {
	return t.claims
}
