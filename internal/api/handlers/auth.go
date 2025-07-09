package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/OmerYesilkaya/fileuploader/internal/api"
	"github.com/OmerYesilkaya/fileuploader/internal/utils"
	"github.com/google/uuid"
)

type User struct {
	ID           string
	Name         string
	Email        string
	PasswordHash string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AuthHandler struct {
	Ctx *api.AppContext
}

func (h *AuthHandler) HandleLogin(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		utils.Error(w, http.StatusBadRequest, "invalid JSON:"+err.Error())
		return
	}

	// Check if user exists
	var user User
	stmt := "SELECT * FROM user WHERE email = ?"
	err = h.Ctx.DB.QueryRow(stmt, req.Email).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.PasswordHash,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			utils.Error(w, http.StatusUnauthorized, "invalid_email_or_password")
			return
		}
		utils.Error(w, http.StatusInternalServerError, "user check failed:"+err.Error())
		return
	}

	// Verify password
	hashedPassword := user.PasswordHash
	err = utils.CheckPasswordHash(req.Password, hashedPassword)
	if err != nil {
		utils.Error(w, http.StatusUnauthorized, "invalid_email_or_password")
		return
	}

	// Generate JWT token
	token, err := utils.GenerateJWT(user.ID)
	if err != nil {
		utils.Error(w, http.StatusInternalServerError, "failed to generate token:"+err.Error())
		return
	}

	response := map[string]string{"token": token, "email": user.Email, "name": user.Name, "id": user.ID}
	err = utils.ResponseSuccess(w, http.StatusOK, "login successful", response)
	if err != nil {
		utils.Error(w, http.StatusInternalServerError, "failed to encode response:"+err.Error())
	}
}

type SignupRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Name     string `json:"name"`
}

func (h *AuthHandler) HandleSignup(w http.ResponseWriter, r *http.Request) {
	var req SignupRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		utils.Error(w, http.StatusBadRequest, "invalid JSON:"+err.Error())
		return
	}

	// Check if user already exists
	var userExists bool
	stmt := "SELECT EXISTS(SELECT 1 FROM user WHERE email = ?)"
	err = h.Ctx.DB.QueryRow(stmt, req.Email).Scan(&userExists)
	if err != nil {
		utils.Error(w, http.StatusInternalServerError, "user check failed:"+err.Error())
		return
	}
	if userExists {
		utils.Error(w, http.StatusConflict, "user_already_exists")
		return
	}

	// Hash the password
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		utils.Error(w, http.StatusInternalServerError, "failed to hash password:"+err.Error())
		return
	}

	// Insert new user
	stmt = "INSERT INTO user (id, email, password_hash, name) VALUES (?, ?, ?, ?)"
	id := uuid.New()
	_, err = h.Ctx.DB.Exec(stmt, id, req.Email, hashedPassword, req.Name)
	if err != nil {
		utils.Error(w, http.StatusInternalServerError, "failed to create user: "+err.Error())
		return
	}

	response := map[string]string{"message": "user created successfully"}
	err = utils.ResponseSuccess(w, http.StatusCreated, "user created successfully", response)
	if err != nil {
		utils.Error(w, http.StatusInternalServerError, "failed to encode response:"+err.Error())
		return
	}
}
