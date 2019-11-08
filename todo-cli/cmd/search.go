package cmd

import (
	"context"
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

type commandLineSearchArgs struct {
	pattern string
	tags    []string
	states  []string
}

var searchArgs = commandLineSearchArgs{}

var searchCmd = &cobra.Command{
	Use:   `search --pattern=<pattern> --tags=<tags> --states=<>`,
	Short: "Search todo",
	Run: func(createCmd *cobra.Command, args []string) {
		client, err := cmdLine.client()
		if err != nil {
			log.Errorf("grpc client: %v\n", err)
			os.Exit(1)
		}
		ctx, cancel := context.WithTimeout(context.Background(), cmdLine.timeout)
		defer cancel()

		// Create the todo
		resp, err := client.Search(ctx, searchArgs.pattern, searchArgs.tags, searchArgs.states)

		if err != nil {
			log.Errorf("grpc client: %v\n", err)
			os.Exit(1)
		}

		for _, todo := range resp {
			log.Infof("Todo:%v", todo)
		}
	},
}

func init() {
	searchCmd.Flags().StringVarP(&searchArgs.pattern, "pattern", "", "", "The pattern on description ex:.*")
	searchCmd.Flags().StringSliceVarP(&searchArgs.states, "states", "", []string{}, "The states NOT_STARTED, IN_PROGRESS or DONE")
	searchCmd.Flags().StringSliceVarP(&searchArgs.tags, "tags", "", []string{}, "The tags")
}
