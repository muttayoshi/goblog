package routes

import (
	"github.com/gorilla/mux"
	"github.com/muttayoshi/goblog/handlers"
)

// SetupRouter initializes the router and sets up the routes
func SetupRouter() *mux.Router {
	r := mux.NewRouter()

	// Define your routes here
	r.HandleFunc("/posts", handlers.GetPosts).Methods("GET")
	r.HandleFunc("/post", handlers.CreatePost).Methods("POST")
	r.HandleFunc("/post/{id}", handlers.GetPostByID).Methods("GET")
	//r.HandleFunc("/posts/{id}", handlers.UpdatePost).Methods("PUT")
	//r.HandleFunc("/posts/{id}", handlers.DeletePost).Methods("DELETE")

	return r
}
