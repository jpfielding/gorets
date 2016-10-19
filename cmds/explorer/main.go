package main

import (
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
	conns := map[string]explorer.Connection{}
	explorer.JSONLoad("/tmp/gorets/connections.json", &conns)

	http.Handle("/", http.FileServer(http.Dir(*react)))

	cors := explorer.NewCors("*")

	http.HandleFunc("/api/login", explorer.Gzip(cors.Wrap(explorer.Connect(conns))))
	http.HandleFunc("/api/metadata", explorer.Gzip(cors.Wrap(explorer.Metadata(conns))))
	http.HandleFunc("/api/search", explorer.Gzip(cors.Wrap(explorer.Search(conns))))
	http.HandleFunc("/api/object", explorer.Gzip(cors.Wrap(explorer.GetObject(conns))))

	log.Println("Server starting: http://localhost:" + *port)
	log.Fatal(http.ListenAndServe(":"+*port, nil))
}
