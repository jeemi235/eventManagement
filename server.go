package main

import (
	"10/graph"
	"log"
	"net/http"
	"os"

	dataloaders "10/dataloaders"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

const defaultPort = "5053"

func main() {
	//creating new chi router

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(dataloaders.LoadMiddleware)
	//r.Use(middlewares.DbContext)

	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	//srv := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{}}))

	//This will open graphql playground
	r.Handle("/", playground.Handler("GraphQL plaground: ", "/query"))
	r.Handle("/query", handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{}})))

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
