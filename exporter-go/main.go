package main

import (
	"context"
	"encoding/json"
	"flag"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	restclient "k8s.io/client-go/rest"

	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	controllerClient "sigs.k8s.io/controller-runtime/pkg/client"
)

// StatusViolation represents each violation under status
type StatusViolation struct {
	Kind              string `json:"kind"`
	Name              string `json:"name"`
	Namespace         string `json:"namespace,omitempty"`
	Message           string `json:"message"`
	EnforcementAction string `json:"enforcementAction"`
}

type WrappedStatusViolation struct {
	*StatusViolation
	ConstraintKind string
	ConstraintName string
}

const (
	constraintsGV           = "constraints.gatekeeper.sh/v1beta1"
	constraintsGroup        = "constraints.gatekeeper.sh"
	constraintsGroupVersion = "v1beta1"
)

var (
	listenAddress = flag.String("web.listen-address", ":9141",
		"Address to listen on for telemetry")
	metricsPath = flag.String("web.telemetry-path", "/metrics",
		"Path under which to expose metrics")

	inCluster = flag.Bool("incluster", false,
		"Does the exporter run within a K8S cluster, when true it will try to look for K8S service account details in the usual location.")

	namespace = "opa_scorecard"
	// Metrics
	up = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "up"),
		"Was the last OPA violation query successful.",
		nil, nil,
	)
	violation = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "violations"),
		"OPA violations for all constraints",
		[]string{"kind", "name", "violating_kind", "violating_name", "violating_namespace", "violation_msg", "violation_enforcement"}, nil,
	)
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
	ch <- up
	ch <- violation
}

func (e *Exporter) Collect(ch chan<- prometheus.Metric) {

	ch <- prometheus.MustNewConstMetric(
		up, prometheus.GaugeValue, 1,
	)
	for _, m := range metrics {
		ch <- m
	}

}

func createConfig() (*restclient.Config, error) {
	if *inCluster {
		log.Println("Using incluster K8S client")
		return rest.InClusterConfig()
	} else {
		home, err := os.UserHomeDir()
		if err != nil {
			log.Println("Could not find user HomeDir" + err.Error())
			return nil, err
		}

		kubeconfig := filepath.Join(home, ".kube", "config")

		// use the current context in kubeconfig
		config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		return config, nil
	}
}

func createKubeClient() (*kubernetes.Clientset, error) {

	config, err := createConfig()
	if err != nil {
		log.Println(err)
		return nil, err
	}

	// create the clientset
	clientset, err := kubernetes.NewForConfig(config)

	if err != nil {
		log.Println(err)
		return nil, err
	}
	return clientset, nil
}

func createKubeClientGroupVersion() (controllerClient.Client, error) {
	config, err := createConfig()
	if err != nil {
		log.Println(err)
		return nil, err
	}

	client, err := controllerClient.New(config, controllerClient.Options{})
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return client, nil
}

func getConstraintViolations() ([]WrappedStatusViolation, error) {
	client, err := createKubeClient()
	if err != nil {
		return nil, err
	}

	constraints, err := client.ServerResourcesForGroupVersion(constraintsGV)
	if err != nil {
		return nil, err
	}

	cClient, err := createKubeClientGroupVersion()
	if err != nil {
		return nil, err
	}

	ret := []WrappedStatusViolation{}
	for _, r := range constraints.APIResources {
		canList := false
		for _, verb := range r.Verbs {
			if verb == "list" {
				canList = true
				break
			}
		}

		if !canList {
			continue
		}
		actual := &unstructured.UnstructuredList{}
		actual.SetGroupVersionKind(schema.GroupVersionKind{
			Group:   constraintsGroup,
			Kind:    r.Kind,
			Version: constraintsGroupVersion,
		})

		err = cClient.List(context.TODO(), actual)
		if err != nil {
			return nil, err
		}

		if len(actual.Items) > 0 {
			for _, constraint := range actual.Items {
				kind := constraint.GetKind()
				name := constraint.GetName()
				namespace := constraint.GetNamespace()
				log.Default().Printf("Kind:%s, Name:%s, Namespace:%s \n", kind, name, namespace)
				var obj map[string]interface{} = constraint.Object
				var status map[string]interface{}
				data, err := json.Marshal(obj["status"])
				if err != nil {
					log.Println(err)
					continue
				}
				json.Unmarshal(data, &status)
				if status["totalViolations"].(float64) > 0 {

					var violations []interface{}
					data, err := json.Marshal(status["violations"])
					if err != nil {
						return nil, err
					}
					json.Unmarshal(data, &violations)
					for _, violation := range violations {
						data, err := json.Marshal(violation)
						if err != nil {
							return nil, err
						}
						var viol StatusViolation
						json.Unmarshal(data, &viol)
						ret = append(ret, WrappedStatusViolation{ConstraintKind: constraint.GetKind(), ConstraintName: constraint.GetName(), StatusViolation: &viol})
					}
				}
			}
		}

	}
	return ret, nil
}

func (e *Exporter) startTimer() {
	ticker = time.NewTicker(10 * time.Second)
	go func() {
		for {
			select {
			case <-done:
				return
			case t := <-ticker.C:
				log.Println("Tick at", t)
				violations, err := getConstraintViolations()
				if err != nil {
					log.Printf("%+v\n", err)
				}

				tmpMetrics := make([]prometheus.Metric, 0)

				for _, v := range violations {
					//"kind", "name", "namespace", "msg"
					metric := prometheus.MustNewConstMetric(violation, prometheus.GaugeValue, 1, v.ConstraintKind, v.ConstraintName, v.Kind, v.Name, v.Namespace, v.Message, v.EnforcementAction)
					tmpMetrics = append(tmpMetrics, metric)
				}
				metrics = tmpMetrics
			}
		}
	}()
}

func main() {
	flag.Parse()

	exporter := NewExporter()
	exporter.startTimer()
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
