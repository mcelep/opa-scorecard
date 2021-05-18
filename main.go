package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Exporter struct {
	mirthEndpoint, mirthUsername, mirthPassword string
}

func NewExporter(mirthEndpoint string, mirthUsername string, mirthPassword string) *Exporter {
	return &Exporter{
		mirthEndpoint: mirthEndpoint,
		mirthUsername: mirthUsername,
		mirthPassword: mirthPassword,
	}
}

func (e *Exporter) Describe(ch chan<- *prometheus.Desc) {
	ch <- up
	//ch <- messagesReceived
}

func (e *Exporter) Collect(ch chan<- prometheus.Metric) {

	ch <- prometheus.MustNewConstMetric(
		up, prometheus.GaugeValue, 1,
	)

}

var (
	listenAddress = flag.String("web.listen-address", ":9141",
		"Address to listen on for telemetry")
	metricsPath = flag.String("web.telemetry-path", "/metrics",
		"Path under which to expose metrics")

	namespace = "opa_scorecard"
	// Metrics
	up = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "up"),
		"Was the last Mirth query successful.",
		nil, nil,
	)
	messagesReceived = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "messages_received_total"),
		"How many messages have been received (per channel).",
		[]string{"channel"}, nil,
	)
)

func main() {
	flag.Parse()

	exporter := NewExporter("test", "test", "test")
	prometheus.Unregister(prometheus.NewGoCollector())
	prometheus.MustRegister(exporter)

	http.Handle(*metricsPath, promhttp.Handler())
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`<html>
             <head><title>OPA ScoreCard Exporter</title></head>
             <body>
             <h1>OPA ScoreCard Exporter</h1>
             <p><a href='` + *metricsPath + `'>Metrics</a></p>
             </body>
             </html>`))
	})
	log.Fatal(http.ListenAndServe(*listenAddress, nil))
}
