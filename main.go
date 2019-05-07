package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/cirocosta/flight_recorder/collectors"
	"github.com/cirocosta/flight_recorder/db"
	"github.com/cirocosta/flight_recorder/exporter"
	"github.com/concourse/flag"
	"github.com/jessevdk/go-flags"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/vito/twentythousandtonnesofcrudeoil"
)

var (
	version = "dev"
	config  = struct {
		TelemetryPath string `long:"path" default:"/" description:"path to serve metrics"`
		ListenAddress string `long:"address" default:":9000" description:"address to listen for prometheus scraping"`

		Postgres flag.PostgresConfig `group:"PostgreSQL Configuration" namespace:"postgres"`
	}{}
)

func handleSignals(exp *exporter.Exporter, database *db.Db) {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGTERM, syscall.SIGINT)
	<-sigChan

	exp.Stop()
	database.Close()
}

func main() {
	parser := flags.NewParser(&config, flags.HelpFlag|flags.PassDoubleDash)
	parser.NamespaceDelimiter = "-"
	twentythousandtonnesofcrudeoil.TheEnvironmentIsPerfectlySafe(parser, "FR_")

	_, err := parser.Parse()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	database, err := db.New(config.Postgres.ConnectionString())
	if err != nil {
		panic(err)
	}

	exp := exporter.Exporter{
		TelemetryPath: config.TelemetryPath,
		ListenAddress: config.ListenAddress,
		Collectors: []prometheus.Collector{
			&collectors.WorkersByState{Db: database},
			&collectors.PipelinesPerTeam{Db: database},
			&collectors.Resources{Db: database},
			&collectors.Teams{Db: database},
		},
	}

	go handleSignals(&exp, database)

	err = exp.Listen()
	if err != nil {
		panic(err)
	}
}
