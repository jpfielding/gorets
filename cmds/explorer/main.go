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

	http.Handle("/", http.FileServer(http.Dir(*react)))

	// gorilla rpc
	s := rpc.NewServer()
	s.RegisterCodec(json.NewCodec(), "application/json")
	s.RegisterService(&explorer.MetadataService{}, "")
	s.RegisterService(&explorer.SearchService{}, "")
	s.RegisterService(&explorer.ObjectService{}, "")

	cors := handlers.CORS(
		handlers.AllowedMethods([]string{"OPTIONS", "POST", "GET", "HEAD"}),
		handlers.AllowedHeaders([]string{"Accept", "Content-Type", "Content-Length", "Accept-Encoding", "X-CSRF-Token", "Authorization", "Access-Control-Allow-Origin"}),
	)
	// rpc calls
	http.Handle("/rpc", handlers.CompressHandler(cors(s)))

	// websocket wire logs  TODO fix this or remove it
	// http.Handle("/wirelog", explorer.WireLogSocket(explorer.WirelogUpgrader))

	log.Println("Server starting: http://localhost:" + *port)
	log.Fatal(http.ListenAndServe(":"+*port, nil))
}
