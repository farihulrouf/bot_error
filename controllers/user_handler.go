package controllers

import (
	"errors"
	"regexp"
	// "strings"
	"net/http"
	"database/sql"
	"encoding/json"
	"wagobot.com/db"
	"wagobot.com/base"
	"wagobot.com/model"
	"golang.org/x/crypto/bcrypt"
)

func GetUserHandler(w http.ResponseWriter, r *http.Request) {

	// tokenStr := r.Header.Get("Authorization")
	// if tokenStr == "" {
	// 	http.Error(w, "Authorization header missing", http.StatusUnauthorized)
	// 	return
	// }

	// tokenStr = strings.TrimPrefix(tokenStr, "Bearer ")
	// claims, err := base.ParseToken(tokenStr)
	// if err != nil {
	// 	http.Error(w, "Invalid token", http.StatusUnauthorized)
	// 	return
	// }

	// userIDStr, ok := claims["username"].(string)
	// if !ok {
	// 	http.Error(w, "Invalid user ID in token", http.StatusUnauthorized)
	// 	return
	// }

	username := base.CurrentUser.Username
	user, err := db.GetUserByID(username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// http.Error(w, "User not found", http.StatusNotFound)
			base.SetResponse(w, http.StatusNotFound, "User not found")
			return
		} else {
			base.SetResponse(w, http.StatusInternalServerError, "Failed to get user data")
			return
			// http.Error(w, "Failed to get user data", http.StatusInternalServerError)
		}
	}

	base.SetResponse(w, 200, user)
}

func UserUpdateHandler(w http.ResponseWriter, r *http.Request) {
	// // Ambil token dari header Authorization
	// tokenStr := r.Header.Get("Authorization")
	// if tokenStr == "" {
	// 	http.Error(w, "Authorization header missing", http.StatusUnauthorized)
	// 	return
	// }

	// // Hapus prefix "Bearer " dari token
	// tokenStr = strings.TrimPrefix(tokenStr, "Bearer ")
	// claims, err := base.ParseToken(tokenStr)
	// if err != nil {
	// 	http.Error(w, "Invalid token", http.StatusUnauthorized)
	// 	return
	// }

	// // Ambil username dari claims JWT
	// username, ok := claims["username"].(string)
	// if !ok {
	// 	http.Error(w, "Invalid username in token", http.StatusUnauthorized)
	// 	return
	// }

	var req model.User

	username := base.CurrentUser.Username

	// Decode body request JSON ke struct model.User
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		// http.Error(w, "Bad Request", http.StatusBadRequest)
		base.SetResponse(w, http.StatusBadRequest, "Bad Request")
		return
	}

	// Validasi: Format email
	if req.Email != "" && !isValidEmail(req.Email) {
		// http.Error(w, "Invalid email format", http.StatusBadRequest)
		base.SetResponse(w, http.StatusBadRequest, "Bad Request")
		return
	}

	// Cek apakah current password sama dengan yang tersimpan
	if req.CurrentPassword != "" && req.NewPassword != "" {
		user, err := db.GetUserByUsername(username)
		if err != nil {
			// http.Error(w, "Failed to retrieve user data", http.StatusInternalServerError)
			base.SetResponse(w, http.StatusBadRequest, "Failed to retrieve user data")
			return
		}

		err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.CurrentPassword))
		if err != nil {
			// http.Error(w, "Current password is incorrect", http.StatusUnauthorized)
			base.SetResponse(w, http.StatusUnauthorized, "Current password is incorrect")
			return
		}

		// Hash password baru
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
		if err != nil {
			// http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			base.SetResponse(w, http.StatusInternalServerError, "Generate password error")
			return
		}

		// Update profil dengan password baru yang ter-hash
		err = db.UpdateUserProfile(username, req.FirstName, req.LastName, req.Email, string(hashedPassword))
		if err != nil {
			// http.Error(w, "Failed to update profile", http.StatusInternalServerError)
			base.SetResponse(w, http.StatusInternalServerError, "Failed to update profile")
			return
		}

		// w.WriteHeader(http.StatusOK)
		// w.Write([]byte("Profile updated successfully with password"))
		base.SetResponse(w, http.StatusOK, "Profile updated successfully")
		return
	}

	// Update profil tanpa mengubah password
	err = db.UpdateUserProfile(username, req.FirstName, req.LastName, req.Email, req.NewPassword)
	if err != nil {
		// http.Error(w, "Failed to update profile", http.StatusInternalServerError)
		base.SetResponse(w, http.StatusInternalServerError, "Failed to update profile")
		return
	}

	base.SetResponse(w, http.StatusOK, "Profile updated successfully")

	// w.WriteHeader(http.StatusOK)
	// w.Write([]byte("Profile updated successfully without password"))
}

func isValidEmail(email string) bool {
	emailRegex := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	match, _ := regexp.MatchString(emailRegex, email)
	return match // Placeholder for demo purpose
}
