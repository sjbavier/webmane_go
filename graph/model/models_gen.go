// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package model

type Mutation struct {
}

type Query struct {
}

type SongInput struct {
	ID          *string `json:"id,omitempty"`
	Path        string  `json:"path"`
	Title       *string `json:"title,omitempty"`
	Artist      *string `json:"artist,omitempty"`
	Album       *string `json:"album,omitempty"`
	Genre       *string `json:"genre,omitempty"`
	ReleaseYear *string `json:"release_year,omitempty"`
}
