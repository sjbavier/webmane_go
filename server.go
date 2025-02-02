package main

import (
	"log"
	"net/http"
	"os"
	"webmane_go/cmd"
	"webmane_go/db"
	"webmane_go/graph"
	"webmane_go/music"

	"github.com/go-chi/chi"
	"github.com/gorilla/websocket"
	"github.com/rs/cors"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
)

const defaultPort = "8080"

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	// initialize database pool
	dbPool, db_err := db.ConnectToDb()

	if db_err != nil {
		log.Fatalf("Connecting to database failed:\n %v", db_err)
	}
	defer dbPool.Close()

	// command line
	resolver := &graph.Resolver{DB: dbPool}
	rootCmd := cmd.CommandContextWrapper(dbPool, resolver)
	if cmd_err := rootCmd.Execute(); cmd_err != nil {
		log.Fatalf("command failed to execute  %v", cmd_err)
	}

	router := chi.NewRouter()

	// Add CORS middleware around every request
	// See https://github.com/rs/cors for full option listing
	router.Use(cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5173"},
		AllowCredentials: true,
		Debug:            true,
	}).Handler)

	srv := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{DB: dbPool}}))

	srv.AddTransport(&transport.Websocket{
		Upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				// Check against your desired domains here
				return r.Host == "webmane.net"
			},
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
		},
	})

	router.Handle("/", playground.Handler("GraphQL playground", "/query"))
	router.Handle("/query", srv)
	router.HandleFunc("/music", music.GetMusic(dbPool))
	// http.HandleFunc("/music/seed", music.SeedMusic(dbPool))

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}
