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

func NewAPIError(e *APIError, w http.ResponseWriter) {
	w.WriteHeader(e.Status)
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")

	err := json.NewEncoder(w).Encode(e)
	if err != nil {
		log.Println("[API ERROR]: The website encountered an unexpected error.")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

/*func NewAPIResponse(success bool, msg string, data interface{}) *APIResponse {
	return &APIResponse{
		Success: success,
		Message: msg,
		Data: data,
	}
}*/

// TODO: Use this for both the APIResponse and APIError type
func NewAPIResponse(res *APIResponse, w http.ResponseWriter, code int) {
	w.WriteHeader(code)
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")

	err := json.NewEncoder(w).Encode(res)
	if err != nil {
		log.Println("[API RESPONSE]: The website encountered an unexpected error.")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

/*
type BaseController struct {
	ControllerServiceProvider
	*database.DB
}*/
