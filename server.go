package main

import (
	"context" // Often needed for Ent operations like migrations
	"log"
	"net/http"
	"os"
	"webmane_go/cmd"
	"webmane_go/db" // Keep db import

	// Import ent for migration check if used
	"webmane_go/graph"
	"webmane_go/music" // Import music package

	"github.com/go-chi/chi"
	"github.com/gorilla/websocket"
	"github.com/rs/cors"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"

	_ "github.com/jackc/pgx/v4/stdlib"
)

// ▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒
// ▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒░░░░░░░░░░▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒
// ▒▒▒▒▒▒▒▒▒▒▒▒▒░░░░░░░░░░░░░░░░░░░▒▒▒▒▒▒▒▒▒▒▒▒
// ▒▒▒▒▒▒▒▒▒▒░░░░░░░░░░░░░░░░░░░░░░░░░▒▒▒▒▒▒▒▒▒
// ▒▒▒▒▒▒▒▒░░░░░░░░░░░░░░░░░░░░░░░░░░░░░▒▒▒▒▒▒▒
// ▒▒▒▒▒▒░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░▄░░▒▒▒▒▒
// ▒▒▒▒▒░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░██▌░░▒▒▒▒
// ▒▒▒▒░░░░░░░░░░░░░░░░░░░░░░░░░░░▄▄███▀░░░░▒▒▒
// ▒▒▒░░░░░░░░░░░░░░░░░░░░░░░░░░░█████░▄█░░░░▒▒
// ▒▒░░░░░░░░░░░░░░░░░░░░░░░░░░▄████████▀░░░░▒▒
// ▒▒░░░░░░░░░░░░░░░░░░░░░░░░▄█████████░░░░░░░▒
// ▒░░░░░░░░░░░░░░░░░░░░░░░░░░▄███████▌░░░░░░░▒
// ▒░░░░░░░░░░░░░░░░░░░░░░░░▄█████████░░░░░░░░▒
// ▒░░░░░░░░░░░░░░░░░░░░░▄███████████▌░░░░░░░░▒
// ▒░░░░░░░░░░░░░░░▄▄▄▄██████████████▌░░░░░░░░▒
// ▒░░░░░░░░░░░▄▄███████████████████▌░░░░░░░░░▒
// ▒░░░░░░░░░▄██████████████████████▌░░░░░░░░░▒
// ▒░░░░░░░░████████████████████████░░░░░░░░░░▒
// ▒█░░░░░▐██████████▌░▀▀███████████░░░░░░░░░░▒
// ▐██░░░▄██████████▌░░░░░░░░░▀██▐█▌░░░░░░░░░▒▒
// ▒██████░█████████░░░░░░░░░░░▐█▐█▌░░░░░░░░░▒▒
// ▒▒▀▀▀▀░░░██████▀░░░░░░░░░░░░▐█▐█▌░░░░░░░░▒▒▒
// ▒▒▒▒▒░░░░▐█████▌░░░░░░░░░░░░▐█▐█▌░░░░░░░▒▒▒▒
// ▒▒▒▒▒▒░░░░███▀██░░░░░░░░░░░░░█░█▌░░░░░░▒▒▒▒▒
// ▒▒▒▒▒▒▒▒░▐██░░░██░░░░░░░░▄▄████████▄▒▒▒▒▒▒▒▒
// ▒▒▒▒▒▒▒▒▒██▌░░░░█▄░░░░░░▄███████████████████
// ▒▒▒▒▒▒▒▒▒▐██▒▒░░░██▄▄███████████████████████
// ▒▒▒▒▒▒▒▒▒▒▐██▒▒▄████████████████████████████
// ▒▒▒▒▒▒▒▒▒▒▄▄████████████████████████████████
// ████████████████████████████████████████████

const defaultPort = "8080"

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	// --- Initialize Ent Client ---
	entClient, err := db.InitializeEntClient()
	if err != nil {
		log.Fatalf("Failed to initialize Ent client: %v", err)
	}
	defer entClient.Close()

	// --- Optional: Run migrations (Development/Initial Setup) ---
	// Uncomment this block if you want Ent to create/update tables on startup
	// based on your ent/schema definitions. Use Atlas for production.
	log.Println("Checking database schema...")
	if err := entClient.Schema.Create(
		context.Background(),
		// Add migration options if needed, e.g.,
		// migrate.WithDropIndex(true),
		// migrate.WithDropColumn(true),
	); err != nil {
		log.Fatalf("Failed creating/updating schema resources: %v", err)
	}
	log.Println("Database schema is up-to-date.")
	// --- End Optional Migration ---

	// --- Command Line Setup ---
	resolver := &graph.Resolver{EntClient: entClient}
	// Pass only the resolver to the wrapper
	rootCmd := cmd.CommandContextWrapper(resolver)

	// Execute command logic (e.g., if "seed" argument is passed)
	// Cobra's Execute handles argument parsing and running the correct command.
	// If no command args are given, it might just configure and proceed.
	// If a command like "seed" runs, it might exit here after completion.
	if cmd_err := rootCmd.Execute(); cmd_err != nil {
		// Log the error and exit if a command failed.
		// Don't proceed to start the server if a CLI command was meant to run and failed.
		log.Fatalf("Command execution failed: %v", cmd_err)
		// os.Exit(1) // Exit explicitly if Execute doesn't always exit on command run
	}
	// If Execute doesn't exit after running a command (like 'help'),
	// the server startup below will still happen. This might be desired or not.

	// --- HTTP Server Setup ---
	router := chi.NewRouter()

	router.Use(cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5173", "http://127.0.0.1:5173"}, // Add other origins if needed
		AllowCredentials: true,
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		MaxAge:           300,   // Maximum value not ignored by any browser
		Debug:            false, // Set to true for debugging CORS issues
	}).Handler)

	// --- GraphQL Server ---
	srv := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: resolver}))

	srv.AddTransport(&transport.Websocket{
		Upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				// TODO: Make origin check more robust and configurable
				origin := r.Header.Get("Origin")
				allowedOrigins := map[string]bool{
					"http://localhost:5173":    true,
					"http://127.0.0.1:5173":    true,
					"http://webmane.net":       true, // Example production origin
					"https://webmane.net":      true, // Example production origin
					"http://localhost:" + port: true, // Allow same origin for playground/dev
				}
				return allowedOrigins[origin]
			},
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
		},
	})

	router.Handle("/", playground.Handler("GraphQL playground", "/query"))
	router.Handle("/query", srv)

	// --- REST Endpoint for Music Streaming ---
	// Pass the entClient to the refactored GetMusic handler
	router.HandleFunc("/music", music.GetMusic(entClient))

	// --- Remove or Refactor Seed Endpoint ---
	// The HTTP seed endpoint might be redundant now with the CLI command.
	// If still needed, it must be refactored to use the entClient.
	// router.HandleFunc("/music/seed", music.SeedMusic(entClient)) // Example if refactored

	log.Printf("GraphQL playground available at http://localhost:%s/", port)
	log.Printf("Starting server on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}
