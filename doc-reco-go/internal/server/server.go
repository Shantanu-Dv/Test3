package server

import (
	"doc-reco-go/internal/config"
	"doc-reco-go/internal/middleware"
	"doc-reco-go/internal/provider"
	"fmt"
	"log"
	"net/http"

	"gopkg.in/DataDog/dd-trace-go.v1/contrib/gorilla/mux"
)

type Server struct {
	Router *mux.Router //root domain router
}

func (server *Server) InitializeServer() error {
	err := config.LoadConfig()
	if err != nil {
		log.Fatalln(err)
	}

	r := mux.NewRouter(mux.WithServiceName(config.Config.Datadog.ServiceName))
	r.Use(middleware.SentryMiddleware)
	r.Use(middleware.LoggingMiddleware)

	server.Router = r
	server.InitializeRoutes()
	return provider.InitializeProvider()
}

func (server *Server) Run() {
	defer provider.ReleaseProviderResources()

	err := server.InitializeServer()
	if err != nil {
		panic(err)
	}

	port := config.Config.ServerPort
	s := &http.Server{
		Addr:    fmt.Sprintf(":%s", port),
		Handler: server.Router,
	}
	log.Printf("server running on port %s\n", port)
	log.Fatalln(s.ListenAndServe())
}
