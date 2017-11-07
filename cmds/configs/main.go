package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/gorilla/rpc"
	gson "github.com/gorilla/rpc/json"
	"github.com/jpfielding/gorets/config"
)

func main() {
	port := flag.String("port", "8888", "http port")
	configPath := flag.String("configs", "source-configs", "The configurations for this service")

	flag.Parse()

	cfgs, err := config.ImportFrom(*configPath)
	if err != nil {
		panic(err)
	}
	fmt.Printf("loaded %d configs\n", len(cfgs))

	// gorilla rpc
	s := rpc.NewServer()
	s.RegisterCodec(gson.NewCodec(), "application/json")
	s.RegisterService(&config.ConfigService{Configs: cfgs}, "")

	cors := handlers.CORS(
		handlers.AllowedMethods([]string{"OPTIONS", "POST", "GET", "HEAD"}),
		handlers.AllowedHeaders([]string{"Accept", "Content-Type", "Content-Length", "Accept-Encoding", "X-CSRF-Token", "Authorization"}),
	)
	// rpc calls
	http.Handle("/rpc", handlers.CompressHandler(cors(s)))

	log.Println("Server starting: http://localhost:" + *port)
	log.Fatal(http.ListenAndServe(":"+*port, nil))
}
