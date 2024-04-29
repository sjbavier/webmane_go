package graph

import (
	"webmane_go/graph/model"

	"github.com/jackc/pgx/v4/pgxpool"
)

//go:generate go run github.com/99designs/gqlgen generate

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	DB         *pgxpool.Pool
	MusicStore map[string]model.Song
}
