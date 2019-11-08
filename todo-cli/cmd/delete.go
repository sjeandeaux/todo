package cmd

import (
	"context"
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var idDelete string

var deleteCmd = &cobra.Command{
	Use:   `delete --id=<title>`,
	Short: "Delete a todo by ID",
	Run: func(createCmd *cobra.Command, args []string) {
		client, err := cmdLine.client()
		if err != nil {
			log.Errorf("grpc client: %v\n", err)
			os.Exit(1)
		}
		ctx, cancel := context.WithTimeout(context.Background(), cmdLine.timeout)
		defer cancel()

		// Create the todo
		resp, err := client.Delete(ctx, idDelete)

		if err != nil {
			log.Errorf("grpc client: %v\n", err)
			os.Exit(1)
		}

		log.Infof("Deleted:%t", resp)

	},
}

func init() {
	deleteCmd.Flags().StringVarP(&idDelete, "id", "", "", "The ID of todo")
}
