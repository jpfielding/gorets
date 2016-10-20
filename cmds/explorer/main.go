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
	http.HandleFunc("/api/login", custom(explorer.Connect()))
	http.HandleFunc("/api/metadata", custom(explorer.Metadata()))

	// gorilla rpc
	s := rpc.NewServer()
	s.RegisterCodec(json.NewCodec(), "application/json")
	s.RegisterService(new(explorer.ConnectionService), "")
	s.RegisterService(new(explorer.MetadataService), "")
	s.RegisterService(new(explorer.SearchService), "")
	s.RegisterService(new(explorer.ObjectService), "")

	cors := handlers.CORS(
		handlers.AllowedMethods([]string{"OPTIONS", "POST", "GET", "HEAD"}),
		handlers.AllowedHeaders([]string{"Accept", "Content-Type", "Content-Length", "Accept-Encoding", "X-CSRF-Token", "Authorization"}),
		// handlers.AllowedOriginValidator(func(origin string) bool {
		// 	return true
		// }),
	)
	http.Handle("/rpc", cors(s))
	// http.Handle("/rpc", cors(s))
	// http.Handle("/rpc", handlers.CompressHandler(cors(s)))
	log.Println("Server starting: http://localhost:" + *port)
	log.Fatal(http.ListenAndServe(":"+*port, nil))
}

// TODO kill this when custom handlers go away
func custom(wrapped http.HandlerFunc) http.HandlerFunc {
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
		wrapped.ServeHTTP(w, r)
	}
}
