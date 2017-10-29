package controllers

import (
	"net/http"
	"time"

	"fmt"
	"strconv"

	"github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/steffen25/golang.zone/app"
	"github.com/steffen25/golang.zone/models"
	"github.com/steffen25/golang.zone/repositories"
	"github.com/steffen25/golang.zone/services"
	"github.com/steffen25/golang.zone/util"
)

type PostController struct {
	*app.App
	repositories.PostRepository
	repositories.UserRepository
}

type PostPaginator struct {
	Total        int     `json:"total"`
	PerPage      int     `json:"perPage"`
	CurrentPage  int     `json:"currentPage"`
	LastPage     int     `json:"lastPage"`
	From         int     `json:"from"`
	To           int     `json:"to"`
	FirstPageUrl string  `json:"firstPageUrl"`
	LastPageUrl  string  `json:"lastPageUrl"`
	NextPageUrl  *string `json:"nextPageUrl"`
	PrevPageUrl  *string `json:"prevPageUrl"`
}

func NewPostController(a *app.App, pr repositories.PostRepository, ur repositories.UserRepository) *PostController {
	return &PostController{a, pr, ur}
}

func (pc *PostController) GetAll(w http.ResponseWriter, r *http.Request) {
	httpScheme := "https://"
	total, _ := pc.PostRepository.GetTotalPostCount()
	page := r.URL.Query().Get("page")
	pageInt, err := strconv.Atoi(page)
	if err != nil {
		pageInt = 1
	}
	perPage := r.URL.Query().Get("perpage")
	perPageInt, err := strconv.Atoi(perPage)
	if err != nil || perPageInt < 1 || perPageInt > 100 {
		perPageInt = 10
	}
	offset := (pageInt - 1) * perPageInt
	to := pageInt * perPageInt
	if to > total {
		to = total
	}

	from := offset + 1
	totalPages := (total-1)/perPageInt + 1
	prevPage := pageInt - 1
	firstPageUrl := fmt.Sprintf(httpScheme+r.Host+r.URL.Path+"?page=%d", 1)
	lastPageString := fmt.Sprintf(httpScheme+r.Host+r.URL.Path+"?page=%d", totalPages)
	var prevPageUrl string
	var nextPageUrl string
	if prevPage > 0 && prevPage < totalPages {
		prevPageUrl = fmt.Sprintf(httpScheme+r.Host+r.URL.Path+"?page=%d", prevPage)
	}

	nextPage := pageInt + 1
	if nextPage <= totalPages {
		nextPageUrl = fmt.Sprintf(httpScheme+r.Host+r.URL.Path+"?page=%d", nextPage)
	}

	posts, err := pc.PostRepository.Paginate(perPageInt, offset)
	if err != nil {
		NewAPIError(&APIError{false, "Could not fetch posts", http.StatusBadRequest}, w)
		return
	}

	if len(posts) == 0 {
		NewAPIResponse(&APIResponse{Success: false, Message: "Could not find posts", Data: posts}, w, http.StatusNotFound)
		return
	}

	postPaginator := APIPagination{
		total,
		perPageInt,
		pageInt,
		totalPages,
		from,
		to,
		firstPageUrl,
		lastPageString,
		nextPageUrl,
		prevPageUrl,
	}

	NewAPIResponse(&APIResponse{Success: true, Data: posts, Pagination: &postPaginator}, w, http.StatusOK)
}

func (pc *PostController) GetById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		NewAPIError(&APIError{false, "Invalid request", http.StatusBadRequest}, w)
		return
	}
	post, err := pc.PostRepository.FindById(id)
	if err != nil {
		NewAPIError(&APIError{false, "Could not find post", http.StatusNotFound}, w)
		return
	}

	NewAPIResponse(&APIResponse{Success: true, Data: post}, w, http.StatusOK)
}

func (pc *PostController) GetBySlug(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	slug := vars["slug"]
	post, err := pc.PostRepository.FindBySlug(slug)
	if err != nil {
		NewAPIError(&APIError{false, "Could not find post", http.StatusNotFound}, w)
		return
	}

	NewAPIResponse(&APIResponse{Success: true, Data: post}, w, http.StatusOK)
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

	title = util.CleanZalgoText(title)

	if len(title) < 10 {
		NewAPIError(&APIError{false, "Title is too short", http.StatusBadRequest}, w)
		return
	}

	slug := util.GenerateSlug(title)
	if len(slug) == 0 {
		NewAPIError(&APIError{false, "Title is invalid", http.StatusBadRequest}, w)
		return
	}

	body, err := j.GetString("body")
	if err != nil {
		NewAPIError(&APIError{false, "Content is required", http.StatusBadRequest}, w)
		return
	}

	body = util.CleanZalgoText(body)

	if len(body) < 10 {
		NewAPIError(&APIError{false, "Body is too short", http.StatusBadRequest}, w)
		return
	}

	post := &models.Post{
		Title:     title,
		Slug:      slug,
		Body:      body,
		CreatedAt: time.Now(),
		UserID:    uid,
	}

	err = pc.PostRepository.Create(post)
	if err != nil {
		NewAPIError(&APIError{false, "Could not create post", http.StatusBadRequest}, w)
		return
	}

	// TODO: Change this maybe put the user object into a context and get the author from there.
	u, err := pc.UserRepository.FindById(uid)
	if err != nil {
		NewAPIError(&APIError{false, "Content is required", http.StatusBadRequest}, w)
		return
	}
	post.Author = u.Name

	defer r.Body.Close()
	NewAPIResponse(&APIResponse{Success: true, Message: "Post created", Data: post}, w, http.StatusOK)
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

	title = util.CleanZalgoText(title)

	if len(title) < 10 {
		NewAPIError(&APIError{false, "Title is too short", http.StatusBadRequest}, w)
		return
	}

	slug := util.GenerateSlug(title)
	if len(slug) == 0 {
		NewAPIError(&APIError{false, "Title is invalid", http.StatusBadRequest}, w)
		return
	}

	body, err := j.GetString("body")
	if err != nil {
		NewAPIError(&APIError{false, "Content is required", http.StatusBadRequest}, w)
		return
	}

	body = util.CleanZalgoText(body)

	if len(body) < 10 {
		NewAPIError(&APIError{false, "Body is too short", http.StatusBadRequest}, w)
		return
	}

	post.UserID = uid
	post.UpdatedAt = mysql.NullTime{Time: time.Now(), Valid: true}
	post.Title = title
	post.Body = body
	post.Slug = slug
	err = pc.PostRepository.Update(post)
	if err != nil {
		NewAPIError(&APIError{false, "Could not update post", http.StatusBadRequest}, w)
		return
	}

	NewAPIResponse(&APIResponse{Success: true, Message: "Post updated", Data: post}, w, http.StatusOK)
}
