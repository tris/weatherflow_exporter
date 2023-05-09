package main

import (
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/tris/weatherflow_exporter/exporter"
)

const (
	// reuse same port from https://github.com/nalbury/tempest-exporter, for now
	defaultPort = 6969
)

func main() {
	http.HandleFunc("/metrics", promhttp.HandlerFor(prometheus.DefaultGatherer, promhttp.HandlerOpts{}).ServeHTTP)

	http.HandleFunc("/scrape", exporter.NewWeatherExporter().ScrapeHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = strconv.Itoa(defaultPort)
	}

	log.Printf("Listening on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
