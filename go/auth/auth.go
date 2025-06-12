package auth

import (
	"net/http"
	"os"
	"regexp"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"github.com/google/uuid"
)

var usernameRegex = regexp.MustCompile(`^[a-zA-Z0-9_-]{3,15}$`)

func Login() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var secret = os.Getenv("JWT_SECRET")
		if r.Method != http.MethodPost {
			er := http.StatusMethodNotAllowed
			http.Error(w, "Invalid method", er)
			return
		}

		username := r.FormValue("username")

		if !usernameRegex.MatchString(username) {
			er := http.StatusNotAcceptable
			http.Error(w, "Username Length Invalid Or You Used Special Characters", er)
			return

		}
		userID := uuid.NewString()

		claims := jwt.MapClaims{
			"UserName": username,
			"UserId":   userID,
			"exp":      time.Now().Add(24 * time.Hour).Unix(),
			"iat":      time.Now().Unix(),
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

		signed, err := token.SignedString([]byte(secret))
		if err != nil {
			http.Error(w, "Failed To Login", http.StatusInternalServerError)
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:     "auth_token",
			Value:    signed,
			Expires:  time.Now().Add(24 * time.Hour),
			SameSite: http.SameSiteNoneMode, // For Development, Prod will change it back to lax
			HttpOnly: true,
			Path:     "/",
			Secure:   true,
		})

		w.WriteHeader(http.StatusCreated)
		return

	}

}

// This function is basically just checking that we have that cookie so i can render the login page if i dont

// Proper validation will happen in middleware which will be connected to routes that need it
func IsLoggedIn() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		token, err := r.Cookie("auth_token")
		println(token)

		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		return

	}
}
