package base

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"
	// "reflect"
	"time"
	"encoding/json"

	"github.com/dgrijalva/jwt-go"
	"github.com/joho/godotenv" // Import godotenv package
	"golang.org/x/crypto/bcrypt"

	"wagobot.com/model"
)

var signingKey []byte
var CurrentUser model.User

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
func GenerateToken(username string, id int) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	
	if id != 0 {
		claims["id"] = id
	}

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
			SetResponse(w, http.StatusUnauthorized, "Unauthorized")
			return
		}

		// Check if the token starts with "Bearer "
		if !strings.HasPrefix(authHeader, "Bearer ") {
			SetResponse(w, http.StatusUnauthorized, "Unauthorized")
			// http.Error(w, "Unauthorized: Invalid token format", http.StatusUnauthorized)
			return
		}

		// Remove "Bearer " prefix from the token string
		tokenString := authHeader[len("Bearer "):]

		// // Parse and validate the token
		claims, err := ParseToken(tokenString)
		if err != nil {
			SetResponse(w, http.StatusUnauthorized, "Unauthorized")
			// http.Error(w, "Unauthorized: Invalid token", http.StatusUnauthorized)
			return
		}

		CurrentUser = model.User {
			ID: int(claims["id"].(float64)),
			Username: claims["username"].(string),
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

	tokenString, err := GenerateToken("optimasi", 0)
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

	token, err := GenerateToken(secretKey, 0)
	if err != nil {
		return "", fmt.Errorf("could not create token: %w", err)
	}
	return token, nil
}

func ValidateRequest(r *http.Request, datatype interface{}) (interface{}, error) {
    contentType := r.Header.Get("Content-Type")

	if contentType == "application/json" {
        // Handle JSON body
		fmt.Println("THe body", r)
        err := json.NewDecoder(r.Body).Decode(&datatype)
		if err != nil {
			return datatype, nil
		} else {
			return datatype, err
		}
    }

    // Handle form data
    err := r.ParseForm()
    if err != nil {
        return datatype, err
    }

    // // Use reflection to map form values to struct fields
    // val := reflect.ValueOf(datatype).Elem()
    // for i := 0; i < val.NumField(); i++ {
    //     field := val.Type().Field(i)
    //     formTag := field.Tag.Get("form")
    //     if formTag != "" {
    //         formValue := r.FormValue(formTag)
    //         fieldVal := val.Field(i)
    //         if fieldVal.CanSet() {
    //             switch fieldVal.Kind() {
    //             case reflect.String:
    //                 fieldVal.SetString(formValue)
    //             case reflect.Int:
    //                 if intValue, err := strconv.Atoi(formValue); err == nil {
    //                     fieldVal.SetInt(int64(intValue))
    //                 }
    //             }
    //         }
    //     }
    // }

    return datatype, nil
}


func SetResponse(w http.ResponseWriter, statusCode int, data interface{}) {

	status := "success"

	response := map[string]interface{}{
		"status": status,
		"data": data,
	}

	if statusCode == 0 {
		statusCode = 200
		response = map[string]interface{}{
			"version": data,
		}
	} else if statusCode != 200 {
		status = "error"
		response = map[string]interface{}{
			"status": status,
			"message": data,
		}
	}

    // Marshal the data into a pretty JSON format
    jsonResponse, err := json.MarshalIndent(response, "", "")
    if err != nil {
        // If marshalling fails, respond with a 500 Internal Server Error
        http.Error(w, "Failed.", http.StatusInternalServerError)
        return
    }

    // Set the content type to application/json
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(statusCode)
    
    // Write the JSON response
    w.Write(jsonResponse)
}
