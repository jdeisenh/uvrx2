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

const prefix = "uvrx2"

type al struct {
	name string
	Idx  uint16
	Sub  uint8
}

var listof = []struct {
	name  string
	desc  string
	param string
	atrl  []al
}{
	{
		"input_temp",
		"Temperatur Brenner Vorlauf",
		"sensor",
		[]al{
			{"Brenner VL", 8272, 0},
			{"Aussen", 8272, 1},
			{"HK 1", 8272, 2},
			{"HK 2", 8272, 3},
			{"HK 3", 8272, 4},
		},
	},
	{
		"mischer_pct",
		"Mischer Öffnung",
		"heizkreis",
		[]al{
			{"Büro", 11529,2},
			{"Werkstatt", 11529, 4},
			{"Fussbodenhzg", 11529, 3},
		},
	},
	{
		"vorlaufsoll_temp",
		"Vorlauf Solltemperatur",
		"heizkreis",
		[]al{
			{"Büro", 11521,2},
			{"Werkstatt", 11521, 4},
			{"Fussbodenhzg", 11521, 3},
		},
	},
	{
		"brenner_temp",
		"Brenner Temperatur",
		"typ",
		[]al{
			{"Ist", 11019,5},
			{"Soll", 11031, 5},
			{"Aktiv", 11521, 5},
		},
	},
	{
		"ausgang_aktiv",
		"Ausgang",
		"art",
		[]al{
			{"Brenner", 8400, 5},
			{"Pumpe 1", 8400, 4},
			{"Pumpe 2", 8400, 3},
			{"Pumpe 3", 8400, 2},
		},
	},
	{
		"ausgang_laufzeit",
		"Ausgang Einschaltzeit",
		"art",
		[]al{
			{"Brenner", 8402, 5},
			{"Pumpe 1", 8402, 4},
			{"Pumpe 2", 8402, 3},
			{"Pumpe 3", 8402, 2},
		},
	},
}

// Implements prometheus.Collector
type CustomCollector struct {
	client *uvrx2.Client
	desc   []*prometheus.Desc
}

func (cm *CustomCollector) Collect(ch chan<- prometheus.Metric) {

	for k, x := range listof {
		for _, z := range x.atrl {
			value, e := uvrx2.NewElement(cm.client, z.Idx, z.Sub).Read()
			if e == nil {
				ch <- prometheus.MustNewConstMetric(
					cm.desc[k],
					prometheus.GaugeValue,
					value.Float64(),
					z.name,
				)
			}
		}
	}
}

func (cm *CustomCollector) Describe(dc chan<- *prometheus.Desc) {
	prometheus.DescribeByCollect(cm, dc)
}

func NewCustomCollector(client *uvrx2.Client) *CustomCollector {

	var list []*prometheus.Desc
	for _, x := range listof {
		list = append(list, prometheus.NewDesc(prefix+"_"+x.name, x.desc, []string{x.param}, nil))
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
