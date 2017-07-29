package controllers

import (
	"net/http"
	"encoding/json"

	"github.com/steffen25/golang.zone/repositories"
	"github.com/steffen25/golang.zone/services"
)

type AuthController struct {
	*repositories.UserRepository
}

func NewAuthController(uc *repositories.UserRepository) *AuthController {
	return &AuthController{uc}
}

func (ac *AuthController) Login(w http.ResponseWriter, r *http.Request) {

	j, err := GetJSON(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		err := NewAPIError(false, "Invalid request", http.StatusBadRequest)
		json.NewEncoder(w).Encode(err)
		return
	}
	email, err := j.GetString("email")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		err := NewAPIError(false, "Email is required", http.StatusBadRequest)
		json.NewEncoder(w).Encode(err)
		return
	}
	if ok := IsEmail(email); !ok {
		w.WriteHeader(http.StatusBadRequest)
		err := NewAPIError(false, "You must provide a valid email address", http.StatusBadRequest)
		json.NewEncoder(w).Encode(err)
		return
	}
	u, err := ac.UserRepository.FindByEmail(email)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		err := NewAPIError(false, "Could not find user", http.StatusBadRequest)
		json.NewEncoder(w).Encode(err)
		return
	}

	pw, err := j.GetString("password")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		err := NewAPIError(false, "Password is required", http.StatusBadRequest)
		json.NewEncoder(w).Encode(err)
		return
	}

	if ok := u.CheckPassword(pw); !ok {
		w.WriteHeader(http.StatusBadRequest)
		err := NewAPIError(false, "Password do not match", http.StatusBadRequest)
		json.NewEncoder(w).Encode(err)
		return
	}

	t, err := services.GenerateJWT(u)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		err := NewAPIError(false, "Something went wrong", http.StatusInternalServerError)
		json.NewEncoder(w).Encode(err)
		return
	}

	json.NewEncoder(w).Encode(APIResponse{Success: true, Message: "Login successful", Data: t})
}