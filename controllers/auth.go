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
	AccessToken string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func NewAuthController(uc *repositories.UserRepository) *AuthController {
	return &AuthController{uc}
}

// TODO: Create a refresh token
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

	accessToken, err := services.GenerateJWT(u)
	if err != nil {
		NewAPIError(&APIError{false, "Something went wrong", http.StatusBadRequest}, w)
		return
	}

	refreshToken, _ := services.GenerateRefreshToken(u)

	NewAPIResponse(&APIResponse{Success: true, Message: "Login successful", Data: Token{accessToken, refreshToken}}, w, http.StatusOK)
}
func (ac *AuthController) Refresh(w http.ResponseWriter, r *http.Request) {
	uid := int(r.Context().Value("userId").(float64))
	u, err := ac.UserRepository.FindById(uid)
	if err != nil {
		NewAPIError(&APIError{false, "Could not find user", http.StatusBadRequest}, w)
		return
	}

	accessToken, err := services.GenerateJWT(u)
	if err != nil {
		NewAPIError(&APIError{false, "Something went wrong", http.StatusBadRequest}, w)
		return
	}

	refreshToken, _ := services.GenerateRefreshToken(u)

	NewAPIResponse(&APIResponse{Success: true, Message: "Refresh successful", Data: Token{accessToken, refreshToken}}, w, http.StatusOK)
}