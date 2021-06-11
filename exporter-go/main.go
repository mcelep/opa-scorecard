package main

import (
	"flag"
	"log"
	"net/http"
	"time"

	"github.com/mcelep/opa_scorecard_exporter/pkg/opa"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	listenAddress = flag.String("web.listen-address", ":9141",
		"Address to listen on for telemetry")
	metricsPath = flag.String("web.telemetry-path", "/metrics",
		"Path under which to expose metrics")

	inCluster = flag.Bool("incluster", false,
		"Does the exporter run within a K8S cluster, when true it will try to look for K8S service account details in the usual location.")

	ticker  *time.Ticker
	done    = make(chan bool)
	metrics = []prometheus.Metric{}
)

type Exporter struct {
}

func NewExporter() *Exporter {
	return &Exporter{}
}

func (e *Exporter) Describe(ch chan<- *prometheus.Desc) {
	ch <- opa.Up
	ch <- opa.ConstraintViolation
	ch <- opa.ConstraintInformation
}

func (e *Exporter) Collect(ch chan<- prometheus.Metric) {

	ch <- prometheus.MustNewConstMetric(
		opa.Up, prometheus.GaugeValue, 1,
	)
	for _, m := range metrics {
		ch <- m
	}

}

func (e *Exporter) startScheduled(t time.Duration) {
	ticker = time.NewTicker(t)
	go func() {
		for {
			select {
			case <-done:
				return
			case t := <-ticker.C:
				log.Println("Tick at", t)
				constraints, err := opa.GetConstraints(inCluster)
				if err != nil {
					log.Printf("%+v\n", err)
				}
				allMetrics := make([]prometheus.Metric, 0)
				violationMetrics := opa.ExportViolations(constraints)
				allMetrics = append(allMetrics, violationMetrics...)

				constraintInformationMetrics := opa.ExportConstraintInformation(constraints)
				allMetrics = append(allMetrics, constraintInformationMetrics...)

				metrics = allMetrics
			}
		}
	}()
}

func main() {
	flag.Parse()

	exporter := NewExporter()
	exporter.startScheduled(10 * time.Second)
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
