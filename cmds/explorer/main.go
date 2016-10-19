package main

import (
	"context"
	"flag"
	"log"
	"net/http"

	"github.com/jpfielding/gorets/explorer"
)

func main() {
	port := flag.String("port", "8000", "http port")
	react := flag.String("react", "../../explorer/client/build", "ReactJS path")

	flag.Parse()

	// TODO this needs to be bound to a client cookie
	conn := &explorer.Connection{
		WireLogFile: "/tmp/gorets/wire.conn.log",
	}
	// TODO deal with contexts in the web appropriately
	ctx := context.Background()
	http.Handle("/", http.FileServer(http.Dir(*react)))

	cors := explorer.NewCors("*")

	http.HandleFunc("/api/login", explorer.Gzip(cors.Wrap(explorer.Login(ctx, conn))))
	http.HandleFunc("/api/metadata", explorer.Gzip(cors.Wrap(explorer.Metadata(ctx, conn))))
	http.HandleFunc("/api/search", explorer.Gzip(cors.Wrap(explorer.Search(ctx, conn))))
	http.HandleFunc("/api/object", explorer.Gzip(cors.Wrap(explorer.GetObject(ctx, conn))))

	log.Println("Server starting: http://localhost:" + *port)
	log.Fatal(http.ListenAndServe(":"+*port, nil))
}
