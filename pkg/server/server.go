package server

import (
	"final_project/pkg/api"
	"log"
	"net/http"
	"os"
)

func Run() error {

	port := os.Getenv("TODO_PORT")

	api.Init()

	webDir := "./web"

	http.Handle("/", http.FileServer(http.Dir(webDir)))

	log.Println("Server running on port:", port)
	return http.ListenAndServe(":"+port, nil)
}
