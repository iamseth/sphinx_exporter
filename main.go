package main

import (
	"flag"
	"net/http"
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/log"
	"github.com/yunge/sphinx"
)

var (
	// Version will be set at build time.
	Version       = "0.0.0.dev"
	listenAddress = flag.String("web.listen-address", "0.0.0.0:9161", "Address to listen on for web interface and telemetry.")
	metricsPath   = flag.String("web.telemetry-path", "/metrics", "Path under which to expose metrics.")
	sphinxHost    = flag.String("sphinx.host", "localhost", "Hostname or IP address for Sphinx.")
	sphinxPort    = flag.Int("sphinx.port", 9312, "TCP port for Sphinx.")
)

const namespace = "sphinx"

// Exporter collects Sphinx metrics. It implements prometheus.Collector.
type Exporter struct {
	host string
	port int
	up   prometheus.Gauge
}

// NewExporter returns a new Sphinx exporter.
func NewExporter(host string, port int) *Exporter {
	return &Exporter{
		host: host,
		port: port,
		up:   newGauge("up", "Was the last scrape of Sphinx successful?"),
	}
}

// Describe describes all the metrics exported by the exporter.
func (e *Exporter) Describe(ch chan<- *prometheus.Desc) {
	metricCh := make(chan prometheus.Metric)
	doneCh := make(chan struct{})
	go func() {
		for m := range metricCh {
			ch <- m.Desc()
		}
		close(doneCh)
	}()
	e.Collect(metricCh)
	close(metricCh)
	<-doneCh
}

// Collect implements prometheus.Collector.
func (e *Exporter) Collect(ch chan<- prometheus.Metric) {

	opts := &sphinx.Options{
		Host:    e.host,
		Port:    e.port,
		Timeout: 5000,
	}
	sc := sphinx.NewClient(opts)
	if err := sc.Error(); err != nil {
		log.Error(err)
		e.up.Set(0)
		ch <- e.up
		return
	}
	resp, err := sc.Status()
	if err != nil {
		e.up.Set(0)
		log.Error(err)
		ch <- e.up
		return
	}
	e.up.Set(1)
	ch <- e.up
	for _, row := range resp {
		k, v := row[0], row[1]
		// If value is set to OFF, just skip it.
		if v == "OFF" {
			continue
		}
		// Convert all values to float64.
		f, err := strconv.ParseFloat(v, 64)
		if err != nil {
			log.Errorf("Unable to parse value %s for key %s", v, k)
			continue
		}
		switch k {
		case "uptime":
			ch <- counter("uptime", "testing", f)
		case "connections":
			ch <- counter("connections", "testing", f)
		case "maxed_out":
			ch <- counter("maxed_out", "testing", f)
		case "command_search":
			ch <- counter("command_search", "testing", f)
		case "command_excerpt":
			ch <- counter("command_excerpt", "testing", f)
		case "command_update":
			ch <- counter("command_update", "testing", f)
		case "command_delete":
			ch <- counter("command_delete", "testing", f)
		case "command_keywords":
			ch <- counter("command_keywords", "testing", f)
		case "command_persist":
			ch <- counter("command_persist", "testing", f)
		case "command_status":
			ch <- counter("command_status", "testing", f)
		case "command_flushattrs":
			ch <- counter("command_flushattrs", "testing", f)
		case "agent_connect":
			ch <- counter("agent_connect", "testing", f)
		case "agent_retry":
			ch <- counter("agent_retry", "testing", f)
		case "queries":
			ch <- counter("queries", "testing", f)
		case "dist_queries":
			ch <- counter("dist_queries", "testing", f)
		case "query_wall":
			ch <- counter("query_wall", "testing", f)
		case "query_cpu":
			ch <- counter("query_cpu", "testing", f)
		case "dist_wall":
			ch <- counter("dist_wall", "testing", f)
		case "dist_local":
			ch <- counter("dist_local", "testing", f)
		case "dist_wait":
			ch <- counter("dist_wait", "testing", f)
		case "query_reads":
			ch <- counter("query_reads", "testing", f)
		case "query_readkb":
			ch <- counter("query_readkb", "testing", f)
		case "query_readtime":
			ch <- counter("query_readtime", "testing", f)
		}
	}
}

func counter(name, description string, value float64) prometheus.Counter {
	c := prometheus.NewCounter(
		prometheus.CounterOpts{
			Namespace: namespace,
			Name:      name,
			Help:      description,
		},
	)
	c.Set(value)
	return c
}

func newGauge(name, description string) prometheus.Gauge {
	return prometheus.NewGauge(
		prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      name,
			Help:      description,
		},
	)
}

func main() {
	flag.Parse()
	exporter := NewExporter(*sphinxHost, *sphinxPort)
	prometheus.MustRegister(exporter)
	http.Handle(*metricsPath, prometheus.Handler())
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`<html>
                <head><title>Sphinx Exporter</title></head>
                <body>
                   <h1>Sphinx Exporter</h1>
                   <p><a href='` + *metricsPath + `'>Metrics</a></p>
                   </body>
                </html>
              `))
	})
	log.Infof("Starting sphinx_exporter %s.", Version)
	log.Fatal(http.ListenAndServe(*listenAddress, nil))
}
