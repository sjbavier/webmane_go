package cmd

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log" // Use log package for better output formatting
	"os"
	"path/filepath"
	"runtime"
	"sync"

	// Import base ent package for IsNotFound
	// Import ent/music package (aliased to avoid conflict)
	"webmane_go/graph"
	"webmane_go/graph/model"
	"webmane_go/music" // Keep for BaseDirectory and Extensions

	"github.com/spf13/cobra"
	ffmpeg_go "github.com/u2takey/ffmpeg-go"
)

// CommandContext remains the same
type CommandContext struct {
	Resolver *graph.Resolver
}

// CommandContextWrapper remains the same
func CommandContextWrapper(resolver *graph.Resolver) *cobra.Command {
	ctx := &CommandContext{Resolver: resolver}

	seedCmd := &cobra.Command{
		Use:   "seed",
		Short: "Seed database with music data from the ./music/data/ directory",
		Long: `Scans the ./music/data/ directory recursively for supported audio files (.m4a, .mp4, .mp3, .flac).
Uses ffprobe to extract metadata and ffmpeg to extract embedded cover art (if any).
Upserts the song information into the database using the GraphQL mutation.
If a song exists, existing metadata is preserved unless the scanned file provides non-empty data for a field.
Errors processing individual files will be logged, but the process will attempt to continue with other files.`, // Updated description
		RunE: func(cmd *cobra.Command, args []string) error {
			log.Println("Starting database seed process...")
			err := ctx.seedMusic()
			if err != nil {
				log.Printf("Seeding process finished with errors: %v", err)
				return err
			}
			log.Println("Seeding process completed.")
			return nil
		},
	}
	rootCmd.AddCommand(seedCmd)
	return rootCmd
}

// --- FFProbe structs remain the same ---
type FFProbe struct {
	Streams []Stream `json:"streams"`
	Format  Format   `json:"format"`
}
type Stream struct {
	Index         int         `json:"index"`
	CodecName     string      `json:"codec_name"`
	CodecLongName string      `json:"codec_long_name"`
	CodecType     string      `json:"codec_type"`
	SampleRate    string      `json:"sample_rate"`
	Channels      int         `json:"channels"`
	ChannelLayout string      `json:"channel_layout"`
	Duration      string      `json:"duration"`
	BitRate       string      `json:"bit_rate"`
	Disposition   Disposition `json:"disposition"`
}
type Disposition struct {
	Default         int `json:"default"`
	Dub             int `json:"dub"`
	Original        int `json:"original"`
	Comment         int `json:"comment"`
	Lyrics          int `json:"lyrics"`
	Karaoke         int `json:"karaoke"`
	Forced          int `json:"forced"`
	HearingImpaired int `json:"hearing_impaired"`
	VisualImpaired  int `json:"visual_impaired"`
	CleanEffects    int `json:"clean_effects"`
	AttachedPic     int `json:"attached_pic"`
	TimedThumbnails int `json:"timed_thumbnails"`
	NonDiegetic     int `json:"non_diegetic"`
	Captions        int `json:"captions"`
	Descriptions    int `json:"descriptions"`
	Metadata        int `json:"metadata"`
	Dependent       int `json:"dependent"`
	StillImage      int `json:"still_image"`
}
type Format struct {
	Filename   string `json:"filename"`
	NbStreams  int    `json:"nb_streams"`
	FormatName string `json:"format_name"`
	Duration   string `json:"duration"`
	Size       string `json:"size"`
	BitRate    string `json:"bit_rate"`
	Tags       Tag    `json:"tags"`
}
type Tag struct {
	Title       string `json:"title,omitempty"`
	Artist      string `json:"artist,omitempty"`
	Album       string `json:"album,omitempty"`
	Genre       string `json:"genre,omitempty"`
	ReleaseYear string `json:"date,omitempty"` // Corresponds to "date" tag in ffmpeg
}

// seedMusic function remains largely the same, orchestrating the walk and goroutines
func (ctx *CommandContext) seedMusic() error {
	var wg sync.WaitGroup
	numWorkers := runtime.NumCPU()
	if numWorkers > 8 {
		numWorkers = 8
	}
	log.Printf("Using %d workers for seeding", numWorkers)
	semaphore := make(chan struct{}, numWorkers)
	errorChan := make(chan error, 1)
	var firstProcessingError error

	log.Println("Scanning music directory...")
	errWalk := filepath.Walk(music.BaseDirectory, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Printf("Warning: Error accessing path %q: %v. Skipping.", path, err)
			return nil
		}
		if info.IsDir() {
			return nil
		}

		// only support certain filetypes
		extension := filepath.Ext(path)
		isSupported := false
		for _, ext := range music.Extensions {
			if extension == ext {
				isSupported = true
				break
			}
		}

		if isSupported {
			wg.Add(1)
			semaphore <- struct{}{}

			go func(filePath string) {
				defer wg.Done()
				defer func() { <-semaphore }()

				// *** Call the modified insertSong ***
				err := insertSongWithAdditiveLogic(filePath, ctx) // Renamed for clarity
				if err != nil {
					log.Printf("ERROR processing song %s: %v", filePath, err)
					select {
					case errorChan <- err:
					default:
					}
				}
			}(path)
		}
		return nil
	})

	wg.Wait()
	close(errorChan)

	if errWalk != nil {
		log.Printf("Error during directory walk: %v", errWalk)
		return fmt.Errorf("directory walk failed: %w", errWalk)
	}

	firstProcessingError = <-errorChan

	if firstProcessingError != nil {
		log.Println("Finished walking directory, but encountered errors processing one or more files.")
		return firstProcessingError
	}

	log.Println("Finished seeding walk successfully.")
	return nil
}

// *** New insertSong function specific to seeding with additive logic ***
func insertSongWithAdditiveLogic(path string, cmdCtx *CommandContext) error {
	opCtx := context.Background() // Context for database operations

	metaJson, errProbe := ffmpeg_go.Probe(path)
	if errProbe != nil {
		// Log the error and continue with minimal/default metadata.
		log.Printf("Warning: ffprobe failed for %s, inserting record with minimal metadata: %v", path, errProbe)
		metaJson = "{}" // Use an empty JSON so that unmarshalling works.
	}

	var songMeta FFProbe
	errMar := json.Unmarshal([]byte(metaJson), &songMeta)
	if errMar != nil {
		log.Printf("Warning: failed to unmarshal ffprobe JSON for %s: %v", path, errMar)
		// Continue with an empty songMeta
		songMeta = FFProbe{}
	}

	// Store extracted data
	newTitle := songMeta.Format.Tags.Title
	newArtist := songMeta.Format.Tags.Artist
	newAlbum := songMeta.Format.Tags.Album
	newGenre := songMeta.Format.Tags.Genre
	newReleaseYear := songMeta.Format.Tags.ReleaseYear

	// 2. Extract cover art (ffmpeg) - Allow failure here
	hasCoverArt := false
	for _, stream := range songMeta.Streams {
		if stream.CodecType == "video" && stream.Disposition.AttachedPic == 1 {
			hasCoverArt = true
			break
		}
	}
	var encodedCoverArt string
	if hasCoverArt {
		// Extract cover art directly into memory
		var buf bytes.Buffer
		errExtract := ffmpeg_go.Input(path).
			Output("pipe:1", ffmpeg_go.KwArgs{"map": "0:v?", "frames:v": 1, "f": "mjpeg"}).
			WithOutput(&buf, os.Stdout).
			Run()
		if errExtract != nil {
			fmt.Printf("error extracting cover art: %v\n", errExtract)
		} else {
			// Read and encode cover art data to base64
			if buf.Len() > 0 {
				encodedCoverArt = base64.StdEncoding.EncodeToString(buf.Bytes())
			}
		}
	}

	// // 4. Construct the input for the UpsertSong resolver based on existence and rules
	input := model.SongInput{
		Path: path, // Path is always required
		// Initialize pointers to nil
		Title:       nil,
		Artist:      nil,
		Album:       nil,
		Genre:       nil,
		ReleaseYear: nil,
		CoverArt:    nil,
	}
	if newTitle != "" {
		input.Title = &newTitle
	}
	if newArtist != "" {
		input.Artist = &newArtist
	}
	if newAlbum != "" {
		input.Album = &newAlbum
	}
	if newGenre != "" {
		input.Genre = &newGenre
	}
	if newReleaseYear != "" {
		input.ReleaseYear = &newReleaseYear
	}
	if encodedCoverArt != "" {
		input.CoverArt = &encodedCoverArt
	}

	// 5. Call the standard AdditiveUpsertSong resolver with the prepared input
	_, upsertErr := cmdCtx.Resolver.Mutation().AdditivePathUpsertSong(opCtx, input)

	if upsertErr != nil {
		// The resolver handles the actual DB interaction (create or update)
		return fmt.Errorf("failed to upsert song %s via resolver: %w", path, upsertErr)
	}

	return nil // Indicate success for this song
}

// Note: The original insertSong function is effectively replaced by
// insertSongWithAdditiveLogic for the context of the seed command.
