package uvrx2

import (
	"time"

	"github.com/brutella/can"
	"github.com/brutella/canopen"
	"github.com/brutella/uvr"
)

type Client struct {
	id        uint8
	bus       *can.Bus
	heartbeat chan<- struct{}
}

func NewClient(id uint8, bus *can.Bus) *Client {
	return &Client{id, bus, nil}
}

func (c *Client) Connect(id uint8) error {
	c.heartbeat = canopen.ProduceHeartbeat(c.id, canopen.Operational, c.bus, time.Second*10)
	return uvr.Connect(id, c.id, c.bus)
}

func (c *Client) Disconnect(id uint8) error {
	c.heartbeat <- struct{}{}
	return uvr.Disconnect(id, c.id, c.bus)
}

func (c *Client) Write(b []byte, i canopen.ObjectIndex) error {
	return uvr.WriteToIndex(i, b, c.id, c.bus)
}
