package main

import (
	"log"
	"net/http"
	"os"

	"github.com/jpfielding/gorets/explorer/server"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}
	reactPath := os.Getenv("REACT_PATH")
	if reactPath == "" {
		reactPath = "./"
	}
	http.Handle("/", http.FileServer(http.Dir(reactPath+"client")))
	http.HandleFunc("/api/login", server.Login())

	log.Println("Server starting: http://localhost:" + port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
	log.Println("Server starte: http://localhost:" + port)
}
