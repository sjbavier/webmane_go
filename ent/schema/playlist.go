package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

type Playlist struct {
	ent.Schema
}

// Fields of the Playlist.
func (Playlist) Fields() []ent.Field {
	return []ent.Field{
		// ID is handled by Ent automatically

		field.String("name").
			NotEmpty().
			MaxLen(255), // A reasonable max length for a name

		// Use time.Time for dates/timestamps
		field.Time("last_update").
			Default(time.Now).
			UpdateDefault(time.Now),

		field.Time("last_accessed").
			Optional(), // Can be NULL if never accessed

		// Use string for TEXT type or path/reference
		field.String("cover_art").
			Optional(), // Allows NULL in the database
	}
}

// Edges of the Playlist.
func (Playlist) Edges() []ent.Edge {
	return []ent.Edge{
		// Defines a M2M relationship with Music (Song)
		// A playlist can have many songs.
		edge.To("songs", Music.Type),
	}
}

// Indexes of the Playlist.
func (Playlist) Indexes() []ent.Index {
	return []ent.Index{
		// Index on the name field for faster lookups
		index.Fields("name"),
	}
}

// Annotations of the Playlist.
func (Playlist) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{ // Add this annotation
			Table: "playlists", // Specify your desired table name here
		},
	}
}
