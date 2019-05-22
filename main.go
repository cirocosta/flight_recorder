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
		Disable       []string `long:"disable" short:"d" description:"collectors to disable"`
		TelemetryPath string   `long:"path" default:"/" description:"path to serve metrics"`
		ListenAddress string   `long:"address" default:":9000" description:"address to listen for prometheus scraping"`

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

	collectorsMap := map[string]prometheus.Collector{
		"pipelines": &collectors.Collector{
			Description: prometheus.NewDesc(
				"flight_recorder_pipelines",
				"Number of pipelines set",
				[]string{"team"},
				nil,
			),
			RetrievalFunc: database.Pipelines,
		},
		"containers": &collectors.Collector{
			Description: prometheus.NewDesc(
				"flight_recorder_containers",
				"Number of containers",
				[]string{"type", "worker", "state"},
				nil,
			),
			RetrievalFunc: database.Containers,
		},
		"workers": &collectors.Collector{
			Description: prometheus.NewDesc(
				"flight_recorder_workers",
				"Per-state worker count",
				[]string{"state"},
				nil,
			),
			RetrievalFunc: database.Workers,
		},
		"resources": &collectors.Collector{
			Description: prometheus.NewDesc(
				"flight_recorder_resources",
				"Number of resources configured",
				[]string{"pipeline", "team"},
				nil,
			),
			RetrievalFunc: database.Resources,
		},
		"teams": &collectors.Collector{
			Description: prometheus.NewDesc(
				"flight_recorder_teams",
				"Number of teams",
				nil,
				nil,
			),
			RetrievalFunc: database.Teams,
		},
		"builds": &collectors.Collector{
			Description: prometheus.NewDesc(
				"flight_recorder_builds",
				"Number of builds",
				[]string{"status", "team", "pipeline"},
				nil,
			),
			RetrievalFunc: database.Builds,
		},
		"table-stat": &collectors.Collector{
			Description: prometheus.NewDesc(
				"flight_recorder_pg_table_stat",
				"PostgreSQL user table statistics",
				[]string{"statistic", "table"},
				nil,
			),
			RetrievalFunc: database.TableStat,
		},
		"tables-sizes": &collectors.Collector{
			Description: prometheus.NewDesc(
				"flight_recorder_pg_table_sizes_bytes",
				"PostgreSQL table sizes",
				[]string{"table"},
				nil,
			),
			RetrievalFunc: database.TableSizes,
		},
	}

	for _, collectorToDisable := range config.Disable {
		_, ok := collectorsMap[collectorToDisable]
		if !ok {
			panic(fmt.Errorf("unknown collector to be disabled: " + collectorToDisable))
		}

		fmt.Println("disabled:\t" + collectorToDisable)
		delete(collectorsMap, collectorToDisable)
	}

	var collectors = []prometheus.Collector{}
	for name, collector := range collectorsMap {
		fmt.Println("enabled:\t" + name)
		collectors = append(collectors, collector)
	}

	exp := exporter.Exporter{
		TelemetryPath: config.TelemetryPath,
		ListenAddress: config.ListenAddress,
		Collectors:    collectors,
	}

	go handleSignals(&exp, database)

	fmt.Println("listening on ", config.ListenAddress, config.TelemetryPath)

	err = exp.Listen()
	if err != nil {
		panic(err)
	}
}
