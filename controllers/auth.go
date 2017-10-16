package controllers

import (
	"log"
	"net/http"
	"strconv"

	"github.com/steffen25/golang.zone/app"
	"github.com/steffen25/golang.zone/repositories"
	"github.com/steffen25/golang.zone/services"
)

type AuthController struct {
	App *app.App
	repositories.UserRepository
	jwtService services.JWTAuthService
}

func NewAuthController(a *app.App, us repositories.UserRepository, jwtService services.JWTAuthService) *AuthController {
	return &AuthController{a, us, jwtService}
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

	accessToken, err := ac.jwtService.GenerateAccessToken(u)
	if err != nil {
		NewAPIError(&APIError{false, "Something went wrong", http.StatusBadRequest}, w)
		return
	}

	refreshToken, err := ac.jwtService.GenerateRefreshToken(u)
	if err != nil {
		NewAPIError(&APIError{false, "Something went wrong", http.StatusBadRequest}, w)
		return
	}

	tokens := services.Tokens{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    services.TokenDuration.Seconds(),
		TokenType:    services.TokenType,
	}

	NewAPIResponse(&APIResponse{Success: true, Message: "Login successful", Data: tokens}, w, http.StatusOK)
}

func (ac *AuthController) Logout(w http.ResponseWriter, r *http.Request) {
	tokenString, err := services.GetTokenFromRequest(&ac.App.Config, r)
	if err != nil {
		NewAPIError(&APIError{false, "Something went wrong", http.StatusBadRequest}, w)
		return
	}

	jti, err := services.ExtractJti(&ac.App.Config, tokenString)
	err = ac.App.Redis.Del(jti).Err()
	if err != nil {
		NewAPIError(&APIError{false, "Something went wrong", http.StatusBadRequest}, w)
		return
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
	keys := ac.App.Redis.Keys("*" + userId + ".*")
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
	jti, err := services.ExtractRefreshTokenJti(&ac.App.Config, tokenString)
	if err != nil {
		NewAPIError(&APIError{false, "Something went wrong", http.StatusBadRequest}, w)
		return
	}
	u, err := ac.UserRepository.FindById(uid)
	if err != nil {
		NewAPIError(&APIError{false, "Could not find user", http.StatusBadRequest}, w)
		return
	}
	accessToken, err := ac.jwtService.GenerateAccessToken(u)
	if err != nil {
		NewAPIError(&APIError{false, "Something went wrong", http.StatusBadRequest}, w)
		return
	}

	refreshToken, err := ac.jwtService.GenerateRefreshToken(u)
	if err != nil {
		NewAPIError(&APIError{false, "Something went wrong", http.StatusBadRequest}, w)
		return
	}

	err = ac.App.Redis.Del(jti).Err()
	if err != nil {
		NewAPIError(&APIError{false, "Something went wrong", http.StatusBadRequest}, w)
		return
	}

	tokens := services.Tokens{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    services.TokenDuration.Seconds(),
		TokenType:    services.TokenType,
	}

	NewAPIResponse(&APIResponse{Success: true, Message: "Refresh successful", Data: tokens}, w, http.StatusOK)
}
