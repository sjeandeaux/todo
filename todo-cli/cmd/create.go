package cmd

import (
	"context"
	"os"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/sjeandeaux/todo/pkg/client"
	"github.com/spf13/cobra"
)

var createArgs = &client.ToDo{}

var createCmd = &cobra.Command{
	Use:   `create --title=<title> --description=<description> --state=[NOT_STARTED, IN_PROGRESS, DONE] --tags="tag1,tag2" --reminder=<duration>`,
	Short: "Create a todo",
	Run: func(createCmd *cobra.Command, args []string) {
		client, err := cmdLine.client()
		if err != nil {
			log.Errorf("grpc client: %v\n", err)
			os.Exit(1)
		}
		ctx, cancel := context.WithTimeout(context.Background(), cmdLine.timeout)
		defer cancel()

		// Create the todo
		resp, err := client.Create(ctx, *createArgs)

		if err != nil {
			log.Errorf("grpc client: %v\n", err)
			os.Exit(1)
		}

		log.Infof("ID:%q", resp)

	},
}

func init() {
	createCmd.Flags().StringVarP(&createArgs.Title, "title", "", "title", "The title of todo")
	createCmd.Flags().StringVarP(&createArgs.Description, "description", "", "description", "The description of todo")
	createCmd.Flags().StringVarP(&createArgs.State, "state", "", "NOT_STARTED", "The state [NOT_STARTED, IN_PROGRESS, DONE]")
	createCmd.Flags().StringSliceVarP(&createArgs.Tags, "tags", "", []string{}, "The tags tag1,...,tagN")

	now := time.Now().UTC().Add(1 * time.Hour)
	createCmd.Flags().Int64VarP(&createArgs.Reminder, "reminder", "", now.Unix(), "The reminder")
}
