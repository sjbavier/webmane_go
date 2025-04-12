// Code generated by ent, DO NOT EDIT.

package playlist

import (
	"time"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
)

const (
	// Label holds the string label denoting the playlist type in the database.
	Label = "playlist"
	// FieldID holds the string denoting the id field in the database.
	FieldID = "id"
	// FieldName holds the string denoting the name field in the database.
	FieldName = "name"
	// FieldLastUpdate holds the string denoting the last_update field in the database.
	FieldLastUpdate = "last_update"
	// FieldLastAccessed holds the string denoting the last_accessed field in the database.
	FieldLastAccessed = "last_accessed"
	// FieldCoverArt holds the string denoting the cover_art field in the database.
	FieldCoverArt = "cover_art"
	// EdgeSongs holds the string denoting the songs edge name in mutations.
	EdgeSongs = "songs"
	// Table holds the table name of the playlist in the database.
	Table = "playlists"
	// SongsTable is the table that holds the songs relation/edge. The primary key declared below.
	SongsTable = "playlist_songs"
	// SongsInverseTable is the table name for the Music entity.
	// It exists in this package in order to avoid circular dependency with the "music" package.
	SongsInverseTable = "music_ent"
)

// Columns holds all SQL columns for playlist fields.
var Columns = []string{
	FieldID,
	FieldName,
	FieldLastUpdate,
	FieldLastAccessed,
	FieldCoverArt,
}

var (
	// SongsPrimaryKey and SongsColumn2 are the table columns denoting the
	// primary key for the songs relation (M2M).
	SongsPrimaryKey = []string{"playlist_id", "music_id"}
)

// ValidColumn reports if the column name is valid (part of the table columns).
func ValidColumn(column string) bool {
	for i := range Columns {
		if column == Columns[i] {
			return true
		}
	}
	return false
}

var (
	// NameValidator is a validator for the "name" field. It is called by the builders before save.
	NameValidator func(string) error
	// DefaultLastUpdate holds the default value on creation for the "last_update" field.
	DefaultLastUpdate func() time.Time
	// UpdateDefaultLastUpdate holds the default value on update for the "last_update" field.
	UpdateDefaultLastUpdate func() time.Time
)

// OrderOption defines the ordering options for the Playlist queries.
type OrderOption func(*sql.Selector)

// ByID orders the results by the id field.
func ByID(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldID, opts...).ToFunc()
}

// ByName orders the results by the name field.
func ByName(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldName, opts...).ToFunc()
}

// ByLastUpdate orders the results by the last_update field.
func ByLastUpdate(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldLastUpdate, opts...).ToFunc()
}

// ByLastAccessed orders the results by the last_accessed field.
func ByLastAccessed(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldLastAccessed, opts...).ToFunc()
}

// ByCoverArt orders the results by the cover_art field.
func ByCoverArt(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldCoverArt, opts...).ToFunc()
}

// BySongsCount orders the results by songs count.
func BySongsCount(opts ...sql.OrderTermOption) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborsCount(s, newSongsStep(), opts...)
	}
}

// BySongs orders the results by songs terms.
func BySongs(term sql.OrderTerm, terms ...sql.OrderTerm) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborTerms(s, newSongsStep(), append([]sql.OrderTerm{term}, terms...)...)
	}
}
func newSongsStep() *sqlgraph.Step {
	return sqlgraph.NewStep(
		sqlgraph.From(Table, FieldID),
		sqlgraph.To(SongsInverseTable, FieldID),
		sqlgraph.Edge(sqlgraph.M2M, false, SongsTable, SongsPrimaryKey...),
	)
}
