package middlewares

import (
	"net/http"
	"os"
	authstruct "storegestserver/pkg/features/auth/struct"
	"strconv"

	"github.com/golang-jwt/jwt"
)

func AuthHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var jwtKey = []byte(os.Getenv("JWTSECRET"))

		authHeader := r.Header.Get("auth")
		if authHeader == "" {
			panic(GormError{Code: http.StatusBadRequest, Message: "Auth Header missing", IsGorm: true})
		}

		tokenData := &authstruct.TokenStruct{}
		token, err := jwt.ParseWithClaims(authHeader, tokenData, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})

		if err != nil {
			if err == jwt.ErrSignatureInvalid {
				panic(GormError{Code: http.StatusUnauthorized, Message: "Stop hacking!", IsGorm: true})
			}
			panic(GormError{Code: http.StatusUnauthorized, Message: "Invalid token", IsGorm: true})
		}

		if !token.Valid {
			panic(GormError{Code: http.StatusUnauthorized, Message: "Invalid token", IsGorm: true})
		}

		w.Header().Set("userId", strconv.Itoa(tokenData.Id))

		next.ServeHTTP(w, r)
	})
}
