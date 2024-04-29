package cmd

import (
	"context"
	"fmt"
	"webmane_go/graph"
	"webmane_go/graph/model"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/spf13/cobra"
)

type CommandContext struct {
	DBPool           *pgxpool.Pool
	Resolver         *graph.Resolver
	MutationResolver *graph.MutationResolver
}

func CommandContextWrapper(dbPool *pgxpool.Pool, mutationResolver *graph.MutationResolver) *cobra.Command {
	ctx := &CommandContext{DBPool: dbPool, MutationResolver: mutationResolver}

	seedCmd := &cobra.Command{
		Use:   "seed",
		Short: "seed database with data",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("Database seeded")
			return ctx.seedMusic()
		},
	}
	rootCmd.AddCommand(seedCmd)

	return rootCmd

}

func (ctx *CommandContext) seedMusic() error {
	input := model.SongInput{
		ID:   "1",
		Path: "BEST OF AFROBEATS NAIJA OVERDOSE 13 VIDEO MIX 2022 [Burna Boy, Asake, Ruger, Buga, Finesse, Ckay]-GZOV93NoXSI.m4a",
	}
	song, err := ctx.MutationResolver.UpsertSong(context.Background(), input)
	if err != nil {
		fmt.Printf("Error seeding song %v\n", err)
		return err
	}
	fmt.Printf("Seeded song with ID %s\n", song.Path)
	return nil
}
