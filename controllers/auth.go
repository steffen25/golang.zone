package controllers

import (
	"net/http"

	"github.com/steffen25/golang.zone/repositories"
	"github.com/steffen25/golang.zone/services"
)

type AuthController struct {
	*repositories.UserRepository
}

type Token struct {
	Token string `json:"token"`
}

func NewAuthController(uc *repositories.UserRepository) *AuthController {
	return &AuthController{uc}
}

func (ac *AuthController) Authenticate(w http.ResponseWriter, r *http.Request) {

	j, err := GetJSON(r.Body)
	if err != nil {
		NewAPIError(&APIError{false, "Invalid request", http.StatusBadRequest}, w)
		return
	}
	email, err := j.GetString("email")
	if err != nil {
		NewAPIError(&APIError{false, "Email is required", http.StatusBadRequest}, w)
		return
	}
	if ok := IsEmail(email); !ok {
		NewAPIError(&APIError{false, "You must provide a valid email address", http.StatusBadRequest}, w)
		return
	}
	u, err := ac.UserRepository.FindByEmail(email)
	if err != nil {
		NewAPIError(&APIError{false, "Incorrect email or password", http.StatusBadRequest}, w)
		return
	}

	pw, err := j.GetString("password")
	if err != nil {
		NewAPIError(&APIError{false, "Password is required", http.StatusBadRequest}, w)
		return
	}

	if ok := u.CheckPassword(pw); !ok {
		NewAPIError(&APIError{false, "Incorrect email or password", http.StatusBadRequest}, w)
		return
	}

	t, err := services.GenerateJWT(u)
	if err != nil {
		NewAPIError(&APIError{false, "Something went wrong", http.StatusBadRequest}, w)
		return
	}

	NewAPIResponse(&APIResponse{Success: true, Message: "Login successful", Data: Token{t}}, w, http.StatusOK)
}