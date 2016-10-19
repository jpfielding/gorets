package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/jpfielding/gorets/explorer"
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
	conn := &explorer.Connection{
		WireLogFile: "/tmp/gorets/wire.conn.log",
	}
	// TODO deal with contexts in the web appropriately
	ctx := context.Background()
	http.Handle("/", http.FileServer(http.Dir(reactPath)))
	http.HandleFunc("/api/login", explorer.Login(ctx, conn))
	http.HandleFunc("/api/metadata", explorer.Metadata(ctx, conn))
	http.HandleFunc("/api/search", explorer.Search(ctx, conn))

	log.Println("Server starting: http://localhost:" + port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
