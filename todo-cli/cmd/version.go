package cmd

import (
	log "github.com/sirupsen/logrus"
	"github.com/sjeandeaux/todo/pkg/information"
	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   `version`,
	Short: "version a of the application",
	Run: func(createCmd *cobra.Command, args []string) {
		log.WithField("data", information.MetaDataValue).Infoln("MetaData")
	},
}
