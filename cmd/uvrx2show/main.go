package main

import (
	"flag"
	"log"

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

	uvrx2.Show(c)

	c.Disconnect(uvrID)
	bus.Disconnect()
}
