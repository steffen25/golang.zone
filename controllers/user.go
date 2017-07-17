package controllers

import (
	"net/http"
	"github.com/steffen25/golang.zone/models"
	"encoding/json"
	"time"
)

type UserController struct {
	BaseController
}

func (uc *UserController) Get(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	user := models.User{
		ID:1,
		Name: "Hej",
		Email: "rootΩ@localhost.com",
		Password: "123456",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	json.NewEncoder(w).Encode(user)
}
