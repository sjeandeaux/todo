package main

import (
	"flag"

	log "github.com/sirupsen/logrus"

	"github.com/sjeandeaux/ori/pkg/information"
)

type commandLine struct {
	host string
	port string
}

var cmdLine = &commandLine{}

func init() {
	flag.StringVar(&cmdLine.host, "host", "localhost", "The grpc host")
	flag.StringVar(&cmdLine.port, "port", "8080", "The grpc port")
	flag.Parse()
}

func main() {
	log.Println(information.Print())
}
