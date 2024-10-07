package main

import (
	"flag"
	"fmt"
	"github.com/brutella/can"
	"github.com/brutella/canopen"
	"github.com/brutella/canopen/sdo"
	"log"
)

func HandleCANopen(frame can.Frame) {
	log.Printf("%X % X\n", frame.ID, frame.Data)
}

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
	// bus.SubscribeFunc(HandleCANopen)
	go bus.ConnectAndPublish()

	nodeID := uint8(*clientId)
	uvrID := uint8(*serverId)

	c := NewClient(nodeID, bus)
	c.Connect(uvrID)

	Show(c)

	c.Disconnect(uvrID)
	bus.Disconnect()
}

type Element struct {
	c    *Client
	caoi canopen.ObjectIndex
}

func NewElement(c *Client, idx uint16, sub uint8) *Element {
	el := &Element{
		c:    c,
		caoi: canopen.NewObjectIndex(idx, sub),
	}
	return el
}

func (el *Element) Read() (Cell, error) {
	buf, err := Readbuf(el.caoi, el.c.id, el.c.bus)
	return Cell(buf), err
}

func Readbuf(idx canopen.ObjectIndex, nodeID uint8, bus *can.Bus) ([]byte, error) {
	upload := sdo.Upload{
		ObjectIndex:   idx,
		RequestCobID:  uint16(SSDOClientToServer2) + uint16(nodeID),
		ResponseCobID: uint16(SSDOServerToClient2) + uint16(nodeID),
	}

	b, err := upload.Do(bus)

	if err != nil {
		return nil, err
	}
	return b, nil
}

// Show dumps values for 'interesting" data
func Show(c *Client) {

	for _, device := range interestingdata {
		fmt.Printf("Device: %s\n", device.device)
		for _, m := range device.descriptors {

			name := m.name
			if m.getname {
				n, e := NewElement(c, m.idx+4096, m.sub).Read()
				if e == nil {
					name = n.String()
				}
			}
			got, e := NewElement(c, m.idx, m.sub).Read()
			if e == nil {
				fmt.Printf("\t%-30.30s: %s\n", name, got)
			}
		}
	}
}
