package main

import (
	"flag"
	"log"
	"fmt"

	"github.com/brutella/can"

	"github.com/jdeisenh/uvrx2"
)

func main() {
	var (
		clientId = flag.Int("client_id", 16, "id of the client; range from [1...254]")
		serverId = flag.Int("server_id", 32, "id of the server to which the client connects to: range from [1...254]")
		iface    = flag.String("if", "can0", "name of the can network interface")
	)

	flag.Parse()

	log.SetFlags(log.LstdFlags | log.Lshortfile)

	bus, err := can.NewBusForInterfaceWithName(*iface)

	if err != nil {
		log.Fatal(err)
	}
	go bus.ConnectAndPublish()

	nodeID := uint8(*clientId)
	uvrID := uint8(*serverId)

	c := uvrx2.NewClient(nodeID, bus)
	c.Connect(uvrID)

	Show(c)

	c.Disconnect(uvrID)
	bus.Disconnect()
}

// Show dumps values for 'interesting" data
func Show(c *uvrx2.Client) {

        for _, device := range uvrx2.Interestingdata {
                fmt.Printf("Device: %s\n", device.Device)
                for _, m := range device.Descriptors {

                        name := m.Name
                        if m.Getname {
                                n, e := uvrx2.NewElement(c, m.Idx+4096, m.Sub).Read()
                                if e == nil {
                                        name = n.String()
                                }
                        }
                        got, e := uvrx2.NewElement(c, m.Idx, m.Sub).Read()
                        if e == nil {
                                fmt.Printf("\t%-30.30s: %s\n", name, got)
                        }
                }
        }
}



