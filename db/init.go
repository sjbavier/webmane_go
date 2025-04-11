package db

import (
	"fmt"
	"log"
	"os"

	"webmane_go/ent" // Import the generated ent package

	// Import the pgx stdlib driver needed by Ent for PostgreSQL
	"entgo.io/ent/dialect"
	_ "github.com/jackc/pgx/v4/stdlib"
)

// InitializeEntClient connects to the database and returns an Ent client.
func InitializeEntClient() (*ent.Client, error) {
	// It's better practice to get the URL from an environment variable
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		// Fallback to a default for local development if the env var isn't set
		log.Println("WARNING: DATABASE_URL environment variable not set. Using default local connection.")
		dbURL = "postgres://user:password@localhost:5432/webmane_go?sslmode=disable" // Added sslmode=disable, often needed locally
	}

	// Use ent.Open to create the client
	client, err := ent.Open(dialect.Postgres, dbURL)
	if err != nil {
		// Use fmt.Errorf for better error wrapping
		return nil, fmt.Errorf("failed opening connection to postgres via Ent: %w", err)
	}

	// Optionally, you can add a Ping check here in development
	// if err := client.Ping(context.Background()); err != nil {
	//  	return nil, fmt.Errorf("failed pinging database via Ent: %w", err)
	// }

	fmt.Println("Connected to database via Ent client")
	return client, nil
}

// --- Keep the old function temporarily if needed by cmd or music packages ---
// You should refactor cmd and music packages to use *ent.Client eventually
// and then remove this function and the direct pgxpool dependency.

// import (
// 	"context"
// 	"github.com/jackc/pgx/v4/pgxpool"
// )
//
// func ConnectToDb() (*pgxpool.Pool, error) {
// 	DATABASE_URL := "postgres://user:password@localhost:5432/webmane_go"
// 	pool, err := pgxpool.Connect(context.Background(), DATABASE_URL)
//
// 	if err != nil {
// 		fmt.Sprintf("Error connecting to db %v", err)
// 		return pool, err
// 	}
// 	fmt.Println("Connected to database (legacy pgxpool)")
// 	return pool, nil
// }

// --- End Temporary Section ---

