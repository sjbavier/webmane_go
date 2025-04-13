package utils

import (
	"strconv"
	"time" // Import time package
	"webmane_go/ent"
	"webmane_go/graph/model"
)

// ISO8601Format defines the standard time format for GraphQL responses
const ISO8601Format = "2006-01-02T15:04:05Z07:00"

func MapEntMusicToModelSong(entMusic *ent.Music) *model.Song {
	if entMusic == nil {
		return nil
	}
	// Handle potential nil pointers for optional string fields from Ent
	// (Ent uses value types, so they won't be nil, but might be empty strings if cleared)
	// The GraphQL model uses *string, so we need to take addresses.
	// If the DB field was nullable and Ent generated Nillable fields, the logic would differ slightly.

	// Helper function to get pointer to string or nil if empty
	stringPtr := func(s string) *string {
		if s == "" {
			return nil // Return nil if the string is empty (assuming empty means not set)
		}
		return &s
	}

	return &model.Song{
		ID:          strconv.Itoa(entMusic.ID),
		Path:        entMusic.Path,
		LastUpdate:  entMusic.LastUpdate.Format(ISO8601Format), // Use consistent format
		Title:       stringPtr(entMusic.Title),
		Artist:      stringPtr(entMusic.Artist),
		Album:       stringPtr(entMusic.Album),
		Genre:       stringPtr(entMusic.Genre),
		ReleaseYear: stringPtr(entMusic.ReleaseYear),
		CoverArt:    stringPtr(entMusic.CoverArt),
	}
}

// MapEntPlaylistToModelPlaylist converts an Ent Playlist (with potentially loaded songs)
// to a GraphQL Playlist model.
func MapEntPlaylistToModelPlaylist(entPlaylist *ent.Playlist) *model.Playlist {
	if entPlaylist == nil {
		return nil
	}

	// Helper function to format optional time
	formatOptionalTime := func(t time.Time) *string {
		if t.IsZero() { // Check if the time is the zero value (often indicates NULL in DB)
			return nil
		}
		formatted := t.Format(ISO8601Format)
		return &formatted
	}

	// Helper function to get pointer to string or nil if empty
	stringPtr := func(s string) *string {
		if s == "" {
			return nil // Return nil if the string is empty
		}
		return &s
	}

	// Map the songs edge if it's loaded
	var modelSongs []*model.Song
	if entPlaylist.Edges.Songs != nil { // Check if the edge was loaded
		modelSongs = make([]*model.Song, 0, len(entPlaylist.Edges.Songs)) // Pre-allocate slice
		for _, entSong := range entPlaylist.Edges.Songs {
			// Reuse the existing song mapper
			mappedSong := MapEntMusicToModelSong(entSong)
			if mappedSong != nil { // Ensure the mapped song isn't nil
				modelSongs = append(modelSongs, mappedSong)
			}
		}
	} else {
		// If edges were not loaded, return an empty (but non-nil) slice as per GraphQL schema [Song!]!
		modelSongs = make([]*model.Song, 0)
	}

	return &model.Playlist{
		ID:           strconv.Itoa(entPlaylist.ID),
		Name:         entPlaylist.Name,
		LastUpdate:   entPlaylist.LastUpdate.Format(ISO8601Format), // Required field
		LastAccessed: formatOptionalTime(entPlaylist.LastAccessed), // Optional field
		CoverArt:     stringPtr(entPlaylist.CoverArt),              // Optional field
		Songs:        modelSongs,                                   // Mapped songs slice
	}
}
