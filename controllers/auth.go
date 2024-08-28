package controllers

import (
	"errors"
	"strings"
	"net/http"
	"database/sql"
	"encoding/json"

	"golang.org/x/crypto/bcrypt"
	"wagobot.com/base"
	"wagobot.com/db"
	"wagobot.com/model"
	"wagobot.com/response"
)

// Register handles user registration.
func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	var user model.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		// http.Error(w, "Bad Request", http.StatusBadRequest)
		base.SetResponse(w, http.StatusBadRequest, "Bad request")
		return
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		// http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		base.SetResponse(w, http.StatusInternalServerError, "Password error")
		return
	}
	user.Password = string(hashedPassword)

	// Save the user to the database

	err = db.CreateUser(user.Username, user.Password, user.Email, user.FirstName, user.LastName, user.Url)
	if err != nil {
		// http.Error(w, "Failed to register user", http.StatusInternalServerError)
		base.SetResponse(w, http.StatusInternalServerError, "Failed to register user")
		return
	}

	// w.WriteHeader(http.StatusOK)
	// w.Write([]byte("User registered successfully"))
	base.SetResponse(w, http.StatusOK, "User registered successfully")
}

// @Summary Login user
// @Description Logs in a user with username and password
// @Tags Auth
// @Accept json
// @Produce json
// @Param credentials body response.Credentials true "Username and Password"
// @Success 200 {object} response.TokenResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 401 {object} response.ErrorResponse
// @Router /api/login [post]
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var credentials response.Credentials

	if err := json.NewDecoder(r.Body).Decode(&credentials); err != nil {
		// http.Error(w, "Bad Request", http.StatusBadRequest)
		base.SetResponse(w, http.StatusBadRequest, "Bad Request")
		return
	}

	// Retrieve user from the database by username
	user, err := db.GetUserByUsername(credentials.Username)
	if err != nil {
		// http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		base.SetResponse(w, http.StatusUnauthorized, "Invalid credentials")
		return
	}

	// Compare the stored hashed password with the provided password
	if err := base.ComparePassword(user.Password, credentials.Password); err != nil {
		// http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		base.SetResponse(w, http.StatusUnauthorized, "Invalid credentials")
		return
	}

	// Generate JWT token
	token, err := base.GenerateToken(user.Username, user.ID)
	if err != nil {
		// http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		base.SetResponse(w, http.StatusInternalServerError, "Failed to generate token")
		return
	}

	// // Respond with the generated token
	response := response.TokenResponse{Token: token}
	// w.Header().Set("Content-Type", "application/json")
	// if err := json.NewEncoder(w).Encode(response); err != nil {
	// 	http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	// 	return
	// }

	base.SetResponse(w, http.StatusOK, response)
}

func CreateToken(w http.ResponseWriter, r *http.Request) {
	// Generate a new JWT token
	token, err := base.CreateNewToken()
	if err != nil {
		// http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		base.SetResponse(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	response := response.TokenResponse{Token: token}
	base.SetResponse(w, http.StatusOK, response)

	// w.Header().Set("Content-Type", "application/json")
	// w.WriteHeader(http.StatusOK)
	// w.Write([]byte(`{"token": "` + token + `"}`))
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

func UpdateWbhookURLHandler(w http.ResponseWriter, r *http.Request) {
	tokenStr := r.Header.Get("Authorization")
	if tokenStr == "" {
		http.Error(w, "Authorization header missing", http.StatusUnauthorized)
		return
	}

	tokenStr = strings.TrimPrefix(tokenStr, "Bearer ")
	//ParseToken
	claims, err := base.ParseToken(tokenStr)
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

