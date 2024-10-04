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

type Cell []byte

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

func printNumberType(sort byte, val int32) string {
	switch sort {
	case 1:
		return fmt.Sprintf("%g°", float32(val)/10)
	case 4:
		return fmt.Sprintf("%02d:%02d:%02d", val/3600,val/60%60,val%60)
	case 8:
		return fmt.Sprintf("%g%%", float32(val)/10)
	case 0x2b:
		return fmt.Sprintf("bool %d", val)
	case 0x3c:
		return fmt.Sprintf("%02d:%02d", val/60,val%60)
	default:
		return fmt.Sprintf("I %x %d", sort, val)
	}
}

func printStringType(sort byte, rest string) string {
	switch sort {
	case 0:
		return rest
	default:
		return fmt.Sprintf("%s (#%d)", rest, sort)
	}
}

func (b Cell) String() string {
	switch b[0] {
	case 31:
		return printStringType(b[1], printableASCIIString(b[2:]))
	case 0xd0, 0x50, 0xc0, 0xb0:
		return printNumberType(b[1], parseInt(b[2:]))
	default:
		return fmt.Sprintf("%X %v", b[0], b[1:])
	}
}

type descriptor struct {
	name    string
	idx     uint16
	sub     uint8
	getname bool
}

var dumpdata = []struct {
	device      string
	descriptors []descriptor
}{
	{
		"Eingänge",
		[]descriptor{
			// Eingänge
			{"T. Kessel VL", 8272, 0, false},
			{"T. Aussen", 8272, 1, false},
			{"T. HK VL 1", 8272, 2, false},
			{"T. HK VL 2", 8272, 3, false},
			{"T. HK VL 3", 8272, 4, false},
		}}, {
		"HKR #1 Büro",
		[]descriptor{
			// HKR #1
			{"Betrieb", 10299, 2, true},
			{"Vorlaufsoll 1", 11521, 2, true},
			{"Mischer 1", 11529, 2, true},
			{"Betriebsart 1", 11535, 2, true},
			{"Absenk temp", 10256, 2, true},
			{"Normal temp", 10257, 2, true},
		}}, {
		"HKR #2 Fussboden",
		[]descriptor{

			// HKR #2 Fussboden
			{"Modus 3", 10299, 3, true},
			{"Vorlaufsoll 2", 11521, 3, true},
			{"Mischer 1", 11529, 3, true},
			{"Absenk temp", 10256, 3, true},
			{"Normal temp", 10257, 3, true},
			{"Betriebsart 1", 11535, 3, true},
		}}, {
		"HKD #3 Werkstatt",
		[]descriptor{
			{"Modus 3", 10299, 4, true},
			{"Vorlaufsoll 3", 11521, 4, true},
			{"Betriebsart 1", 11535, 4, true},
			{"Offset", 10269, 4, true},
		}}, {
		"Anforderung Heizung",
		[]descriptor{
			{"Anforderung Heizung", 11019, 5, true},
			{"Anforderung Soll", 11031, 5, true},
			{"Anforderung", 11521, 5, true},
		}}, {
		"Mischer 1",
		[]descriptor{
			{"Mischer 1 Laufzeit", 8348, 0, true},
		}}, {
		"Pumpe #1",
		[]descriptor{
			{"Brenner Zustand", 8400, 4, true},
			{"Brenner Laufzeit", 8402, 4, true},
		}}, {
		"Pumpe #2",
		[]descriptor{
			{"Brenner Zustand", 8400, 3, true},
			{"Brenner Laufzeit", 8402, 3, true},
		}}, {
		"Pumpe #3",
		[]descriptor{
			{"Brenner Zustand", 8400, 2, true},
			{"Brenner Laufzeit", 8402, 2, true},
		}}, {
		"Brenner",
		[]descriptor{
			{"Brenner Zustand", 8400, 5, true},
			{"Brenner Laufzeit", 8402, 5, true},
		}}, {
		"Zeitprogramm",
		[]descriptor{
			{"Ein", 10312, 1, true},
			{"Aus", 10347, 1, true},
		}},
}

func Show(c *Client) {

	for _, device := range dumpdata {
		log.Printf("Device: %s", device.device)
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
				log.Printf("%-30.30s: %s", name, got)
			}
		}
	}
}
