package music

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
	"webmane_go/graph/model"

	"github.com/jackc/pgx/v4/pgxpool"
)

const BaseDirectory = "./music/data/"

var Extensions = []string{".m4a", ".mp4", ".mp3", ".flac"}

func GetMusic(dbpool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := context.Background()
		// extract the relative file path and validate
		musicId := r.URL.Query().Get("id")
		if musicId == "" {
			http.Error(w, "Missing music ID", http.StatusBadRequest)
			return
		}

		var song model.Song
		var lastUpdateTime time.Time
		sql := `SELECT id, path, last_update, title, artist, album, genre, release_year, cover_art FROM MUSIC WHERE id = $1`
		err := dbpool.QueryRow(ctx, sql, musicId).Scan(&song.ID, &song.Path, &lastUpdateTime, &song.Title, &song.Artist, &song.Album, &song.Genre, &song.ReleaseYear, &song.CoverArt)
		if err != nil {
			http.Error(w, "Error getting music: "+err.Error(), http.StatusInternalServerError)
			return
		}

		song.LastUpdate = lastUpdateTime.Format(time.RFC3339)
		extension := filepath.Ext(song.Path)
		// Set the appropriate content type for FLAC audio
		switch {
		case extension == ".flac":
			w.Header().Set("Content-Type", "audio/flac")
		case extension == ".mp4":
			w.Header().Set("Content-Type", "audio/mp4")
		case extension == ".m4a":
			w.Header().Set("Content-Type", "audio/m4a")
		case extension == ".mp3":
			w.Header().Set("Content-Type", "audio/mp3")
		default:
			http.Error(w, "Unsupported file type", http.StatusUnsupportedMediaType)
			return
		}

		// Open the audio file
		file, err := os.Open(song.Path)
		if err != nil {
			http.Error(w, "File not found", http.StatusNotFound)
			return
		}
		// Get the file size
		stat, err := file.Stat()
		if err != nil {
			http.Error(w, "Could not get file info", http.StatusInternalServerError)
			return
		}
		fileSize := stat.Size()

		// Handle HTTP range requests
		rangeHeader := r.Header.Get("Range")
		if rangeHeader == "" {
			w.Header().Set("Content-Length", fmt.Sprint(fileSize))
			http.ServeContent(w, r, song.Path, stat.ModTime(), file)
			return
		}

		ranges := strings.Split(rangeHeader, "=")
		if len(ranges) != 2 {
			http.Error(w, "Invalid range header", http.StatusBadRequest)
			return
		}

		rangeParts := strings.Split(ranges[1], "-")
		start := int64(0)
		if rangeParts[0] != "" {
			start, err = strconv.ParseInt(rangeParts[0], 10, 64)
			if err != nil {
				http.Error(w, "Invalid range start", http.StatusBadRequest)
				return
			}
		}
		end := fileSize - 1
		if len(rangeParts) > 1 && rangeParts[1] != "" {
			end, err = strconv.ParseInt(rangeParts[1], 10, 64)
			if err != nil {
				http.Error(w, "Invalid range end", http.StatusBadRequest)
				return
			}
		}

		if start > end || start < 0 || end >= fileSize {
			http.Error(w, "Invalid range", http.StatusRequestedRangeNotSatisfiable)
			return
		}

		w.Header().Set("Content-Range", fmt.Sprintf("bytes %d-%d/%d", start, end, fileSize))
		w.Header().Set("Content-Length", fmt.Sprint(end-start+1))
		w.WriteHeader(http.StatusPartialContent)

		_, err = file.Seek(start, io.SeekStart)
		if err != nil {
			http.Error(w, "Error seeking file", http.StatusInternalServerError)
			return
		}

		_, err = io.CopyN(w, file, end-start+1)
		if err != nil {
			http.Error(w, "Error serving file", http.StatusInternalServerError)
		}
	}
}
