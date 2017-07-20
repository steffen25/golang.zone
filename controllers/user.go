package controllers

import (
	"net/http"
	"encoding/json"
	"github.com/gorilla/mux"
	"strconv"
	"regexp"
	"time"
	"fmt"

	"github.com/steffen25/golang.zone/repositories"
	"github.com/steffen25/golang.zone/models"
)

// Embed a UserDAO/Repository thingy
type UserController struct {
	*repositories.UserRepository
}

func NewUserController(uc *repositories.UserRepository) *UserController {
	return &UserController{uc}
}

func (uc *UserController) HelloWorld(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	fmt.Fprint(w, "Hello world")
}

func (uc *UserController) Create(w http.ResponseWriter, r *http.Request) {
	// check if content type is json, validate valid email
	params := make(map[string]string)
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&params)
	if err != nil {
		err := NewAPIError(false, "Invalid request")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(err)
		return
	}

	// Validate the length of the body since some users could send a big payload
	/*required := []string{"name", "email", "password"}
	if len(params) != len(required) {
		err := NewAPIError(false, "Invalid request")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(err)
		return
	}*/
	required := []string{"name", "email", "password"}
	for _, v := range required {
		_, aerr := checkParam(params, v)
		if aerr != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(aerr)
			return
		}
	}
	var name, email, pw string = params["name"], params["email"], params["password"]

	// check if the required fields is in the request that has been sent
	// TODO: Implement something like this and embed in a basecontroller https://stackoverflow.com/a/23960293/2554631
	if len(name) < 2 || len(name) > 32 {
		w.WriteHeader(http.StatusBadRequest)
		aerr := NewAPIError(false, "Name must be between 2 and 32 characters")
		json.NewEncoder(w).Encode(aerr)
		return
	}

	const email_regex = `^([\w\.\_]{2,10})@(\w{1,}).([a-z]{2,4})$`
	if m, _ := regexp.MatchString(email_regex, email); !m {
		aerr := NewAPIError(false, "You must provide a valid email address")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(aerr)
		return
	}

	exists := uc.UserRepository.Exists(email)
	if exists {
		w.WriteHeader(http.StatusBadRequest)
		aerr := NewAPIError(false, "The email address is already in use")
		json.NewEncoder(w).Encode(aerr)
		return
	}

	if len(pw) < 6 {
		w.WriteHeader(http.StatusBadRequest)
		aerr := NewAPIError(false, "Password must not be less than 6 characters")
		json.NewEncoder(w).Encode(aerr)
		return
	}

	u := &models.User{
		Name: name,
		Email: email,
		CreatedAt: time.Now(),
	}
	u.SetPassword(pw)

	err = uc.UserRepository.Create(u)
	if err != nil {
		return
	}

	defer r.Body.Close()

	json.NewEncoder(w).Encode(APIResponse{Success: true, Message: "User created"})
}

func (uc *UserController) GetAll(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	users, err := uc.UserRepository.GetAll()
	if err != nil {
		// something went wrong
		err := APIError{Success: false, Message: "Could not fetch users"}
		json.NewEncoder(w).Encode(err)
		return
	}

	data := APIResponse{Data:users}
	json.NewEncoder(w).Encode(data)
}

func (uc *UserController) GetById(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		err := APIError{Success: false, Message: "Invalid request"}
		json.NewEncoder(w).Encode(err)
		return
	}
	user, err := uc.UserRepository.FindById(id)
	if err != nil {
		// user was not found
		w.WriteHeader(http.StatusBadRequest)
		err := APIError{Success: false, Message: "Could not find user"}
		json.NewEncoder(w).Encode(err)
		return
	}
	data := APIResponse{Data:user}
	json.NewEncoder(w).Encode(data)
}

func checkParam(haystack map[string]string, needle string) (string, *APIError) {
	if n, ok := haystack[needle]; ok {
		return n, nil
	}
	err := NewAPIError(false, "Missing "+needle+" from request")
	return "", err
}

func NewAPIError(success bool, msg string) *APIError {
	return &APIError{
		Success: success,
		Message: msg,
	}
}