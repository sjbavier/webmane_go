package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// Music holds the schema definition for the Music entity.
type Music struct {
	ent.Schema
}

// Fields of the Music.
func (Music) Fields() []ent.Field {
	return []ent.Field{
		// ID is handled by Ent automatically (int/int64)

		field.String("path").
			MaxLen(500).
			Unique(). // Corresponds to unique constraint
			NotEmpty(), // Assuming path should not be empty

		// Use time.Time for dates/timestamps
		// UpdateDefault sets the value on every update automatically
		field.Time("last_update").
			Default(time.Now).
			UpdateDefault(time.Now),


		field.String("title").
			MaxLen(500).
			Optional(), // Allows NULL in the database

		field.String("artist").
			MaxLen(500).
			Optional(),

		field.String("album").
			MaxLen(500).
			Optional(),

		field.String("genre").
			MaxLen(500).
			Optional(),

		// Keep as string based on your SQL, though int might be better if always numeric
		field.String("release_year").
			MaxLen(50).
			Optional(),

		// Use string for TEXT type
		field.String("cover_art").
			Optional(),
			// If you need truly unlimited text (like TEXT in Postgres),
			// you might need SchemaType, but often default string is fine.
			// SchemaType(map[string]string{
			// 	dialect.Postgres: "text",
			// }),
	}
}

// Edges of the Music.
func (Music) Edges() []ent.Edge {
	// No relationships defined in your SQL schema yet
	return nil
}

// Indexes of the Music.
func (Music) Indexes() []ent.Index {
	return []ent.Index{
		// Corresponds to: create index if not exists idx_music_path on public.music (path);
		index.Fields("path"),
	}
}
