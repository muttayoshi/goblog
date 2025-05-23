package handlers

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/muttayoshi/goblog/database"
	"github.com/muttayoshi/goblog/models"
	"net/http"
	"strconv"
	"time"
)

// GetPosts godoc
// @Summary Get all posts
// @Description Get all posts
// @Tags posts
// @Accept json
// @Produce json
// @Param page query int false "Page number"
// @Param limit query int false "Number of posts per page"
// @Success 200 {object} models.PaginatedPosts
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /posts/ [get]
func GetPosts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	pageStr := r.URL.Query().Get("page")
	limitStr := r.URL.Query().Get("limit")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 {
		limit = 5
	}
	offset := (page - 1) * limit

	var total int
	err = database.DB.QueryRow("SELECT COUNT(*) FROM posts").Scan(&total)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	rows, err := database.DB.Query(
		"SELECT id, title, content, created_at, updated_at FROM posts LIMIT ? OFFSET ?",
		limit, offset,
	)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var posts []models.Post
	for rows.Next() {
		var p models.Post
		rows.Scan(&p.ID, &p.Title, &p.Content, &p.CreatedAt, &p.UpdatedAt)
		posts = append(posts, p)
	}

	resp := map[string]interface{}{
		"data":        posts,
		"page":        page,
		"limit":       limit,
		"total":       total,
		"total_pages": (total + limit - 1) / limit,
	}

	json.NewEncoder(w).Encode(resp)
}

// CreatePost creates a new post
// @Summary Create a new post
// @Description Create a new post
// @Tags posts
// @Accept json
// @Produce json
// @Param post body models.Post true "Post"
// @Success 201 {object} models.Post
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /post/ [post]
func CreatePost(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var p models.Post
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	stmt, err := database.DB.Prepare("INSERT INTO posts (title, content) VALUES (?, ?)")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer stmt.Close()

	res, err := stmt.Exec(p.Title, p.Content)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	id, _ := res.LastInsertId()
	p.ID = int(id)
	p.CreatedAt = time.Now().Format(time.RFC3339)
	p.UpdatedAt = time.Now().Format(time.RFC3339)
	json.NewEncoder(w).Encode(p)
}

// GetPostByID retrieves a post by its ID
// @Summary Get post by ID
// @Description Get post by ID
// @Tags posts
// @Accept json
// @Produce json
// @Param id path int true "Post ID"
// @Success 200 {object} models.Post
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /post/{id} [get]
func GetPostByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var p models.Post
	err = database.DB.QueryRow("SELECT id, title, content, created_at, updated_at FROM posts WHERE id = ?", id).Scan(&p.ID, &p.Title, &p.Content, &p.CreatedAt, &p.UpdatedAt)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(p)
}
