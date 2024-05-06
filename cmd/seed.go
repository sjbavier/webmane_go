package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
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

	errWalk := filepath.Walk(music.BaseDirectory, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		// Skip directories
		if info.IsDir() {
			return nil
		}
		// data, err := ffprobe.GetProbeData()
		metaJson, errProbe := ffmpeg_go.Probe(path)
		if errProbe != nil {
			fmt.Printf("error with probe: %v", errProbe)
		}

		var songMeta FFProbe
		errMar := json.Unmarshal([]byte(metaJson), &songMeta)
		if errMar != nil {

			fmt.Printf("error with unmarshal: %v", errMar)
		}
		input := model.SongInput{
			Path:        path,
			Title:       &songMeta.Format.Tags.Title,
			Artist:      &songMeta.Format.Tags.Artist,
			Album:       &songMeta.Format.Tags.Album,
			Genre:       &songMeta.Format.Tags.Genre,
			ReleaseYear: &songMeta.Format.Tags.ReleaseYear,
		}

		_song, err := ctx.Resolver.Mutation().UpsertSong(context.Background(), input)

		if err != nil {
			fmt.Printf("Error seeding song %v\n", err)
			return err
		}
		fmt.Printf("song seeded %v", _song)
		fmt.Printf("song seeded meta %v", songMeta)

		return nil
	})

	if errWalk != nil {
		fmt.Printf("Error seeding song %v\n", errWalk)
		return errWalk

	}

	// input := model.SongInput{
	// 	Path: "4BEST OF AFROBEATS NAIJA OVERDOSE 13 VIDEO MIX 2022 [Burna Boy, Asake, Ruger, Buga, Finesse, Ckay]-GZOV93NoXSI.m4a",
	// }

	// song, err := ctx.Resolver.Mutation().UpsertSong(context.Background(), input)

	// if err != nil {
	// 	fmt.Printf("Error seeding song %v\n", err)
	// 	return err
	// }

	// fmt.Println("Seeded song with path: ")
	// fmt.Printf("%s\n", song.Path)
	return nil
}
