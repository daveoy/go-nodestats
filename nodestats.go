package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

func main() {
	var listenAddress = flag.String("listen-address", ":9666", "The address to listen on for HTTP requests.")
	// set up in cluster configs for kubernetes client
	config, err := rest.InClusterConfig()
	if err != nil {
		panic(err.Error())
	}
	// create the client
	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}
	// set up metrics we're tracking
	var labels = []string{"namespace", "node"}
	var podsPerNodeNamespace = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "pods_per_node",
		Help: "Number of Pods hosted on each Node",
	}, labels)
	prometheus.MustRegister(podsPerNodeNamespace)
	// here we do the work async
	go func() {
		fmt.Printf("entering main Goroutine")
		for {
			fmt.Printf("clearing gauge")
			podsPerNodeNamespace.Reset()
			// get pods in all the namespaces by omitting namespace
			// Or specify namespace to get pods in particular namespace
			pods, err := client.CoreV1().Pods("").List(context.TODO(), metav1.ListOptions{})
			if err != nil {
				panic(err.Error())
			}
			for _, item := range pods.Items {
				fmt.Printf("%+v %+v %+v\n", item.Name, item.Namespace, item.Spec.NodeName)
				podsPerNodeNamespace.WithLabelValues(item.Namespace, item.Spec.NodeName).Add(1)
			}
			time.Sleep(60 * time.Second)
		}
	}()
	// Expose the registered metrics via HTTP.
	http.Handle("/metrics", promhttp.HandlerFor(
		prometheus.DefaultGatherer,
		promhttp.HandlerOpts{
			// Opt into OpenMetrics to support exemplars.
			EnableOpenMetrics: false,
		},
	))
	log.Fatal(http.ListenAndServe(*listenAddress, nil))
}
