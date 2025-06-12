package auth

import (
	"net/http"
	"os"

	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var secret = os.Getenv("JWT_SECRET")
		auth_token, err := r.Cookie("auth_token")
		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		_, err = jwt.Parse(auth_token.Value, func(token *jwt.Token) (interface{}, error) {
			return secret, nil
		}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))
		if err != nil {
			http.Error(w, "Invalid Token", http.StatusUnauthorized)

			return
		}

		// could just pass the decoded token in context but decoding it isnt that expensive or anything

		// i suppose it would fall under the DRY principle though
		next(w, r)

	}

}

//todo rate limiter middleware gonna use token bucket
