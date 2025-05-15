package handlers

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/muttayoshi/goblog/database"
	"github.com/muttayoshi/goblog/models"
	"net/http"
	"strconv"
)

// GetPosts godoc
// @Summary Get all posts
// @Description Get all posts
// @Tags posts
// @Accept json
// @Produce json
// @Success 200 {array} models.Post
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /posts [get]
func GetPosts(w http.ResponseWriter, r *http.Request) {
	rows, err := database.DB.Query("SELECT id, title, content, created_at, updated_at FROM posts")
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
	json.NewEncoder(w).Encode(posts)
}

func CreatePost(w http.ResponseWriter, r *http.Request) {
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
	json.NewEncoder(w).Encode(p)
}

func GetPostByID(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var p models.Post
	err = database.DB.QueryRow("SELECT id, title, content FROM posts WHERE id = ?", id).Scan(&p.ID, &p.Title, &p.Content)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(p)
}
