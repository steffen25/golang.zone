package controllers

import (
	"github.com/steffen25/golang.zone/repositories"
	"net/http"
	"github.com/steffen25/golang.zone/models"
	"time"
	"regexp"
	"strings"
)

type PostController struct {
	repositories.PostRepository
}

func NewPostController(pr repositories.PostRepository) *PostController {
	return &PostController{pr}
}

func (pc *PostController) GetAll(w http.ResponseWriter, r *http.Request) {
	posts, err := pc.PostRepository.GetAll()
	if err != nil {
		// something went wrong
		NewAPIError(&APIError{false, "Could not fetch posts", http.StatusBadRequest}, w)
		return
	}

	NewAPIResponse(&APIResponse{Success: true, Data: posts}, w, http.StatusOK)
}

func (pc *PostController) Create (w http.ResponseWriter, r *http.Request) {
	uid := int(r.Context().Value("userId").(float64))

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


func generateSlug(title string) string {
	re := regexp.MustCompile("[^a-z0-9]+")
	return strings.Trim(re.ReplaceAllString(strings.ToLower(title), "-"), "-")
}