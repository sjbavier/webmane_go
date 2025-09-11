package cmd

import (
	"flag"
	"fmt"
	"io"
	"log" // Use log package for better output formatting
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"webmane_go/graph"

	"github.com/kkdai/youtube/v2"
	"github.com/spf13/cobra"
)

func DlUrl(resolver *graph.Resolver) *cobra.Command {
	ctx := &CommandContext{Resolver: resolver}
	return &cobra.Command{
		Use:   "dl",
		Short: "download a video and convert to m4a",
		Long: "Downloads video, converts it to an m4a, adds it to the music library",
		RunE: func(cmd *cobra.Command, args []string) error {
			log.Println("Starting database seed process...")
			err := ctx.downloadUrl(args)
			if err != nil {
				log.Printf("Seeding process finished with errors: %v", err)
				return err
			}
			log.Println("Seeding process completed.")
			return nil
		},
	}
}


func (ctx *CommandContext) downloadUrl(args []string) error {
	// Parse command-line flags
	fmt.Printf("args: %v\n", args)
	url := flag.String("url", "", "https://www.youtube.com/watch?v=cD6qkQjTHq0")
	outArg := flag.String("out", "", "Output filename (e.g., audio.m4a). If empty, defaults to <video_title>.m4a")
	flag.Parse()

	if *url == "" {
		fmt.Fprintln(os.Stderr, "Usage: go run main.go -url <YouTubeURL> [-out <filename>]")
		os.Exit(1)
	}

	// Create a client
	client := youtube.Client{}

	// Fetch video metadata
	video, err := client.GetVideo(*url)
	if err != nil {
		log.Fatalf("Error fetching video info: %v\n", err)
		return err
	}

	// Determine output filename
	outputFilename := *outArg
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
	stream, _, err := client.GetStream(video, bestM4AFormat)
	if err != nil {
		log.Fatalf("Error getting stream: %v\n", err)
		return err
	}
	defer stream.Close()

	// Create output file
	file, err := os.Create(outputFilename)
	if err != nil {
		log.Fatalf("Error creating file: %v\n", err)
		return err
	}
	defer file.Close()

	// Copy stream into file
	n, err := io.Copy(file, stream)
	if err != nil {
		log.Fatalf("Error writing to file: %v\n", err)
		return err
	}

	fmt.Printf("Successfully downloaded %d bytes to %s\n", n, outputFilename)
	fmt.Printf("Format details:  MimeType: %s, Bitrate: %dkbps\n",
		bestM4AFormat.MimeType, bestM4AFormat.Bitrate/1000)
	
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
