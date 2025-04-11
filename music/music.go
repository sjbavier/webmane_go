package music

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	// "time" // No longer needed directly here
	"webmane_go/ent"       // Import ent
	"webmane_go/ent/music" // Import ent music package
	// "webmane_go/graph/model" // No longer needed here
	// "github.com/jackc/pgx/v4/pgxpool" // Remove pgxpool import
	// Needed for IsNotFound check
)

const BaseDirectory = "./music/data/" // Keep this if relevant for other parts

var Extensions = []string{".m4a", ".mp4", ".mp3", ".flac"} // Keep this if relevant

// GetMusic now accepts an *ent.Client
func GetMusic(client *ent.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := context.Background()
		// extract the music ID and validate
		musicIdStr := r.URL.Query().Get("id")
		if musicIdStr == "" {
			http.Error(w, "Missing music ID", http.StatusBadRequest)
			return
		}

		// Convert ID string to int for Ent
		musicId, err := strconv.Atoi(musicIdStr)
		if err != nil {
			http.Error(w, "Invalid music ID format", http.StatusBadRequest)
			return
		}

		// Fetch the music record using Ent client
		// We only strictly need the path here, but fetching the whole record is fine.
		entMusic, err := client.Music.
			Query().
			Where(music.ID(musicId)). // Filter by ID
			Only(ctx)                 // Expect exactly one result

		if err != nil {
				// Handle other potential database errors
			http.Error(w, "Error getting music: "+err.Error(), http.StatusInternalServerError)
			return
		}

		// Use the path from the fetched Ent record
		filePath := entMusic.Path
		extension := filepath.Ext(filePath)

		// Set the appropriate content type
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
			// This case might indicate bad data if the extension was validated on insert
			log.Printf("Warning: Serving file with unsupported extension '%s' for path: %s", extension, filePath)
			// Fallback or error - deciding based on requirements. Let's try application/octet-stream
			w.Header().Set("Content-Type", "application/octet-stream")
			// Or uncomment below to return an error:
			// http.Error(w, "Unsupported file type in database record", http.StatusUnsupportedMediaType)
			// return
		}

		// Open the audio file using the path from the database
		file, err := os.Open(filePath)
		if err != nil {
			// Log the specific error for debugging
			log.Printf("Error opening file %s: %v", filePath, err)
			// Check if the error is specifically "file not found"
			if os.IsNotExist(err) {
				http.Error(w, "File not found on server", http.StatusNotFound)
			} else {
				http.Error(w, "Could not open file", http.StatusInternalServerError)
			}
			return
		}
		defer file.Close() // Ensure file is closed

		// Get the file size
		stat, err := file.Stat()
		if err != nil {
			http.Error(w, "Could not get file info", http.StatusInternalServerError)
			return
		}
		fileSize := stat.Size()

		// --- Range request handling (remains the same) ---
		rangeHeader := r.Header.Get("Range")
		if rangeHeader == "" {
			w.Header().Set("Content-Length", fmt.Sprint(fileSize))
			// Use http.ServeContent for simplicity when not handling ranges explicitly
			// It handles caching headers (like ETag, Last-Modified) automatically
			http.ServeContent(w, r, filepath.Base(filePath), stat.ModTime(), file)
			return
		}

		ranges := strings.Split(rangeHeader, "=")
		if len(ranges) != 2 || ranges[0] != "bytes" {
			http.Error(w, "Invalid range header format", http.StatusBadRequest)
			return
		}

		rangeParts := strings.Split(ranges[1], "-")
		startStr := rangeParts[0]
		endStr := ""
		if len(rangeParts) > 1 {
			endStr = rangeParts[1]
		}

		start := int64(0)
		if startStr != "" {
			start, err = strconv.ParseInt(startStr, 10, 64)
			if err != nil || start < 0 {
				http.Error(w, "Invalid range start", http.StatusBadRequest)
				return
			}
		}

		end := fileSize - 1
		if endStr != "" {
			end, err = strconv.ParseInt(endStr, 10, 64)
			if err != nil || end < start || end >= fileSize {
				http.Error(w, "Invalid range end", http.StatusBadRequest)
				return
			}
		}

		// Check if the range is satisfiable
		if start >= fileSize {
			w.Header().Set("Content-Range", fmt.Sprintf("bytes */%d", fileSize))
			http.Error(w, "Range not satisfiable", http.StatusRequestedRangeNotSatisfiable)
			return
		}

		w.Header().Set("Content-Range", fmt.Sprintf("bytes %d-%d/%d", start, end, fileSize))
		w.Header().Set("Content-Length", fmt.Sprint(end-start+1))
		w.WriteHeader(http.StatusPartialContent)

		_, err = file.Seek(start, io.SeekStart)
		if err != nil {
			// Log error, but might be too late to send HTTP error header
			log.Printf("Error seeking file %s to %d: %v", filePath, start, err)
			// Attempt to send error if possible, otherwise connection might just close
			http.Error(w, "Error seeking file", http.StatusInternalServerError)
			return
		}

		// Copy the requested range
		written, err := io.CopyN(w, file, end-start+1)
		if err != nil {
			// Log error, especially if written != expected length
			log.Printf("Error serving file range (%d bytes written) for %s: %v", written, filePath, err)
			// Can't reliably send HTTP error here if headers are already sent
		}
	}
}

// --- SeedMusic function ---
// This function likely needs refactoring too if CreateAllMusic/RefreshAllMusic
// interact with the DB. Assuming they might be part of the CLI now (cmd/seed.go),
// this HTTP endpoint might be redundant or needs updating to use the Ent client
// if it's still required.

// If this endpoint is still needed and should trigger the seeding:
// It should probably call methods on the Resolver or directly use the Ent client.
// For now, leaving it as is but noting it needs review/refactoring/removal.

// import "github.com/jackc/pgx/v4/pgxpool" // Keep if SeedMusic is kept temporarily

// func SeedMusic(pool *pgxpool.Pool) http.HandlerFunc { // Or change to accept *ent.Client
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		// ... existing logic ...
// 		// If CreateAllMusic/RefreshAllMusic need DB access, pass the client/pool
// 	}
// }
// func CreateAllMusic(w http.ResponseWriter, r *http.Request) { /* ... */ }
// func RefreshAllMusic(w http.ResponseWriter, r *http.Request) { /* ... */ }
