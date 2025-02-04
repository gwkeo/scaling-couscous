package auth

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	"github.com/gwkeo/scaling-couscous/utils"
	"net/http"
)

func AuthenticationMW(secret string) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			tokenString := r.Header.Get("X-Auth-Token")
			if tokenString == "" {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			token, err := utils.ParseToken(tokenString, secret)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			if !token.Valid {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			claims := token.Claims.(jwt.MapClaims)

			context.Set(r, "id", claims["id"])
			context.Set(r, "role", claims["role"])
			next.ServeHTTP(w, r)
		})
	}
}
