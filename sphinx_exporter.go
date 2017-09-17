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
	listenAddress = flag.String("web.listen-address", "0.0.0.0:9247", "Address to listen on for web interface and telemetry.")
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
		up:   newGauge("up", "1 if we're able to scrape metrics, otherwise 0."),
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
			ch <- counter("uptime", "Time in seconds searchd has been running.", f)
		case "connections":
			ch <- counter("connections", "Total number of connections made since startup.", f)
		case "maxed_out":
			ch <- counter("maxed_out", "Number of times connections were maxed out.", f)
		case "command_search":
			ch <- counter("command_search", "Sum of search commands since startup.", f)
		case "command_excerpt":
			ch <- counter("command_excerpt", "Sum of excerpt commands since startup.", f)
		case "command_update":
			ch <- counter("command_update", "Sum of update commands since startup.", f)
		case "command_delete":
			ch <- counter("command_delete", "Sum of delete commands since startup.", f)
		case "command_keywords":
			ch <- counter("command_keywords", "Sum of keywords commands since startup.", f)
		case "command_persist":
			ch <- counter("command_persist", "Sum of persist commands since startup.", f)
		case "command_status":
			ch <- counter("command_status", "Sum of status commands since startup.", f)
		case "command_flushattrs":
			ch <- counter("command_flushattrs", "Sum of flushattrs commands since startup.", f)
		case "agent_connect":
			ch <- counter("agent_connect", "TODO", f)
		case "agent_retry":
			ch <- counter("agent_retry", "TODO", f)
		case "queries":
			ch <- counter("queries", "Total number of queries run against Sphinx.", f)
		case "dist_queries":
			ch <- counter("dist_queries", "TODO", f)
		case "query_wall":
			ch <- counter("query_wall", "Total time running queries.", f)
		case "query_cpu":
			ch <- counter("query_cpu", "TODO", f)
		case "dist_wall":
			ch <- counter("dist_wall", "Total time running distributed queries.", f)
		case "dist_local":
			ch <- counter("dist_local", "TODO", f)
		case "dist_wait":
			ch <- counter("dist_wait", "TODO", f)
		case "query_reads":
			ch <- counter("query_reads", "TODO", f)
		case "query_readkb":
			ch <- counter("query_readkb", "TODO", f)
		case "query_readtime":
			ch <- counter("query_readtime", "TODO", f)
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
