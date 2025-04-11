package graph

import (
	"context"
	"fmt"
	"strconv"

	// "time" // time might not be needed directly here anymore
	"webmane_go/ent"       // Import ent
	"webmane_go/ent/music" // Import generated music package
	"webmane_go/graph/model"
)

// Helper function to convert *ent.Music to *model.Song
func mapEntMusicToModelSong(entMusic *ent.Music) *model.Song {
	if entMusic == nil {
		return nil
	}
	return &model.Song{
		ID:          strconv.Itoa(entMusic.ID), // Convert int ID to string
		Path:        entMusic.Path,
		LastUpdate:  entMusic.LastUpdate.Format("2006-01-02T15:04:05Z07:00"), // Format time
		Title:       &entMusic.Title,         // Ent uses value types for optional strings
		Artist:      &entMusic.Artist,
		Album:       &entMusic.Album,
		Genre:       &entMusic.Genre,
		ReleaseYear: &entMusic.ReleaseYear,
		CoverArt:    &entMusic.CoverArt,
	}
}

// UpsertSong is the resolver for the upsertSong field.
func (r *mutationResolver) UpsertSong(ctx context.Context, input model.SongInput) (*model.Song, error) {
	// Ent handles last_update automatically via UpdateDefault(time.Now)

	if input.ID != nil && *input.ID != "" {
		// --- UPDATE ---
		id, err := strconv.Atoi(*input.ID)
		if err != nil {
			return nil, fmt.Errorf("invalid ID format: %w", err)
		}

		updater := r.EntClient.Music.UpdateOneID(id).
			SetPath(input.Path) // Path is required

		// Set optional fields only if provided
		if input.Title != nil {
			updater.SetTitle(*input.Title)
		} else {
			updater.ClearTitle() // Explicitly clear if needed, or SetNillableTitle if schema allows
		}
		if input.Artist != nil {
			updater.SetArtist(*input.Artist)
		} else {
			updater.ClearArtist()
		}
		if input.Album != nil {
			updater.SetAlbum(*input.Album)
		} else {
			updater.ClearAlbum()
		}
		if input.Genre != nil {
			updater.SetGenre(*input.Genre)
		} else {
			updater.ClearGenre()
		}
		if input.ReleaseYear != nil {
			updater.SetReleaseYear(*input.ReleaseYear)
		} else {
			updater.ClearReleaseYear()
		}
		if input.CoverArt != nil {
			updater.SetCoverArt(*input.CoverArt)
		} else {
			updater.ClearCoverArt()
		}

		updatedMusic, err := updater.Save(ctx)
		if err != nil {
			// Handle potential errors like Not Found or constraint violations
			return nil, fmt.Errorf("error updating song %s: %w", *input.ID, err)
		}
		return mapEntMusicToModelSong(updatedMusic), nil

	} else {
		// --- INSERT ---
		creator := r.EntClient.Music.Create().
			SetPath(input.Path) // Path is required

		// Set optional fields
		if input.Title != nil {
			creator.SetTitle(*input.Title)
		}
		if input.Artist != nil {
			creator.SetArtist(*input.Artist)
		}
		if input.Album != nil {
			creator.SetAlbum(*input.Album)
		}
		if input.Genre != nil {
			creator.SetGenre(*input.Genre)
		}
		if input.ReleaseYear != nil {
			creator.SetReleaseYear(*input.ReleaseYear)
		}
		if input.CoverArt != nil {
			creator.SetCoverArt(*input.CoverArt)
		}

		newMusic, err := creator.Save(ctx)
		if err != nil {
			// Handle potential errors like constraint violations (e.g., unique path)
			return nil, fmt.Errorf("error creating song for path %s: %w", input.Path, err)
		}
		return mapEntMusicToModelSong(newMusic), nil
	}
}

// Song is the resolver for the song field.
func (r *queryResolver) Song(ctx context.Context, id string) (*model.Song, error) {
	intID, err := strconv.Atoi(id)
	if err != nil {
		return nil, fmt.Errorf("invalid ID format: %w", err)
	}

	entMusic, err := r.EntClient.Music.Get(ctx, intID)
	if err != nil {
		// Handle ent.IsNotFound(err) specifically if needed
		return nil, fmt.Errorf("error fetching song %s: %w", id, err)
	}

	return mapEntMusicToModelSong(entMusic), nil
}

// Music is the resolver for the music field.
func (r *queryResolver) Music(ctx context.Context, pageNumber *int, pageSize *int, searchText *string) (*model.MusicResponse, error) {
	defaultPageNumber := 1
	defaultPageSize := 10

	pn := defaultPageNumber
	if pageNumber != nil {
		pn = *pageNumber
	}
	ps := defaultPageSize
	if pageSize != nil {
		ps = *pageSize
	}

	offset := (pn - 1) * ps

	// Base query
	query := r.EntClient.Music.Query()

	// Apply search filter if provided
	if searchText != nil && *searchText != "" {
		st := *searchText
		query = query.Where(
			music.Or(
				music.PathContainsFold(st), // Case-insensitive contains (like ILIKE '%...%')
				music.TitleContainsFold(st),
				music.ArtistContainsFold(st),
				music.AlbumContainsFold(st),
				// Add other searchable fields if needed (Genre, Year?)
			),
		)
	}

	// Get total count *before* applying limit/offset
	totalItemsCount, err := query.Count(ctx)
	if err != nil {
		return nil, fmt.Errorf("error getting total item count: %w", err)
	}

	// Apply pagination and fetch results
	entSongs, err := query.
		Limit(ps).
		Offset(offset).
		// OrderBy(ent.Asc(music.FieldTitle)) // Optional: Add default sorting
		All(ctx)

	if err != nil {
		return nil, fmt.Errorf("error fetching music list: %w", err)
	}

	// Map results
	songs := make([]*model.Song, len(entSongs))
	for i, entSong := range entSongs {
		songs[i] = mapEntMusicToModelSong(entSong)
	}

	return &model.MusicResponse{
		Songs:           songs,
		TotalItemsCount: totalItemsCount,
	}, nil
}

// Mutation returns MutationResolver implementation.
func (r *Resolver) Mutation() MutationResolver { return &mutationResolver{r} }

// Query returns QueryResolver implementation.
func (r *Resolver) Query() QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }

