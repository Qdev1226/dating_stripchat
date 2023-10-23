package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/alexedwards/scs/v2"
	"github.com/kriskelly/dating-app-example/internal/dgraph"
	"github.com/kriskelly/dating-app-example/internal/graph"
	"github.com/kriskelly/dating-app-example/internal/graph/generated"
)

const defaultPort = "3000"

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	sessionManager := scs.New()
	sessionManager.Lifetime = 24 * time.Hour

	dgraphClient := dgraph.NewClient()
	dgraphClient.Connect()
	defer dgraphClient.Close()

	resolver := graph.NewResolver(sessionManager, dgraphClient)
	srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: resolver}))

	mux := http.NewServeMux()
	mux.HandleFunc("/", playground.Handler("GraphQL playground", "/query"))
	mux.Handle("/query", srv)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, sessionManager.LoadAndSave(mux)))
}
