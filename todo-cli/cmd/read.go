package cmd

import (
	"context"
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var idRead string

var readCmd = &cobra.Command{
	Use:   `read --id=<title>`,
	Short: "Read a todo by ID",
	Run: func(createCmd *cobra.Command, args []string) {
		client, err := cmdLine.client()
		if err != nil {
			log.Errorf("grpc client: %v\n", err)
			os.Exit(1)
		}
		ctx, cancel := context.WithTimeout(context.Background(), cmdLine.timeout)
		defer cancel()

		// Create the todo
		resp, err := client.Read(ctx, idRead)

		if err != nil {
			log.Errorf("grpc client: %v\n", err)
			os.Exit(1)
		}

		log.Infof("Todo:%v", resp)

	},
}

func init() {
	readCmd.Flags().StringVarP(&idRead, "id", "", "", "The ID of todo")
}
