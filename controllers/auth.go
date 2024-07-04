package controllers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	//"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
	"wagobot.com/auth"
	"wagobot.com/db"
	"wagobot.com/model"
)

// Register handles user registration.
func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	var user model.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	user.Password = string(hashedPassword)

	// Save the user to the database
	err = db.CreateUser(user.Username, user.Password, user.Email, user.FirstName, user.LastName, user.Url)
	if err != nil {
		http.Error(w, "Failed to register user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("User registered successfully"))
}

// Login handles user login.
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var credentials struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&credentials); err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	// Retrieve user from the database by username
	user, err := db.GetUserByUsername(credentials.Username)
	if err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// Compare the stored hashed password with the provided password
	if err := auth.ComparePassword(user.Password, credentials.Password); err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// Generate JWT token
	token, err := auth.GenerateToken(user.Username)
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	// Respond with the generated token
	response := map[string]string{"token": token}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func CreateToken(w http.ResponseWriter, r *http.Request) {
	// Generate a new JWT token
	token, err := auth.CreateNewToken()
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"token": "` + token + `"}`))
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	/*var req model.LogoutRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Call the Logout method
	if err := client.Logout(); err != nil {
		http.Error(w, "Failed to log out user", http.StatusInternalServerError)
		return
	}

	// Respond with success
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "success"})
	*/
}

// Register handles user registration.
func UpdateUserURLHandler(w http.ResponseWriter, r *http.Request) {
	tokenStr := r.Header.Get("Authorization")
	if tokenStr == "" {
		http.Error(w, "Authorization header missing", http.StatusUnauthorized)
		return
	}

	tokenStr = strings.TrimPrefix(tokenStr, "Bearer ")
	//ParseToken
	claims, err := auth.ParseToken(tokenStr)
	if err != nil {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}
	//fmt.Println("Check token claims:", tokenStr)
	//fmt.Println("Check token claims:", claims)

	username, ok := claims["username"].(string)
	//fmt.Println("data", username)
	if !ok {
		http.Error(w, "Invalid usernem in token", http.StatusUnauthorized)
		return
	}

	var req struct {
		Url string `json:"url"`
	}
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	if req.Url == "" {
		http.Error(w, "URL cannot be empty", http.StatusBadRequest)
		return
	}

	err = db.UpdateUserURLWebhook(username, req.Url)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "User not found", http.StatusNotFound)
		} else {
			http.Error(w, "Failed to update URL", http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("User URL updated successfully"))
}



func GetUserHandler(w http.ResponseWriter, r *http.Request) {
	tokenStr := r.Header.Get("Authorization")
	if tokenStr == "" {
		http.Error(w, "Authorization header missing", http.StatusUnauthorized)
		return
	}

	tokenStr = strings.TrimPrefix(tokenStr, "Bearer ")
	claims, err := auth.ParseToken(tokenStr)
	if err != nil {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}


	userIDStr, ok := claims["username"].(string)
	if !ok {
		http.Error(w, "Invalid user ID in token", http.StatusUnauthorized)
		return
	}

	

	user, err := db.GetUserByID(userIDStr)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "User not found", http.StatusNotFound)
		} else {
			http.Error(w, "Failed to get user data", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}