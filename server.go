package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/bortexel/stats-server/database"
	"github.com/bortexel/stats-server/graph"
	"github.com/bortexel/stats-server/graph/generated"
)

const defaultPort = "8080"

func main() {
	uri, ok := os.LookupEnv("MONGO_CONNECTION_URI")
	if !ok {
		log.Fatalln("MONGO_CONNECTION_URI variable is empty")
		return
	}

	err := database.InitDatabase(uri)
	if err != nil {
		log.Println("Unable to init database:", err)
		return
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &graph.Resolver{}}))

	router := chi.NewRouter()
	router.Use(middleware.Logger)

	key, ok := os.LookupEnv("MUTATION_KEY")
	if !ok {
		log.Fatalln("MUTATION_KEY is not set")
	} else {
		router.Use(Authorization(key))
	}

	router.Handle("/", playground.Handler("GraphQL playground", "/query"))
	router.Handle("/query", srv)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}
