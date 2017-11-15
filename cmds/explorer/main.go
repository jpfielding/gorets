package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/gorilla/rpc"
	"github.com/gorilla/rpc/json"
	"github.com/jpfielding/gorets/config"
	"github.com/jpfielding/gorets/explorer"
	"github.com/jpfielding/gowirelog/wirelog"
)

func main() {
	port := flag.String("port", "8000", "http port")
	react := flag.String("react", "../../explorer/client/build", "ReactJS path")
	configService := flag.String("configService", "http://localhost:8888/rpc", "The configuration service")

	flag.Parse()

	http.Handle("/", http.FileServer(http.Dir(*react)))

	// use our client tool to get configs from another microservice
	rpcClient := config.Client{
		EndPoint: *configService,
		Client: http.Client{
			Transport: wirelog.NewHTTPTransport(),
		},
	}
	// adapt this to our config service
	rpcClientFunc := func(args *config.ListArgs) ([]config.Config, error) {
		reply, err := rpcClient.List(*args)
		return reply.Configs, err
	}

	// gorilla rpc
	s := rpc.NewServer()
	s.RegisterCodec(json.NewCodec(), "application/json")
	// TODO remove this connection service and replace with configservice (or reference the service from the react)
	s.RegisterService(&config.RPCService{Configs: rpcClientFunc}, "ConfigService")
	s.RegisterService(&explorer.MetadataService{}, "")
	s.RegisterService(&explorer.SearchService{}, "")
	s.RegisterService(&explorer.ObjectService{}, "")

	cors := handlers.CORS(
		handlers.AllowedMethods([]string{"OPTIONS", "POST", "GET", "HEAD"}),
		handlers.AllowedHeaders([]string{"Accept", "Content-Type", "Content-Length", "Accept-Encoding", "X-CSRF-Token", "Authorization", "Access-Control-Allow-Origin"}),
	)
	// rpc calls
	http.Handle("/rpc", handlers.CompressHandler(cors(s)))

	// websocket wire logs
	http.Handle("/wirelog", explorer.WireLogSocket(explorer.WirelogUpgrader))

	log.Println("Server starting: http://localhost:" + *port)
	log.Fatal(http.ListenAndServe(":"+*port, nil))
}
