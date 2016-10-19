package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/gorilla/rpc"
	"github.com/gorilla/rpc/json"
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

	// first pass
	http.HandleFunc("/api/login", explorer.Gzip(cors.Wrap(explorer.Connect())))
	http.HandleFunc("/api/metadata", explorer.Gzip(cors.Wrap(explorer.Metadata())))
	http.HandleFunc("/api/search", explorer.Gzip(cors.Wrap(explorer.Search())))
	http.HandleFunc("/api/object", explorer.Gzip(cors.Wrap(explorer.GetObject())))

	// newer gorilla rpc
	s := rpc.NewServer()
	s.RegisterCodec(json.NewCodec(), "application/json")
	s.RegisterService(new(explorer.ConnectionService), "")
	http.Handle("/rpc", handlers.CompressHandler(handlers.CORS()(s)))

	log.Println("Server starting: http://localhost:" + *port)
	log.Fatal(http.ListenAndServe(":"+*port, nil))
}
