package exporter

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/tris/weatherflow"
)

const (
	timeoutDuration = 30 * time.Minute
)

type WeatherExporter struct {
	collectors map[string]map[int]*WeatherCollector
	clients    map[string]*weatherflow.Client
	registries map[string]map[int]*prometheus.Registry
	mu         sync.RWMutex
}

func NewWeatherExporter() *WeatherExporter {
	return &WeatherExporter{
		registries: make(map[string]map[int]*prometheus.Registry),
		collectors: make(map[string]map[int]*WeatherCollector),
		clients:    make(map[string]*weatherflow.Client),
	}
}

func prefixedLogger(apiToken string, logFunc weatherflow.Logf) weatherflow.Logf {
	prefix := apiToken
	if len(apiToken) > 8 {
		prefix = apiToken[:8]
	}

	return func(format string, args ...interface{}) {
		formatWithPrefix := fmt.Sprintf("[%s...] %s", prefix, format)
		logFunc(formatWithPrefix, args...)
	}
}

func (we *WeatherExporter) ScrapeHandler(w http.ResponseWriter, r *http.Request) {
	apiToken := r.URL.Query().Get("token")
	deviceID, err := strconv.Atoi(r.URL.Query().Get("device_id"))
	if apiToken == "" || err != nil {
		http.Error(w, "Error: Missing or malformed query parameters 'token' and/or 'device_id'", http.StatusBadRequest)
		return
	}

	we.mu.Lock()
	collectorsForToken, ok := we.collectors[apiToken]
	if !ok {
		collectorsForToken = make(map[int]*WeatherCollector)
		we.collectors[apiToken] = collectorsForToken
	}

	collector, ok := collectorsForToken[deviceID]
	if ok {
		collector.timer.Reset(timeoutDuration)
	} else {
		collector = NewWeatherCollector(deviceID)
		collectorsForToken[deviceID] = collector
		we.initCollectorTimer(collector, apiToken, deviceID)

		client, clientExists := we.clients[apiToken]
		if !clientExists {
			client = weatherflow.NewClient(apiToken, prefixedLogger(apiToken, log.Printf))
			we.clients[apiToken] = client
			client.Start(func(msg weatherflow.Message) {
				collector.update(msg, apiToken)
			})
		}
		client.AddDevice(deviceID)
	}

	registriesForToken, ok := we.registries[apiToken]
	if !ok {
		registriesForToken = make(map[int]*prometheus.Registry)
		we.registries[apiToken] = registriesForToken
	}

	registry, ok := registriesForToken[deviceID]
	if !ok {
		registry = prometheus.NewRegistry()
		registry.MustRegister(collector)
		registriesForToken[deviceID] = registry
	}
	we.mu.Unlock()

	promhttp.HandlerFor(registry, promhttp.HandlerOpts{}).ServeHTTP(w, r)
}

func (we *WeatherExporter) initCollectorTimer(collector *WeatherCollector, apiToken string, deviceID int) {
	collector.timer = time.AfterFunc(timeoutDuration, func() {
		we.mu.Lock()
		defer we.mu.Unlock()

		delete(we.collectors[apiToken], deviceID)
		if len(we.collectors[apiToken]) == 0 {
			delete(we.collectors, apiToken)
		}

		client := we.clients[apiToken]
		client.RemoveDevice(deviceID)
		if client.DeviceCount() == 0 {
			client.Stop()
			delete(we.clients, apiToken)
		}

		we.registries[apiToken][deviceID].Unregister(collector)
		delete(we.registries[apiToken], deviceID)
		if len(we.registries[apiToken]) == 0 {
			delete(we.registries, apiToken)
		}
	})
}
