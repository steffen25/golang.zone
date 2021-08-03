package controllers

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"

	"github.com/go-sql-driver/mysql"
	"github.com/steffen25/golang.zone/app"
	"github.com/steffen25/golang.zone/models"
	"github.com/steffen25/golang.zone/repositories"
	"github.com/steffen25/golang.zone/services"
)

// Embed a UserDAO/Repository thingy
type UserController struct {
	*app.App
	repositories.UserRepository
	repositories.PostRepository
}

func NewUserController(a *app.App, ur repositories.UserRepository, pr repositories.PostRepository) *UserController {
	return &UserController{a, ur, pr}
}

func (uc *UserController) HelloWorld(w http.ResponseWriter, r *http.Request) {
	log.Println(r.Context().Value("userId"))
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	fmt.Fprint(w, "Hello gopher!")
}

func (uc *UserController) Profile(w http.ResponseWriter, r *http.Request) {
	uid, err := services.UserIdFromContext(r.Context())
	if err != nil {
		NewAPIError(&APIError{false, "Something went wrong", http.StatusInternalServerError}, w)
		return
	}

	NewAPIResponse(&APIResponse{Data: uid}, w, http.StatusOK)
}

func (uc *UserController) GetAll(w http.ResponseWriter, r *http.Request) {
	users, err := uc.UserRepository.GetAll()
	if err != nil {
		// something went wrong
		NewAPIError(&APIError{false, "Could not fetch users", http.StatusBadRequest}, w)
		return
	}

	NewAPIResponse(&APIResponse{Success: true, Data: users}, w, http.StatusOK)
}

func (uc *UserController) GetById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		NewAPIError(&APIError{false, "Invalid request", http.StatusBadRequest}, w)
		return
	}
	user, err := uc.UserRepository.FindById(id)
	if err != nil {
		// user was not found
		NewAPIError(&APIError{false, "Could not find user", http.StatusNotFound}, w)
		return
	}

	NewAPIResponse(&APIResponse{Success: true, Data: user}, w, http.StatusOK)
}

func (uc *UserController) Update(w http.ResponseWriter, r *http.Request) {
	uid, err := services.UserIdFromContext(r.Context())
	if err != nil {
		NewAPIError(&APIError{false, "Something went wrong", http.StatusInternalServerError}, w)
		return
	}

	user, err := uc.UserRepository.FindById(uid)
	if err != nil {
		NewAPIError(&APIError{false, "Could not find user", http.StatusBadRequest}, w)
		return
	}

	j, err := GetJSON(r.Body)
	if err != nil {
		NewAPIError(&APIError{false, "Invalid request", http.StatusBadRequest}, w)
		return
	}

	name, err := j.GetString("name")
	if name != "" && err == nil {
		user.Name = name
	}

	newpw, err := j.GetString("newpassword")
	if newpw != "" && err == nil {
		// confirm password
		oldpw, err := j.GetString("oldpassword")
		if err != nil {
			NewAPIError(&APIError{false, "Old password is required", http.StatusBadRequest}, w)
			return
		}
		ok := user.CheckPassword(oldpw)
		if !ok {
			NewAPIError(&APIError{false, "Old password do not match", http.StatusBadRequest}, w)
			return
		}
		if len(newpw) < 6 {
			NewAPIError(&APIError{false, "Password must not be less than 6 characters", http.StatusBadRequest}, w)
			return
		}
		user.SetPassword(newpw)
	}

	user.UpdatedAt = mysql.NullTime{Time: time.Now(), Valid: true}

	err = uc.UserRepository.Update(user)
	if err != nil {
		NewAPIError(&APIError{false, "Could not update user", http.StatusBadRequest}, w)
		return
	}

	authUser := &models.AuthUser{
		User:  user,
		Admin: user.Admin,
	}

	NewAPIResponse(&APIResponse{Success: true, Data: authUser}, w, http.StatusOK)
}

func (uc *UserController) FindPostsByUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		NewAPIError(&APIError{false, "Invalid request", http.StatusBadRequest}, w)
		return
	}
	user, err := uc.UserRepository.FindById(id)
	if err != nil {
		// user was not found
		NewAPIError(&APIError{false, "Could not find user", http.StatusNotFound}, w)
		return
	}

	posts, err := uc.PostRepository.FindByUser(user)
	if err != nil {
		// user was not found
		NewAPIError(&APIError{false, "Could not find user posts", http.StatusNotFound}, w)
		return
	}

	NewAPIResponse(&APIResponse{Success: true, Data: posts}, w, http.StatusOK)
}
