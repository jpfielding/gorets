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
	user := &server.User{
		WireLogFile: "/tmp/gorets/wire.log"
	}
	// TODO deal with contexts in the web appropriately
	ctx := context.Background()
	http.Handle("/", http.FileServer(http.Dir(reactPath)))
	http.HandleFunc("/api/login", server.Login(ctx, user))
	http.HandleFunc("/api/metadata", server.Metadata(ctx, user))

	log.Println("Server starting: http://localhost:" + port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
	log.Println("Server start: http://localhost:" + port)
}
