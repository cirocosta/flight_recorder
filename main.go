package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/cirocosta/concourse_db_exporter/db"
	"github.com/cirocosta/concourse_db_exporter/exporter"
	"github.com/concourse/flag"
	"github.com/jessevdk/go-flags"
)

var (
	version = "dev"
	config  = struct {
		TelemetryPath string `long:"path" default:"/" description:"path to serve metrics"`
		ListenAddress string `long:"address" default:":9000" description:"address to listen for prometheus scraping"`

		Postgres flag.PostgresConfig `group:"PostgreSQL Configuration" namespace:"postgres"`
	}{}
)

func handleSignals(exp *exporter.Exporter) {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGTERM, syscall.SIGINT)
	<-sigChan

	exp.Stop()

}

func main() {
	_, err := flags.Parse(&config)
	if err != nil {
		os.Exit(1)
	}

	exp := exporter.Exporter{
		TelemetryPath: config.TelemetryPath,
		ListenAddress: config.ListenAddress,
	}

	go handleSignals(&exp)

	_, err = db.New(config.Postgres.ConnectionString())
	if err != nil {
		panic(err)
	}

	err = exp.Listen()
	if err != nil {
		panic(err)
	}
}
