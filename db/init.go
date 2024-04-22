package db

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v4/pgxpool"
)

func ConnectToDb() (*pgxpool.Pool, error) {
	DATABASE_URL := "postgres://user:password@localhost:5432/webmane_go"
	pool, err := pgxpool.Connect(context.Background(), DATABASE_URL)

	if err != nil {
		fmt.Sprintf("Error connecting to db %v", err)
		return pool, err
	}
	fmt.Println("Connected to database")
	return pool, nil
}
