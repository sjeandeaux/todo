package main

import (
	log "github.com/sirupsen/logrus"

	"github.com/sjeandeaux/todo/todo-cli/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		log.Panic(err)
	}
}
