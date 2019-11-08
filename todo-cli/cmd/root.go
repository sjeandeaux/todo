package cmd

import (
	"net"
	"time"

	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	log "github.com/sirupsen/logrus"
	"github.com/sjeandeaux/todo/pkg/client"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	grpc "google.golang.org/grpc"
)

type commandLine struct {
	logLevel string
	host     string
	port     string

	timeout time.Duration
}

//the client on todo manager
func (c *commandLine) client() (*client.ToDoManager, error) {
	cc, err := grpc.Dial(net.JoinHostPort(c.host, c.port), grpc.WithInsecure(), grpc.WithUnaryInterceptor(grpc_prometheus.UnaryClientInterceptor),
		grpc.WithStreamInterceptor(grpc_prometheus.StreamClientInterceptor))
	if err != nil {
		return nil, err
	}
	return client.NewToDoManager(cc), nil
}

var cmdLine = &commandLine{}

var rootCmd = &cobra.Command{
	Use:   "todo-cli (create | read | update | delete | search)",
	Short: "A client to manage your todos list.",
	Long:  "A client which communicates with the daemon todod",
}

// Execute executes the root command.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	cobra.OnInitialize(initConfig)
	log.SetFormatter(&log.TextFormatter{})

	rootCmd.Flags().StringVarP(&cmdLine.logLevel, "log-level", "l", log.InfoLevel.String(), "The log level")
	rootCmd.PersistentFlags().StringVarP(&cmdLine.port, "port", "p", "8080", "The port")
	rootCmd.PersistentFlags().DurationVarP(&cmdLine.timeout, "timeout", "t", 3*time.Second, "The timeout when it calls the daemon")
	rootCmd.PersistentFlags().StringVarP(&cmdLine.host, "host", "o", "localhost", "The host")

	viper.BindPFlag("log-level", rootCmd.PersistentFlags().Lookup("log-level"))
	viper.BindPFlag("host", rootCmd.PersistentFlags().Lookup("host"))
	viper.BindPFlag("port", rootCmd.PersistentFlags().Lookup("port"))

	rootCmd.AddCommand(createCmd)
	rootCmd.AddCommand(readCmd)
	rootCmd.AddCommand(updateCmd)
	rootCmd.AddCommand(deleteCmd)
	rootCmd.AddCommand(searchCmd)
	rootCmd.AddCommand(versionCmd)

}

func initConfig() {
	if l, err := log.ParseLevel(cmdLine.logLevel); err == nil {
		log.SetLevel(l)
	} else {
		log.SetLevel(log.InfoLevel)
	}
	viper.AutomaticEnv()
}
