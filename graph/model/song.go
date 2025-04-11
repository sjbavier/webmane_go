// /home/b4v1n4t0r/go_projects/webmane_go/graph/model/song.go
package model

// Song represents the GraphQL Song type.
// Fields correspond to the schema definition in schema.graphqls.
// Optional fields are represented by pointers.
type Song struct {
	ID          string  `json:"id"`           // Corresponds to GraphQL ID! (mapped to string)
	Path        string  `json:"path"`         // Corresponds to GraphQL String!
	LastUpdate  string  `json:"lastUpdate"`   // Corresponds to GraphQL String!
	Title       *string `json:"title"`        // Corresponds to GraphQL String (optional -> pointer)
	Artist      *string `json:"artist"`       // Corresponds to GraphQL String (optional -> pointer)
	Album       *string `json:"album"`        // Corresponds to GraphQL String (optional -> pointer)
	Genre       *string `json:"genre"`        // Corresponds to GraphQL String (optional -> pointer)
	ReleaseYear *string `json:"release_year"` // Corresponds to GraphQL String (optional -> pointer)
	CoverArt    *string `json:"cover_art"`    // Corresponds to GraphQL String (optional -> pointer)
}

// Note: It's generally recommended to let gqlgen generate this file
// (often named models_gen.go) based on your schema.graphqls by running:
// go run github.com/99designs/gqlgen generate
//
// If you have manually created this file, ensure it stays in sync
// with your schema and resolver expectations.
