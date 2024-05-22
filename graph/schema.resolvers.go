package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.44

import (
	"context"
	"fmt"
	"strconv"
	"time"
	"webmane_go/graph/model"
)

// UpsertSong is the resolver for the upsertSong field.
func (r *mutationResolver) UpsertSong(ctx context.Context, input model.SongInput) (*model.Song, error) {
	// panic(fmt.Errorf("not implemented: UpsertSong - upsertSong"))
	var song model.Song
	song.Path = input.Path
	song.LastUpdate = time.Now().Format(time.RFC3339)

	// Check and assign optional fields
	if input.Title != nil {
		song.Title = *input.Title
	}
	if input.Artist != nil {
		song.Artist = *input.Artist
	}
	if input.Album != nil {
		song.Album = *input.Album
	}
	if input.Genre != nil {
		song.Genre = *input.Genre
	}
	if input.ReleaseYear != nil {
		song.ReleaseYear = *input.ReleaseYear
	}

	if input.ID != nil && *input.ID != "" {
		// UPDATE song
		id, _ := strconv.ParseInt(*input.ID, 10, 64)
		sql := `UPDATE music SET path=$2, last_update=$3, title=$4, artist=$5, album=$6, genre=$7, release_year=$8 WHERE id=$1 RETURNING id`
		err := r.DB.QueryRow(ctx, sql, id, song.Path, song.LastUpdate, song.Title, song.Artist, song.Album, song.Genre, song.ReleaseYear).Scan(&song.ID)
		if err != nil {
			return nil, fmt.Errorf("error updating song: %v\nerror: %v", song.Title, err)
		}
	} else {
		// INSERT: no ID
		sql := `INSERT INTO music (path, last_update, title, artist, album, genre, release_year) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id`
		err := r.DB.QueryRow(ctx, sql, song.Path, song.LastUpdate, song.Title, song.Artist, song.Album, song.Genre, song.ReleaseYear).Scan(&song.ID)
		if err != nil {
			return nil, fmt.Errorf("error updating song: %v\nerror: %v", song.Title, err)
		}
	}

	return &song, nil
}

// Song is the resolver for the song field.
func (r *queryResolver) Song(ctx context.Context, id string) (*model.Song, error) {
	panic(fmt.Errorf("not implemented: Song - song"))
}

// Music is the resolver for the music field.
func (r *queryResolver) Music(ctx context.Context, pageNumber *int, pageSize *int) (*model.MusicResponse, error) {
	// panic(fmt.Errorf("not implemented: Music - music"))
	defaultPageNumber := 1
	defaultPageSize := 10

	if pageNumber == nil {
		pageNumber = &defaultPageNumber
	}
	if pageSize == nil {
		pageSize = &defaultPageSize
	}

	offset := (*pageNumber - 1) * *pageSize

	// Query to get the total item count
	countSQL := `SELECT COUNT(*) FROM MUSIC`
	var totalItemsCount int
	err := r.DB.QueryRow(ctx, countSQL).Scan(&totalItemsCount)
	if err != nil {
		return nil, fmt.Errorf("error getting total item count: %v", err)
	}

	sql := `SELECT * FROM MUSIC LIMIT $1 OFFSET $2`
	rows, err := r.DB.Query(ctx, sql, pageSize, offset)
	if err != nil {
		return nil, fmt.Errorf("error getting music")
	}
	defer rows.Close()
	var songs []*model.Song

	for rows.Next() {
		var song model.Song
		var lastUpdateTime time.Time
		err := rows.Scan(&song.ID, &song.Path, &lastUpdateTime, &song.Title, &song.Artist, &song.Album, &song.Genre, &song.ReleaseYear)

		if err != nil {
			return nil, fmt.Errorf("error scanning song %v", err)
		}

		song.LastUpdate = lastUpdateTime.Format(time.RFC3339)
		songs = append(songs, &song)
	}

	if err = rows.Err(); err != nil {

		return nil, fmt.Errorf("error iterating %v", err)
	}
	// Return the MusicResult
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
