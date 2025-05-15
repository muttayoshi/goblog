package main

import (
	"github.com/muttayoshi/goblog/cmd"
	"github.com/muttayoshi/goblog/database"
	_ "github.com/muttayoshi/goblog/docs"
	"github.com/muttayoshi/goblog/log"
	"github.com/muttayoshi/goblog/routes"
	"github.com/swaggo/http-swagger"
	"net/http"
)

// @title Blog API
// @version 1.0
// @description API untuk blog sederhana dengan Golang dan SQLite
// @host localhost:8080
// @BasePath /
func main() {
	database.InitDB("goblog.db")
	r := routes.SetupRouter()
	r.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	log.Println("Starting server on :8080")
	log.Fatal(http.ListenAndServe(":8080", r))

	cmd.Execute()

}
