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
		reactPath = "../../explorer/client"
	}

	// TODO this needs to be bound to a client cookie
	session := &server.Session{
		WireLogFile: "/tmp/gorets/wire.log"
	}
	http.Handle("/", http.FileServer(http.Dir(reactPath)))
	http.HandleFunc("/api/login", server.Login(session))
	http.HandleFunc("/api/metadata", server.Metadata(session))

	log.Println("Server starting: http://localhost:" + port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
	log.Println("Server start: http://localhost:" + port)
}
