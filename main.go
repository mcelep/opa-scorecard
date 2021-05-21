package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	controllerClient "sigs.k8s.io/controller-runtime/pkg/client"
)

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

func createKubeClient() (*kubernetes.Clientset, error) {
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

	// create the clientset
	clientset, err := kubernetes.NewForConfig(config)

	if err != nil {
		log.Println(err)
		return nil, err
	}
	return clientset, nil
}

func createKubeClientGroupVersion() (controllerClient.Client, error) {
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

	//config.GroupVersion = &schema.GroupVersion{Group: constraintsGroup, Version: constraintsGroupVersion}
	//config.NegotiatedSerializer = runtime.NewSimpleNegotiatedSerializer(runtime.SerializerInfo{EncodesAsText: true})
	//client, err := rest.RESTClientFor(config)
	client, err := controllerClient.New(config, controllerClient.Options{})
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return client, nil
}

func getConstraintViolations() {
	client, err := createKubeClient()
	if err != nil {
		log.Println(err)
		return
	}

	constraints, err := client.ServerResourcesForGroupVersion(constraintsGV)
	if err != nil {
		log.Println(err)
		return
	}

	cClient, err := createKubeClientGroupVersion()
	if err != nil {
		log.Println(err)
		return
	}

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
		log.Println(fmt.Sprintf("%s/%s", constraintsGV, r.Name))
		actual := &unstructured.UnstructuredList{}
		actual.SetGroupVersionKind(schema.GroupVersionKind{
			Group:   constraintsGroup,
			Kind:    r.Kind,
			Version: constraintsGroupVersion,
		})
		//key := controllerClient.ObjectKey{Namespace: "", Name: r.Name}
		err = cClient.List(context.TODO(), actual)
		//getResult := client.DiscoveryClient.RESTClient().Get().Resource(r.Name).Do(context.Background())
		if err != nil {
			log.Println(err)
			return
		}
		log.Println(actual)
		// log.Print(getResult)
	}
}

func main() {
	flag.Parse()

	getConstraintViolations()

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
