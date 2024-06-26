package music

import (
	"errors"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
)

const BaseDirectory = "./music/data/"

var Extensions = []string{".m4a", ".mp4", ".mp3", ".flac"}

func GetMusic(w http.ResponseWriter, r *http.Request) {
	// extract the relative file path and validate
	fileQuery := r.URL.Query().Get("file")

	cleanPath := path.Clean(fileQuery)

	var fullPath string
	var extension string
	for _, ext := range Extensions {
		if _, err := os.Stat(BaseDirectory + cleanPath + ext); err != nil {
			if errors.Is(err, fs.ErrNotExist) {
				// path doesn't exist
				continue
			}
		}
		fullPath = filepath.Join(BaseDirectory, fileQuery+ext)
		extension = strings.Split(ext, ".")[1]
		break
	}

	// Set the appropriate content type for FLAC audio
	switch {
	case extension == "flac":
		w.Header().Set("Content-Type", "audio/flac")
	case extension == "mp4":
		w.Header().Set("Content-Type", "audio/mp4")
	case extension == "m4a":
		w.Header().Set("Content-Type", "audio/m4a")
	case extension == "mp3":
		w.Header().Set("Content-Type", "audio/mp3")
	default:
		w.Header().Set("Content-Type", fmt.Sprintf("audio/%v", extension))
	}

	// Open the FLAC file
	file, err := os.Open(fullPath)
	if err != nil {
		http.Error(w, "File not found", http.StatusNotFound)
		return
	}
	defer file.Close()

	// Stream the file to the response
	_, err = io.Copy(w, file)
	if err != nil {
		http.Error(w, "Error streaming file", http.StatusInternalServerError)
	}
}
