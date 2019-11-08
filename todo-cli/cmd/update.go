package cmd

import (
	"context"
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/sjeandeaux/todo/pkg/client"
	"github.com/spf13/cobra"
)

var updateArgs = &client.ToDo{}

var updateCmd = &cobra.Command{
	Use:   `update --id=<id> --title=<title> --description=<description> --state=[NOT_STARTED, IN_PROGRESS, DONE] --tags="tag1,tag2" --reminder=<duration>`,
	Short: "Update a todo",
	Run: func(createCmd *cobra.Command, args []string) {
		client, err := cmdLine.client()
		if err != nil {
			log.Errorf("grpc client: %v\n", err)
			os.Exit(1)
		}
		ctx, cancel := context.WithTimeout(context.Background(), cmdLine.timeout)
		defer cancel()

		// Create the todo
		resp, err := client.Update(ctx, *updateArgs)

		if err != nil {
			log.Errorf("grpc client: %v\n", err)
			os.Exit(1)
		}

		log.Infof("Updated:%t", resp)

	},
}

func init() {
	updateCmd.Flags().StringVarP(&updateArgs.ID, "id", "", "", "The ID of todo")
	updateCmd.Flags().StringVarP(&updateArgs.Title, "title", "", "", "The title of todo")
	updateCmd.Flags().StringVarP(&updateArgs.Description, "description", "", "", "The description of todo")
	updateCmd.Flags().StringVarP(&updateArgs.State, "state", "", "", "The state [NOT_STARTED, IN_PROGRESS, DONE]")
	updateCmd.Flags().StringSliceVarP(&updateArgs.Tags, "tags", "", []string{}, "The tags tag1,...,tagN")
	updateCmd.Flags().Int64VarP(&updateArgs.Reminder, "reminder", "", 0, "The reminder")
}
