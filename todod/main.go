package main

import (
	"context"
	"flag"
	"os"

	log "github.com/sirupsen/logrus"

	"github.com/sjeandeaux/todo/pkg/config"
	"github.com/sjeandeaux/todo/pkg/grpc"
	"github.com/sjeandeaux/todo/pkg/http"
	"github.com/sjeandeaux/todo/pkg/information"
	"github.com/sjeandeaux/todo/pkg/service"
)

type commandLine struct {
	url      string
	logLevel string
	host     string
	grpcPort string
	httpPort string
}

var cmdLine = &commandLine{}

func init() {
	flag.StringVar(&cmdLine.url, "mongo-url", config.LookupEnvOrString("MONGO_URL", "mongodb://localhost:27017@devroot:devroot/?authSource=admin"), "The mongo host")

	flag.StringVar(&cmdLine.host, "host", config.LookupEnvOrString("HOST", "0.0.0.0"), "The grpc host")
	flag.StringVar(&cmdLine.grpcPort, "grpc-port", config.LookupEnvOrString("GRPC_PORT", "8080"), "The grpc port")
	flag.StringVar(&cmdLine.httpPort, "http-port", config.LookupEnvOrString("HTTP_PORT", "8081"), "The http port promotheus or golang debug")
	flag.StringVar(&cmdLine.logLevel, "log-level", config.LookupEnvOrString("LOG_LEVEL", log.InfoLevel.String()), "Log level")
	flag.Parse()
}

func main() {
	if l, err := log.ParseLevel(cmdLine.logLevel); err == nil {
		log.SetLevel(l)
	} else {
		log.SetLevel(log.InfoLevel)
	}

	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)

	log.WithField("data", information.MetaDataValue).Infoln("MetaData")

	ctx := context.Background()
	todoService, err := service.NewToDoServiceServer(ctx, cmdLine.url)
	if err != nil {
		log.Fatal(err)
	}
	defer todoService.Close()

	log.Infof("Starting server GRPC on host:%q port:%q\n", cmdLine.host, cmdLine.grpcPort)
	pGRP, err := grpc.RunServer(ctx, cmdLine.host, cmdLine.grpcPort, todoService)
	if err != nil {
		log.Fatal(err)
	}
	log.Infof("Started server GRPC on host:%q port:%d\n", cmdLine.host, pGRP)

	pHTTP, err := http.RunServer(ctx, cmdLine.host, cmdLine.httpPort)
	if err != nil {
		log.Fatal(err)
	}
	log.Infof("Started server HTTP on host:%q port:%d\n", cmdLine.host, pHTTP)
	<-ctx.Done()
}
