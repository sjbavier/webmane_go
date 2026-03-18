package cmd

import (
	"fmt"
	"io"
	"log" // Use log package for better output formatting
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
	"webmane_go/graph"

	"github.com/kkdai/youtube/v2"
	"github.com/schollz/progressbar/v3"
	"github.com/spf13/cobra"
)

func DlUrl(resolver *graph.Resolver) *cobra.Command {
	ctx := &CommandContext{Resolver: resolver}
	var url, outArg string

	cmd := &cobra.Command{
		Use:   "dl",
		Short: "download a video and convert to m4a",
		Long:  "Downloads video, converts it to an m4a, adds it to the music library",
		RunE: func(cmd *cobra.Command, args []string) error {
			if url == "" {
				return fmt.Errorf("the --url flag is required")
			}

			log.Println("Starting database seed process...")
			// Pass the flag values to the download function
			err := ctx.DownloadUrl(url, outArg)
			if err != nil {
				// The error is already logged in downloadUrl, just return it
				log.Printf("Seeding process finished with errors: %v", err)
				return err
			}
			log.Println("Seeding process completed.")
			return nil
		},
	}

	// Define flags and attach them to the command
	cmd.Flags().StringVar(&url, "url", "", "YouTube URL to download (e.g., https://www.youtube.com/watch?v=videoID)")
	cmd.Flags().StringVar(&outArg, "out", "", "Output filename (e.g., audio.m4a). If empty, defaults to <video_title>.m4a")
	cmd.MarkFlagRequired("url") // Make the --url flag mandatory

	return cmd
}

func (ctx *CommandContext) DownloadUrl(url, outArg string) error {
	// Create a client
	client := youtube.Client{}

	// Fetch video metadata
	video, err := client.GetVideo(url)
	if err != nil {
		log.Fatalf("Error fetching video info: %v\n", err)
		return err
	}

	// Determine output filename
	outputFilename := outArg
	outputPath := os.Getenv("MUSIC_DATA")
	if outputFilename == "" {
		if video.Title != "" {
			sanitizedTitle := sanitizeFilename(video.Title)
			outputFilename = sanitizedTitle + ".m4a"
		} else {
			outputFilename = "audio.m4a" // Fallback if video title is unexpectedly empty
		}
	} else {
		// Ensure the provided output filename has the .m4a extension
		if strings.ToLower(filepath.Ext(outputFilename)) != ".m4a" {
			outputFilename += ".m4a"
		}
	}

	// Select the best M4A audio format
	formats := video.Formats.Type("audio")
	if len(formats) == 0 {
		log.Fatalf("No audio formats found for this video.\n")
		return err
	}

	var bestM4AFormat *youtube.Format
	maxBitrate := -1 // Initialize to ensure any valid bitrate is chosen

	for i := range formats {
		f := &formats[i] // Work with a pointer to the format struct
		// Check if the MimeType indicates an M4A container (audio/mp4)
		if strings.Contains(strings.ToLower(f.MimeType), "audio/mp4") {
			if f.Bitrate > maxBitrate {
				maxBitrate = f.Bitrate
				bestM4AFormat = f
			}
		}
	}

	if bestM4AFormat == nil {
		log.Printf("No M4A (audio/mp4) audio format found for this video.")
		log.Println("Available audio formats:")
		for _, f := range formats {
			log.Printf("  -  MimeType: %s, Bitrate: %dkbps, Quality: %s\n", f.MimeType, f.Bitrate/1000, f.AudioQuality)
		}
		os.Exit(1)
	}

	// Open stream
	stream, size, err := client.GetStream(video, bestM4AFormat)
	if err != nil {
		log.Fatalf("Error getting stream: %v\n", err)
		return err
	}
	defer stream.Close()

	// Create output file
	file, err := os.Create(outputPath + outputFilename)
	if err != nil {
		log.Fatalf("Error creating file: %v\n", err)
		return err
	}
	defer file.Close()

	// Create a new progress bar
	bar := progressbar.NewOptions64(
		size,
		progressbar.OptionSetDescription("Downloading"),
		progressbar.OptionSetWriter(os.Stderr),
		progressbar.OptionShowBytes(true),
		progressbar.OptionSetWidth(10),
		progressbar.OptionThrottle(65*time.Millisecond),
		progressbar.OptionShowCount(),
		progressbar.OptionOnCompletion(func() {
			fmt.Fprint(os.Stderr, "\n")
		}),
		progressbar.OptionSpinnerType(14),
		progressbar.OptionFullWidth(),
	)
	bar.RenderBlank()

	// Copy stream into file
	_, err = io.Copy(io.MultiWriter(file, bar), stream)
	if err != nil {
		log.Fatalf("Error writing to file: %v\n", err)
		return err
	}

	fmt.Printf("Successfully downloaded to %s\n", outputFilename)
	fmt.Printf("Format details:  MimeType: %s, Bitrate: %dkbps\n",
		bestM4AFormat.MimeType, bestM4AFormat.Bitrate/1000)
	
	// insert song in db
	err = InsertSongWithAdditiveLogic(outputPath + outputFilename, ctx)
	if err != nil {
		log.Fatalf("Error Inserting song to database %v\n", err)
		return err
	}

	return nil
}

// sanitizeFilename creates a safe filename from a given string.
// It replaces invalid characters with underscores, collapses multiple spaces,
// trims whitespace, and limits the filename length.
func sanitizeFilename(name string) string {
	// Replace invalid filesystem characters with an underscore
	name = regexp.MustCompile(`[<>:"/\\|?*\x00-\x1F]`).ReplaceAllString(name, "_")
	// Replace multiple spaces with a single space
	name = regexp.MustCompile(`\s+`).ReplaceAllString(name, " ")
	// Trim leading/trailing spaces
	name = strings.TrimSpace(name)
	if name == "" {
		return "untitled"
	}
	// Limit length to prevent issues with maximum filename lengths on some filesystems
	const maxFilenameLength = 200
	if len(name) > maxFilenameLength {
		name = name[:maxFilenameLength]
	}
	return name
}
