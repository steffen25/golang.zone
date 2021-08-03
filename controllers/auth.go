package controllers

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/steffen25/golang.zone/app"
	"github.com/steffen25/golang.zone/models"
	"github.com/steffen25/golang.zone/repositories"
	"github.com/steffen25/golang.zone/services"
	"github.com/steffen25/golang.zone/util"
)

type AuthController struct {
	App *app.App
	repositories.UserRepository
	jwtService services.JWTAuthService
}

func NewAuthController(a *app.App, us repositories.UserRepository, jwtService services.JWTAuthService) *AuthController {
	return &AuthController{a, us, jwtService}
}

func (ac *AuthController) Register(w http.ResponseWriter, r *http.Request) {
	j, err := GetJSON(r.Body)
	if err != nil {
		NewAPIError(&APIError{false, "Invalid request", http.StatusBadRequest}, w)
		return
	}
	defer r.Body.Close()

	name, err := j.GetString("name")
	if err != nil {
		NewAPIError(&APIError{false, "Name is required", http.StatusBadRequest}, w)
		return
	}
	// TODO: Implement something like this and embed in a basecontroller https://stackoverflow.com/a/23960293/2554631
	if len(name) < 2 || len(name) > 32 {
		NewAPIError(&APIError{false, "Name must be between 2 and 32 characters", http.StatusBadRequest}, w)
		return
	}

	email, err := j.GetString("email")
	if err != nil {
		NewAPIError(&APIError{false, "Email is required", http.StatusBadRequest}, w)
		return
	}
	if ok := util.IsEmail(email); !ok {
		NewAPIError(&APIError{false, "You must provide a valid email address", http.StatusBadRequest}, w)
		return
	}
	exists := ac.UserRepository.Exists(email)
	if exists {
		NewAPIError(&APIError{false, "The email address is already in use", http.StatusBadRequest}, w)
		return
	}
	pw, err := j.GetString("password")
	if err != nil {
		NewAPIError(&APIError{false, "Password is required", http.StatusBadRequest}, w)
		return
	}
	if len(pw) < 6 {
		NewAPIError(&APIError{false, "Password must not be less than 6 characters", http.StatusBadRequest}, w)
		return
	}

	u := &models.User{
		Name:      name,
		Email:     email,
		Admin:     false,
		CreatedAt: time.Now(),
	}
	u.SetPassword(pw)

	err = ac.UserRepository.Create(u)
	if err != nil {
		NewAPIError(&APIError{false, "Could not create user", http.StatusBadRequest}, w)
		return
	}

	NewAPIResponse(&APIResponse{Success: true, Message: "User created"}, w, http.StatusOK)
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
	if ok := util.IsEmail(email); !ok {
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

	tokens, err := ac.jwtService.GenerateTokens(u)
	if err != nil {
		NewAPIError(&APIError{false, "Something went wrong", http.StatusBadRequest}, w)
		return
	}

	authUser := &models.AuthUser{
		User:  u,
		Admin: u.Admin,
	}

	data := struct {
		Tokens *services.Tokens `json:"tokens"`
		User   *models.AuthUser `json:"user"`
	}{
		tokens,
		authUser,
	}

	NewAPIResponse(&APIResponse{Success: true, Message: "Login successful", Data: data}, w, http.StatusOK)
}

func (ac *AuthController) Me(w http.ResponseWriter, r *http.Request) {
	uid, err := services.UserIdFromContext(r.Context())
	if err != nil {
		log.Println(err)
		NewAPIError(&APIError{false, "Something went wrong", http.StatusInternalServerError}, w)
		return
	}

	user, err := ac.UserRepository.FindById(uid)
	if err != nil {
		NewAPIError(&APIError{false, "Could not find user", http.StatusNotFound}, w)
		return
	}

	authUser := &models.AuthUser{
		User:  user,
		Admin: user.Admin,
	}

	NewAPIResponse(&APIResponse{Success: true, Message: "My information", Data: authUser}, w, http.StatusOK)
}

func (ac *AuthController) Logout(w http.ResponseWriter, r *http.Request) {
	tokenString, err := services.GetTokenFromRequest(&ac.App.Config, r)
	if err != nil {
		NewAPIError(&APIError{false, "Something went wrong", http.StatusBadRequest}, w)
		return
	}

	tokenHash, err := services.ExtractTokenHash(&ac.App.Config, tokenString)
	if err != nil {
		NewAPIError(&APIError{false, "Something went wrong", http.StatusBadRequest}, w)
		return
	}

	/*jti, err := services.ExtractJti(&ac.App.Config, tokenString)
	if err != nil {
		NewAPIError(&APIError{false, "Something went wrong", http.StatusBadRequest}, w)
		return
	}*/

	keys := ac.App.Redis.Keys("*" + tokenHash + ".*")
	for _, token := range keys.Val() {
		err := ac.App.Redis.Del(token).Err()
		if err != nil {
			log.Printf("Could not delete token: %s ; error: %v", token, err)
			NewAPIError(&APIError{false, "Something went wrong", http.StatusBadRequest}, w)
			return
		}
	}

	NewAPIResponse(&APIResponse{Success: true, Message: "Logout successful"}, w, http.StatusOK)

}

func (ac *AuthController) LogoutAll(w http.ResponseWriter, r *http.Request) {
	uid, err := services.UserIdFromContext(r.Context())
	if err != nil {
		NewAPIError(&APIError{false, "Something went wrong", http.StatusInternalServerError}, w)
		return
	}
	userId := strconv.Itoa(uid)
	keys := ac.App.Redis.Keys("*." + userId + ".*")
	for _, token := range keys.Val() {
		err := ac.App.Redis.Del(token).Err()
		if err != nil {
			log.Printf("Could not delete token: %s ; error: %v", token, err)
		}
	}

	NewAPIResponse(&APIResponse{Success: true, Message: "Logout successful"}, w, http.StatusOK)
}

func (ac *AuthController) RefreshTokens(w http.ResponseWriter, r *http.Request) {
	tokenString, err := services.GetRefreshTokenFromRequest(&ac.App.Config, r)
	if err != nil {
		NewAPIError(&APIError{false, "Something went wrong", http.StatusInternalServerError}, w)
		return
	}
	uid, err := services.UserIdFromContext(r.Context())
	if err != nil {
		NewAPIError(&APIError{false, "Something went wrong", http.StatusInternalServerError}, w)
		return
	}
	tokenHash, err := services.ExtractRefreshTokenHash(&ac.App.Config, tokenString)
	if err != nil {
		NewAPIError(&APIError{false, "Something went wrong", http.StatusBadRequest}, w)
		return
	}
	u, err := ac.UserRepository.FindById(uid)
	if err != nil {
		NewAPIError(&APIError{false, "Could not find user", http.StatusBadRequest}, w)
		return
	}
	tokens, err := ac.jwtService.GenerateTokens(u)
	if err != nil {
		NewAPIError(&APIError{false, "Something went wrong", http.StatusBadRequest}, w)
		return
	}

	keys := ac.App.Redis.Keys("*" + tokenHash + ".*")
	for _, token := range keys.Val() {
		err := ac.App.Redis.Del(token).Err()
		if err != nil {
			log.Printf("Could not delete token: %s ; error: %v", token, err)
			NewAPIError(&APIError{false, "Something went wrong", http.StatusBadRequest}, w)
			return
		}
	}

	NewAPIResponse(&APIResponse{Success: true, Message: "Refresh successful", Data: tokens}, w, http.StatusOK)
}
