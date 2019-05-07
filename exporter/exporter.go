package exporter

import (
	"net"
	"net/http"

	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Exporter struct {
	ListenAddress string
	TelemetryPath string

	listener net.Listener
}

// Listen initiates the HTTP server using the configurations
// provided via ExporterConfig.
//
// This is a blocking method - make sure you either make use of
// goroutines to not block if needed.
func (e *Exporter) Listen() (err error) {
	http.Handle(e.TelemetryPath, promhttp.Handler())

	e.listener, err = net.Listen("tcp", e.ListenAddress)
	if err != nil {
		err = errors.Wrapf(err,
			"failed to listen on address %s", e.ListenAddress)
		return
	}

	err = http.Serve(e.listener, nil)
	if err != nil {
		err = errors.Wrapf(err,
			"failed listening on address %s",
			e.ListenAddress)
		return
	}

	return
}

// Stop closes the tcp listener (if exists).
func (e *Exporter) Stop() (err error) {
	if e.listener == nil {
		return
	}

	err = e.listener.Close()
	return
}
