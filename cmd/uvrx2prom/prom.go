package main

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
)

// Implements prometheus.Collector
type CustomCollector struct {
	m1 *prometheus.Desc
}

func (cm *CustomCollector) Collect(ch chan<- prometheus.Metric) {

	ch <- prometheus.MustNewConstMetric(cm.m1, prometheus.GaugeValue, 1)
}

func (cm *CustomCollector) Describe(dc chan<- *prometheus.Desc) {
	prometheus.DescribeByCollect(cm, dc)
}

func NewCustomCollector() *CustomCollector {
	return &CustomCollector{
		m1: prometheus.NewDesc("m1_metric", "A desc2", nil, nil),
	}
}

func main() {
	prometheus.MustRegister(NewCustomCollector())
	http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe(":8080", nil))
}
