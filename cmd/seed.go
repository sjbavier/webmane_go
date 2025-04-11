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
	"webmane_go/graph"
	"webmane_go/graph/model"
	"webmane_go/music" // Keep for BaseDirectory and Extensions

	// "github.com/jackc/pgx/v4/pgxpool" // Remove pgxpool import
	"github.com/spf13/cobra"
	ffmpeg_go "github.com/u2takey/ffmpeg-go"
)

// CommandContext now only needs the Resolver, which holds the EntClient
type CommandContext struct {
	Resolver *graph.Resolver
}

// CommandContextWrapper now only accepts the Resolver
func CommandContextWrapper(resolver *graph.Resolver) *cobra.Command {
	// ctx := &CommandContext{DBPool: dbPool, Resolver: resolver} // Old
	ctx := &CommandContext{Resolver: resolver} // New: Initialize with Resolver only

	seedCmd := &cobra.Command{
		Use:   "seed",
		Short: "Seed database with music data from the ./music/data/ directory",
		Long: `Scans the ./music/data/ directory recursively for supported audio files (.m4a, .mp4, .mp3, .flac).
Uses ffprobe to extract metadata and ffmpeg to extract embedded cover art (if any).
Upserts the song information into the database using the GraphQL mutation.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			log.Println("Starting database seed process...") // Use log
			err := ctx.seedMusic()
			if err != nil {
				log.Printf("Seeding process failed: %v", err) // Use log
			} else {
				log.Println("Seeding process completed successfully.") // Use log
			}
			return err // Return error for Cobra's handling
		},
	}
	rootCmd.AddCommand(seedCmd)

	// Add other commands to rootCmd here if needed

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

func (ctx *CommandContext) seedMusic() error {
	var wg sync.WaitGroup
	// Use a reasonable number of workers, not necessarily all CPUs if ffmpeg is heavy
	numWorkers := runtime.NumCPU()
	if numWorkers > 8 { // Cap workers to avoid overwhelming system/ffmpeg
		numWorkers = 8
	}
	log.Printf("Using %d workers for seeding", numWorkers)
	semaphore := make(chan struct{}, numWorkers)
	errorChan := make(chan error, 1) // Buffered channel to capture the first error

	errWalk := filepath.Walk(music.BaseDirectory, func(path string, info os.FileInfo, err error) error {
		// Check for errors from filepath.Walk itself
		if err != nil {
			log.Printf("Error accessing path %q: %v\n", path, err)
			// Decide if the error is critical enough to stop walking
			// return err // Uncomment to stop walking on access errors
			return nil // Continue walking even if some paths are inaccessible
		}
		// Skip directories
		if info.IsDir() {
			return nil
		}

		// Check if already received an error, stop adding new tasks
		select {
		case <-errorChan:
			return fmt.Errorf("stopping walk due to previous error") // Stop walking
		default:
			// Continue if no error received yet
		}

		// Check file extension
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
			semaphore <- struct{}{} // Acquire semaphore slot

			go func(filePath string) {
				defer wg.Done()
				defer func() { <-semaphore }() // Release semaphore slot

				// Check again if an error occurred while this goroutine was waiting
				select {
				case <-errorChan:
					return // Don't process if another goroutine failed
				default:
				}

				err := insertSong(filePath, ctx) // Pass the specific error back
				if err != nil {
					log.Printf("Error processing song %s: %v", filePath, err)
					// Try to send the error to the error channel.
					// If the channel is full (an error was already sent), this will do nothing.
					select {
					case errorChan <- err:
						// Signal that an error occurred
					default:
						// An error was already recorded
					}
				}
			}(path) // Pass the current path to the goroutine
		}
		return nil // Continue walking
	})

	// Wait for all started goroutines to finish
	wg.Wait()
	close(errorChan) // Close channel after all goroutines are done

	// Check if any error occurred during the walk or processing
	if errWalk != nil {
		log.Printf("Error during directory walk: %v\n", errWalk)
		return errWalk // Return the walk error
	}

	// Check if any goroutine reported an error
	if firstErr := <-errorChan; firstErr != nil {
		log.Println("Error occurred during song processing.")
		return firstErr // Return the first error encountered by a goroutine
	}

	log.Println("Finished seeding walk.") // Use log
	return nil
}

// insertSong now returns an error
func insertSong(path string, ctx *CommandContext) error {
	// Use context.Background() for now, but consider passing down a cancellable context
	opCtx := context.Background()

	metaJson, errProbe := ffmpeg_go.Probe(path) // Use context-aware version
	if errProbe != nil {
		// Return a more specific error
		return fmt.Errorf("ffprobe failed for %s: %w", path, errProbe)
	}

	var songMeta FFProbe
	errMar := json.Unmarshal([]byte(metaJson), &songMeta)
	if errMar != nil {
		// Return a more specific error
		return fmt.Errorf("failed to unmarshal ffprobe JSON for %s: %w", path, errMar)
	}

	// Check for cover art
	hasCoverArt := false
	videoStreamIndex := -1 // Store the index if needed
	for i, stream := range songMeta.Streams {
		// Check specifically for attached pictures which are often video streams
		if stream.CodecType == "video" && stream.Disposition.AttachedPic == 1 {
			hasCoverArt = true
			videoStreamIndex = i // Store index if needed for extraction mapping
			break
		}
	}

	var encodedCoverArt string
	if hasCoverArt {
		var buf bytes.Buffer
		// Explicitly map the video stream found
		outputArgs := ffmpeg_go.KwArgs{
			"map":      fmt.Sprintf("0:v:%d?", videoStreamIndex), // Map specific stream index if found, fallback if not (?)
			"frames:v": 1,                                        // Extract only one frame
			"f":        "mjpeg",                                   // Output as MJPEG image format
		}
		errExtract := ffmpeg_go.Input(path).
			Output("pipe:1", outputArgs).
			WithOutput(&buf, os.Stderr). // Pipe image data to buffer, ffmpeg logs to stderr
			Run()

		if errExtract != nil {
			// Log the error but don't necessarily stop the whole seed process for one cover art failure
			log.Printf("Warning: failed to extract cover art for %s: %v", path, errExtract)
			// encodedCoverArt remains empty
		} else if buf.Len() > 0 {
			encodedCoverArt = base64.StdEncoding.EncodeToString(buf.Bytes())
		} else {
			log.Printf("Warning: extracted cover art buffer is empty for %s", path)
		}
	}

	// Prepare input for the GraphQL mutation
	// Handle potentially empty strings from metadata gracefully
	title := songMeta.Format.Tags.Title
	artist := songMeta.Format.Tags.Artist
	album := songMeta.Format.Tags.Album
	genre := songMeta.Format.Tags.Genre
	releaseYear := songMeta.Format.Tags.ReleaseYear
	coverArtPtr := &encodedCoverArt
	if encodedCoverArt == "" {
		coverArtPtr = nil // Use nil if no cover art was extracted
	}

	input := model.SongInput{
		Path:        path,
		Title:       &title, // Pass pointers directly
		Artist:      &artist,
		Album:       &album,
		Genre:       &genre,
		ReleaseYear: &releaseYear,
		CoverArt:    coverArtPtr,
	}

	// Call the UpsertSong mutation via the resolver
	// The resolver internally uses the Ent client now.
	_, err := ctx.Resolver.Mutation().UpsertSong(opCtx, input)
	if err != nil {
		// Return a specific error for this song
		return fmt.Errorf("failed to upsert song %s via resolver: %w", path, err)
	}

	// Optional: Log success for this specific file
	// log.Printf("Successfully seeded: %s", path)
	return nil // Indicate success for this song
}

