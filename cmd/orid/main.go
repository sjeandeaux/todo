package main

import (
	"context"
	"flag"

	log "github.com/sirupsen/logrus"

	"github.com/sjeandeaux/ori/pkg/config"
	"github.com/sjeandeaux/ori/pkg/grpc"
	"github.com/sjeandeaux/ori/pkg/information"
	"github.com/sjeandeaux/ori/pkg/service"
)

type commandLine struct {
	config service.Config
	host   string
	port   string
}

var cmdLine = &commandLine{}

func init() {
	flag.StringVar(&cmdLine.config.Host, "mongo-host", config.LookupEnvOrString("MONGO_HOST", ""), "The mongo host")
	flag.StringVar(&cmdLine.config.Login, "mongo-login", config.LookupEnvOrString("MONGO_LOGIN", "devroot"), "The mongo login")
	flag.StringVar(&cmdLine.config.Password, "mongo-password", config.LookupEnvOrString("MONGO_PASSWORD", "devroot"), "The mongo password (it should a secret but out of laziness...")
	flag.StringVar(&cmdLine.config.Port, "mongo-port", config.LookupEnvOrString("MONGO_POST", "27017"), "The mongo port")

	flag.StringVar(&cmdLine.config.Database, "mongo-database", config.LookupEnvOrString("MONGO_DATABASE", "challenge"), "The database")
	flag.StringVar(&cmdLine.config.Collection, "mongo-collection", config.LookupEnvOrString("MONGO_COLLECTION", "todo"), "The todo collection")

	flag.StringVar(&cmdLine.host, "host", config.LookupEnvOrString("HOST", "0.0.0.0"), "The grpc host")
	flag.StringVar(&cmdLine.port, "port", config.LookupEnvOrString("PORT", "8080"), "The grpc port")
	flag.Parse()
}

func main() {
	log.Println(information.Print())
	ctx := context.Background()
	todoService, err := service.NewToDoServiceServer(ctx, cmdLine.config)
	if err != nil {
		log.Fatal(err)
	}
	defer todoService.Close()

	log.Printf("Starting server GRPC on host:%q port:%q\n", cmdLine.host, cmdLine.port)
	pGRP, err := grpc.RunServer(ctx, cmdLine.host, cmdLine.port, todoService)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Started server GRPC on host:%q port:%d\n", cmdLine.host, pGRP)
	<-ctx.Done()
}
