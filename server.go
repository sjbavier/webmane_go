package main

import (
	"log"
	"net/http"
	"os"
	"webmane_go/graph"
	"webmane_go/music"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	// "github.com/jmoiron/sqlx"
)

const defaultPort = "8080"

// var db *sqlx.DB

// func init() {
// 	var err error
// 	db, err := sqlx.Connect("postgres", os.Getenv("DATABASE_URL"))
// 	if err != nil {
// 		log.Fatalf("failed to connect to db: %v", err)
// 	}
// }

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	srv := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{}}))

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)
	http.HandleFunc("/music", music.GetMusic)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
