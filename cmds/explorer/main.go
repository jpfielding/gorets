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

	// TODO remove when migrated over
	http.HandleFunc("/api/login", cors(explorer.Connect()))
	http.HandleFunc("/api/metadata", cors(explorer.Metadata()))

	// gorilla rpc
	s := rpc.NewServer()
	s.RegisterCodec(json.NewCodec(), "application/json")
	s.RegisterService(new(explorer.ConnectionService), "")
	s.RegisterService(new(explorer.MetadataService), "")
	s.RegisterService(new(explorer.SearchService), "")
	s.RegisterService(new(explorer.ObjectService), "")
	// http.Handle("/rpc", handlers.CORS()(s))
	http.Handle("/rpc", handlers.CompressHandler(handlers.CORS()(s)))

	log.Println("Server starting: http://localhost:" + *port)
	log.Fatal(http.ListenAndServe(":"+*port, nil))
}

// TODO kill this when custom handlers go away
func cors(wrapped func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if origin := r.Header.Get("Origin"); origin != "" {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
			w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		}
		// Stop here if its Preflighted OPTIONS request
		if r.Method == "OPTIONS" {
			return
		}
		wrapped(w, r)
	}
}
