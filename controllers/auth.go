package controllers

import (
	"net/http"
	"encoding/json"
	"golang.org/x/crypto/bcrypt"
	"github.com/steffen25/golang.zone/repositories"
)

type AuthController struct {
	*repositories.UserRepository
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
		err := NewAPIError(false, "Email required", http.StatusBadRequest)
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
	err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(pw))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		err := NewAPIError(false, "Password do not match", http.StatusBadRequest)
		json.NewEncoder(w).Encode(err)
		return
	}

	// Login success exchange JWT
	// TODO: JWT token generation
	json.NewEncoder(w).Encode(APIResponse{Success: true, Message: "Login successful"})

}