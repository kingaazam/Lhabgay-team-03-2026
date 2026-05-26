package controllers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"Lhabgay/backend/database"
	"Lhabgay/backend/utils"

	"github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type signupRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Login handles both hardcoded admin login and normal user login.
func Login(w http.ResponseWriter, r *http.Request) {
	var req loginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.Error(w, http.StatusBadRequest, "invalid JSON request")
		return
	}

	req.Email = strings.TrimSpace(strings.ToLower(req.Email))
	if req.Email == "" || req.Password == "" {
		utils.Error(w, http.StatusBadRequest, "email and password are required")
		return
	}

	if req.Email == "admin@lhabgay.com" && req.Password == "admin123" {
		setCookie(w, "session_role", "admin", 24*time.Hour)
		utils.JSON(w, http.StatusOK, map[string]string{
			"role":     "admin",
			"redirect": "admin.html",
		})
		return
	}

	var hashedPassword string
	err := database.DB.QueryRow("SELECT password FROM users WHERE email = $1", req.Email).Scan(&hashedPassword)
	if err == sql.ErrNoRows {
		utils.Error(w, http.StatusUnauthorized, "invalid email or password")
		return
	}
	if err != nil {
		utils.Error(w, http.StatusInternalServerError, "login failed")
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(req.Password)); err != nil {
		utils.Error(w, http.StatusUnauthorized, "invalid email or password")
		return
	}

	setCookie(w, "session_role", "user", 24*time.Hour)
	setCookie(w, "user_email", req.Email, 24*time.Hour)
	utils.JSON(w, http.StatusOK, map[string]string{
		"role":     "user",
		"redirect": "home.html",
	})
}

// Signup creates a normal user account with a bcrypt password hash.
func Signup(w http.ResponseWriter, r *http.Request) {
	var req signupRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.Error(w, http.StatusBadRequest, "invalid JSON request")
		return
	}

	req.Name = strings.TrimSpace(req.Name)
	req.Email = strings.TrimSpace(strings.ToLower(req.Email))
	if req.Name == "" || req.Email == "" || req.Password == "" {
		utils.Error(w, http.StatusBadRequest, "name, email and password are required")
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		utils.Error(w, http.StatusInternalServerError, "could not hash password")
		return
	}

	_, err = database.DB.Exec(
		"INSERT INTO users (name, email, password) VALUES ($1, $2, $3)",
		req.Name,
		req.Email,
		string(hashedPassword),
	)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
			utils.Error(w, http.StatusConflict, "email already exists")
			return
		}
		utils.Error(w, http.StatusInternalServerError, "signup failed")
		return
	}

	utils.JSON(w, http.StatusCreated, map[string]string{"message": "signup successful"})
}

// Logout clears session cookies.
func Logout(w http.ResponseWriter, r *http.Request) {
	clearCookie(w, "session_role")
	clearCookie(w, "user_email")
	utils.JSON(w, http.StatusOK, map[string]string{"message": "logout successful"})
}

func setCookie(w http.ResponseWriter, name, value string, maxAge time.Duration) {
	http.SetCookie(w, &http.Cookie{
		Name:     name,
		Value:    value,
		Path:     "/",
		HttpOnly: true,
		MaxAge:   int(maxAge.Seconds()),
		SameSite: http.SameSiteLaxMode,
	})
}

func clearCookie(w http.ResponseWriter, name string) {
	http.SetCookie(w, &http.Cookie{
		Name:     name,
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		MaxAge:   -1,
		SameSite: http.SameSiteLaxMode,
	})
}
