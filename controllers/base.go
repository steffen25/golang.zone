package controllers

import (
	"io"
	"encoding/json"
	"errors"
	"log"
	"net/http"
)

// BaseController is the common interface for all controllers
// https://github.com/sjoshi6/go-rest-api-boilerplate/blob/master/controllers/base.go
/*type ControllerServiceProvider interface {
	Get(w http.ResponseWriter, r *http.Request)
	GetAll(w http.ResponseWriter, r *http.Request)
	Post(w http.ResponseWriter, r *http.Request)
	Put(w http.ResponseWriter, r *http.Request)
	Delete(w http.ResponseWriter, r *http.Request)
}*/

type APIResponse struct {
	Success bool		`json:"success,omitempty"`
	Message string		`json:"message,omitempty"`
	Data    interface{}	`json:"data,omitempty"`
}

type APIError struct {
	Success bool		`json:"success"`
	Message string		`json:"message"`
	Status int			`json:"status"`
}

type JsonData struct {
	data map[string]interface{}
}


//var NotImplemented = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//	w.Write([]byte("Not Implemented"))
//})

// Inspiration taken from https://github.com/antonholmquist/jason/
// TODO: Move into a util package maybe? or again into some interface the basecontroller is using
func GetJSON(reader io.Reader) (*JsonData, error) {
	j := new(JsonData)
	d := json.NewDecoder(reader)
	d.UseNumber()
	err := d.Decode(&j.data)
	if err != nil {
		log.Print(err)
		return nil, err
	}

	return j, err
}

func (d *JsonData) GetString(key string) (string, error) {
	keys := d.data
	err := errors.New("Could not find key: "+key)
	if v, ok := keys[key]; ok {
		return v.(string), nil
	}

	return "", err
}

func (d *JsonData) GetInt(key string) (int, error) {
	keys := d.data
	err := errors.New("Could not find key: "+key)
	if v, ok := keys[key]; ok {
		return v.(int), nil
	}

	return -1, err
}

func NewAPIError(success bool, msg string, status int) *APIError {
	return &APIError{
		Success: success,
		Message: msg,
		Status: status,
	}
}

// TODO: Use this for both the APIResponse and APIError type
func JsonResponse(res interface{}, w http.ResponseWriter, code int) {

	json, err := json.Marshal(res)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(code)
	w.Header().Set("Content-Type", "application/json")
	w.Write(json)
}

/*
type BaseController struct {
	ControllerServiceProvider
	*database.DB
}*/
