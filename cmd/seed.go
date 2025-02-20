package cmd

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"webmane_go/graph"
	"webmane_go/graph/model"
	"webmane_go/music"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/spf13/cobra"
	ffmpeg_go "github.com/u2takey/ffmpeg-go"
)

type CommandContext struct {
	DBPool   *pgxpool.Pool
	Resolver *graph.Resolver
}

func CommandContextWrapper(dbPool *pgxpool.Pool, resolver *graph.Resolver) *cobra.Command {
	ctx := &CommandContext{DBPool: dbPool, Resolver: resolver}

	seedCmd := &cobra.Command{
		Use:   "seed",
		Short: "seed database with data",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("attempting to seed database")
			return ctx.seedMusic()
		},
	}
	rootCmd.AddCommand(seedCmd)

	return rootCmd

}

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
	ReleaseYear string `json:"date,omitempty"`
}

func (ctx *CommandContext) seedMusic() error {
	// Use a wait group to wait for all goroutines to finish
	var wg sync.WaitGroup
	// Use a semaphore to limit the number of concurrent goroutines
	semaphore := make(chan struct{}, runtime.NumCPU())

	errWalk := filepath.Walk(music.BaseDirectory, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		// Skip directories
		if info.IsDir() {
			return nil
		}

		// grab file extension
		extension := filepath.Ext(path)
		for _, ext := range music.Extensions {
			// looping through the accepted extensions on match insert song, break loop
			if extension == ext {
				// Increment the wait group counter
				wg.Add(1)
				// Acquire a semaphore slot
				semaphore <- struct{}{}
				go func(path string) {
					defer wg.Done()
					// frees semaphore slot
					defer func() { <-semaphore }()
					insertSong(path, ctx)
				}(path)
				break
			}
		}

		return nil
	})

	if errWalk != nil {
		fmt.Printf("Error seeding song %v\n", errWalk)
		return errWalk
	}

	// Wait for all goroutines to finish
	wg.Wait()

	fmt.Println("Done seeding")
	return nil
}

func insertSong(path string, ctx *CommandContext) {
	metaJson, errProbe := ffmpeg_go.Probe(path)
	if errProbe != nil {
		fmt.Printf("error with probe: %v", errProbe)
	}

	var songMeta FFProbe
	errMar := json.Unmarshal([]byte(metaJson), &songMeta)
	if errMar != nil {

		fmt.Printf("error with unmarshal: %v", errMar)
	}
	// Check if the file has a video stream (cover art)
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

	input := model.SongInput{
		Path:        path,
		Title:       &songMeta.Format.Tags.Title,
		Artist:      &songMeta.Format.Tags.Artist,
		Album:       &songMeta.Format.Tags.Album,
		Genre:       &songMeta.Format.Tags.Genre,
		ReleaseYear: &songMeta.Format.Tags.ReleaseYear,
		CoverArt:    &encodedCoverArt,
	}

	_, err := ctx.Resolver.Mutation().UpsertSong(context.Background(), input)

	if err != nil {
		fmt.Printf("Error seeding song %v\n", path)
		fmt.Println(err)
	}
	// fmt.Printf("song seeded %v", _song)
	// fmt.Println("")

}
