package uvrx2

import (
	"github.com/brutella/can"
	"github.com/brutella/canopen"
	"github.com/brutella/canopen/sdo"
	"github.com/brutella/uvr"
)

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
		RequestCobID:  uint16(uvr.SSDOClientToServer2) + uint16(nodeID),
		ResponseCobID: uint16(uvr.SSDOServerToClient2) + uint16(nodeID),
	}

	b, err := upload.Do(bus)

	if err != nil {
		return nil, err
	}
	return b, nil
}
