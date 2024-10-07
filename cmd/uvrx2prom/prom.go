package main

import (
	"flag"
	"net/http"

	"github.com/brutella/can"
	"github.com/jdeisenh/uvrx2"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
)

const prefix = "uvr2"

var listof = []struct {
	name string
	desc string
	Idx  uint16
	Sub  uint8
}{
	{"input_brenner_temp", "Temperatur Brenner Vorlauf", 8272, 0},
	{"input_aussen_temp", "Temperatur Aussen", 8272, 1},
}

// Implements prometheus.Collector
type CustomCollector struct {
	client *uvrx2.Client
	desc   []*prometheus.Desc
}

func (cm *CustomCollector) Collect(ch chan<- prometheus.Metric) {

	for k, x := range listof {
		value, e := uvrx2.NewElement(cm.client, x.Idx, x.Sub).Read()
		if e != nil {
			ch <- prometheus.MustNewConstMetric(
				cm.desc[k],
				prometheus.GaugeValue,
				value.Float64(),
			)
		}
	}
}

func (cm *CustomCollector) Describe(dc chan<- *prometheus.Desc) {
	prometheus.DescribeByCollect(cm, dc)
}

func NewCustomCollector(client *uvrx2.Client) *CustomCollector {

	var list []*prometheus.Desc
	for _, x := range listof {
		list = append(list, prometheus.NewDesc(prefix+"_"+x.name, x.desc, nil, nil))
	}
	res := &CustomCollector{
		client: client,
		desc:   list,
	}
	return res
}

func main() {

	var (
		clientId = flag.Int("client_id", 16, "id of the client; range from [1...254]")
		serverId = flag.Int("server_id", 32, "id of the server to which the client connects to: range from [1...254]")
		iface    = flag.String("if", "can0", "name of the can network interface")
	)

	flag.Parse()

	bus, err := can.NewBusForInterfaceWithName(*iface)

	if err != nil {
		log.Fatal(err)
	}
	go bus.ConnectAndPublish()

	nodeID := uint8(*clientId)
	uvrID := uint8(*serverId)

	c := uvrx2.NewClient(nodeID, bus)
	c.Connect(uvrID)

	prometheus.MustRegister(NewCustomCollector(c))
	http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe(":8080", nil))
}
