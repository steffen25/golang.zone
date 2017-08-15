package controllers

import (
	"net/http"
	"github.com/gorilla/mux"
	"strconv"
	"time"
	"fmt"
	"log"

	"github.com/steffen25/golang.zone/repositories"
	"github.com/steffen25/golang.zone/models"
	"github.com/steffen25/golang.zone/services"
)

// Embed a UserDAO/Repository thingy
type UserController struct {
	repositories.UserRepository
}

func NewUserController(us repositories.UserRepository) *UserController {
	return &UserController{us}
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

	NewAPIResponse(&APIResponse{Data:uid}, w, http.StatusOK)
}

func (uc *UserController) Create(w http.ResponseWriter, r *http.Request) {
	// Validate the length of the body since some users could send a big payload
	/*required := []string{"name", "email", "password"}
	if len(params) != len(required) {
		err := NewAPIError(false, "Invalid request")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(err)
		return
	}*/

	j, err := GetJSON(r.Body)
	if err != nil {
		NewAPIError(&APIError{false, "Invalid request", http.StatusBadRequest}, w)
		return
	}

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
	if ok := IsEmail(email); !ok {
		NewAPIError(&APIError{false, "You must provide a valid email address", http.StatusBadRequest}, w)
		return
	}
	exists := uc.UserRepository.Exists(email)
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
		Name: name,
		Email: email,
		Admin: false,
		CreatedAt: time.Now(),
	}
	u.SetPassword(pw)

	err = uc.UserRepository.Create(u)
	if err != nil {
		NewAPIError(&APIError{false, "Could not create user", http.StatusBadRequest}, w)
		return
	}

	defer r.Body.Close()
	NewAPIResponse(&APIResponse{Success: true, Message: "User created"}, w, http.StatusOK)
}

func (uc *UserController) GetAll(w http.ResponseWriter, r *http.Request) {
	users, err := uc.UserRepository.GetAll()
	if err != nil {
		// something went wrong
		NewAPIError(&APIError{false, "Could not fetch users", http.StatusBadRequest}, w)
		return
	}

	NewAPIResponse(&APIResponse{Success: true, Data:users}, w, http.StatusOK)
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

	NewAPIResponse(&APIResponse{Success: true, Data:user}, w, http.StatusOK)
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

	err = uc.UserRepository.Update(user)
	if err != nil {
		NewAPIError(&APIError{false, "Could not update user", http.StatusBadRequest}, w)
		return
	}

	NewAPIResponse(&APIResponse{Success: true, Data: user}, w, http.StatusOK)
}