package controllers

import (
	"net/http"
	"time"
	"strings"

	"github.com/steffen25/golang.zone/models"
	"github.com/steffen25/golang.zone/services"
	"github.com/steffen25/golang.zone/repositories"
	"github.com/steffen25/golang.zone/app"
	"github.com/gorilla/mux"
	"strconv"
	"github.com/go-sql-driver/mysql"
)

type PostController struct {
	*app.App
	repositories.PostRepository
}

func NewPostController(a *app.App, pr repositories.PostRepository) *PostController {
	return &PostController{a, pr}
}

func (pc *PostController) GetAll(w http.ResponseWriter, r *http.Request) {
	posts, err := pc.PostRepository.GetAll()
	if err != nil {
		NewAPIError(&APIError{false, "Could not fetch posts", http.StatusBadRequest}, w)
		return
	}

	NewAPIResponse(&APIResponse{Success: true, Data: posts}, w, http.StatusOK)
}

func (pc *PostController) Create(w http.ResponseWriter, r *http.Request) {
	uid, err := services.UserIdFromContext(r.Context())
	if err != nil {
		NewAPIError(&APIError{false, "Something went wrong", http.StatusInternalServerError}, w)
		return
	}

	j, err := GetJSON(r.Body)
	if err != nil {
		NewAPIError(&APIError{false, "Invalid request", http.StatusBadRequest}, w)
		return
	}

	title, err := j.GetString("title")
	if err != nil {
		NewAPIError(&APIError{false, "Title is required", http.StatusBadRequest}, w)
		return
	}

	body, err := j.GetString("body")
	if err != nil {
		NewAPIError(&APIError{false, "Content is required", http.StatusBadRequest}, w)
		return
	}

	post := &models.Post{
		Title: title,
		Slug: generateSlug(title),
		Body: body,
		CreatedAt: time.Now(),
		UserID: uid,
	}

	err = pc.PostRepository.Create(post)
	if err != nil {
		NewAPIError(&APIError{false, "Could not create post", http.StatusBadRequest}, w)
		return
	}

	defer r.Body.Close()
	NewAPIResponse(&APIResponse{Success: true, Message: "Post created"}, w, http.StatusOK)
}

func (pc *PostController) Update(w http.ResponseWriter, r *http.Request) {
	uid, err := services.UserIdFromContext(r.Context())
	if err != nil {
		NewAPIError(&APIError{false, "Something went wrong", http.StatusInternalServerError}, w)
		return
	}
	vars := mux.Vars(r)
	postId, err := strconv.Atoi(vars["id"])
	if err != nil {
		NewAPIError(&APIError{false, "Invalid request", http.StatusBadRequest}, w)
		return
	}
	post, err := pc.PostRepository.FindById(postId)
	if err != nil {
		// post was not found
		NewAPIError(&APIError{false, "Could not find post", http.StatusNotFound}, w)
		return
	}

	j, err := GetJSON(r.Body)
	if err != nil {
		NewAPIError(&APIError{false, "Invalid request", http.StatusBadRequest}, w)
		return
	}

	title, err := j.GetString("title")
	if err != nil {
		NewAPIError(&APIError{false, "Title is required", http.StatusBadRequest}, w)
		return
	}

	body, err := j.GetString("body")
	if err != nil {
		NewAPIError(&APIError{false, "Content is required", http.StatusBadRequest}, w)
		return
	}
	if len(strings.Fields(title)) < 2 {
		NewAPIError(&APIError{false, "Title is too short", http.StatusBadRequest}, w)
		return
	}

	if len(strings.Fields(body)) < 5 {
		NewAPIError(&APIError{false, "Body is too short", http.StatusBadRequest}, w)
		return
	}
	post.UserID = uid
	post.UpdatedAt = mysql.NullTime{Time: time.Now(), Valid: true}
	post.Title = title
	post.Body = body
	post.Slug = generateSlug(title)
	err = pc.PostRepository.Update(post)
	if err != nil {
		NewAPIError(&APIError{false, "Could not update post", http.StatusBadRequest}, w)
		return
	}

	NewAPIResponse(&APIResponse{Success: true, Message: "Post updated"}, w, http.StatusOK)
}