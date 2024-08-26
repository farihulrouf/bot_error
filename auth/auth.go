package auth

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/joho/godotenv" // Import godotenv package
	"golang.org/x/crypto/bcrypt"
)

var signingKey []byte

func init() {
	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		fmt.Println("Error loading .env file:", err)
	}

	// Read secret key from environment variable
	secretKey := os.Getenv("SECRET_KEY")
	if secretKey == "" {
		fmt.Println("SECRET_KEY environment variable is not set")
		signingKey = []byte("default_secret_key") // Set a default secret key
	} else {
		signingKey = []byte(secretKey)
	}
}

// GenerateToken generates a JWT token for the given username.
func GenerateToken(username string) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["username"] = username
	claims["exp"] = time.Now().AddDate(1, 0, 0).Unix() // expiration
	//claims["exp"] = time.Now().Add(time.Hour * 24).Unix() // Token expires in 24 hours

	// Generate encoded token and return it
	tokenString, err := token.SignedString(signingKey)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

// ParseToken parses the JWT token and returns the claims if valid.
func ParseToken(tokenString string) (jwt.MapClaims, error) {
	//fmt.Println("data token", tokenString)
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Validate the signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return signingKey, nil
	})
	if err != nil {
		return nil, err
	}

	// Check if the token is valid
	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	// Extract claims
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("invalid claims")
	}

	return claims, nil
}

// JWTMiddleware is a middleware to protect routes using JWT.
func JWTMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Extract token from the request header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Check if the token starts with "Bearer "
		if !strings.HasPrefix(authHeader, "Bearer ") {
			http.Error(w, "Unauthorized: Invalid token format", http.StatusUnauthorized)
			return
		}

		// Remove "Bearer " prefix from the token string
		tokenString := authHeader[len("Bearer "):]

		// Parse and validate the token
		claims, err := ParseToken(tokenString)
		if err != nil {
			http.Error(w, "Unauthorized: Invalid token", http.StatusUnauthorized)
			return
		}

		// Pass the claims to the next handler
		ctx := context.WithValue(r.Context(), "claims", claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func ComparePassword(hashedPassword, password string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		return errors.New("password does not match")
	}
	return nil
}

func CreateNewTokenFromExpired(expiredTokenString string) (string, error) {
	// Parse the expired token to get the claims
	claims, err := parseTokenWithoutValidation(expiredTokenString)
	if err != nil {
		return "", fmt.Errorf("could not parse expired token: %w", err)
	}

	claims["exp"] = time.Now().AddDate(1, 0, 0).Unix() // Token expires in one year

	tokenString, err := GenerateToken("optimasi")
	if err != nil {
		return "", fmt.Errorf("could not create new token: %w", err)
	}

	return tokenString, nil
}

// ParseTokenWithoutValidation parses the JWT token without validating its expiration.
func parseTokenWithoutValidation(tokenString string) (jwt.MapClaims, error) {
	token, _, err := new(jwt.Parser).ParseUnverified(tokenString, jwt.MapClaims{})
	if err != nil {
		return nil, fmt.Errorf("could not parse token: %w", err)
	}

	// Extract claims
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("invalid claims")
	}

	return claims, nil
}

func CreateNewToken() (string, error) {
	secretKey := os.Getenv("SECRET_KEY")
	if secretKey == "" {
		return "", fmt.Errorf("SECRET_KEY is not set in .env file")
	}

	token, err := GenerateToken(secretKey)
	if err != nil {
		return "", fmt.Errorf("could not create token: %w", err)
	}
	return token, nil
}
