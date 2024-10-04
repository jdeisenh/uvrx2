package main

import (
	"fmt"
	"github.com/brutella/can"
	"github.com/brutella/canopen"
	"github.com/brutella/canopen/sdo"
	"strings"
)

func ReadFromIndex(idx canopen.ObjectIndex, nodeID uint8, bus *can.Bus) (interface{}, error) {
	upload := sdo.Upload{
		ObjectIndex:   idx,
		RequestCobID:  uint16(SSDOClientToServer2) + uint16(nodeID),
		ResponseCobID: uint16(SSDOServerToClient2) + uint16(nodeID),
	}

	b, err := upload.Do(bus)

	if err != nil {
		return nil, err
	}
	switch b[0] {
	case 31:
		return fmt.Sprintf("%d %s", b[1], printableASCIIString(b[2:])), nil
	case 0xd0, 0x50:
		return fmt.Sprintf("%d %d", b[1], parseInt32(b[2:])), nil
	case 0xc0:
		return fmt.Sprintf("%d %d", b[1], parseInt(b[2:])), nil
	case 0xb0:
		return fmt.Sprintf("%d %d", b[1], parseInt(b[2:])), nil
	default:
		return fmt.Sprintf("%X %v", b[0], b[1:]), nil
	}
}

func printableASCIIString(b []byte) string {
	var ascii strings.Builder
	for _, b := range b {
		if b >= 32 && b <= 126 || b > 160 {
			ascii.WriteRune(rune(b))
		}
	}

	return ascii.String()
}

func parseInt(b []byte) int32 {
	r := int32(0)
	for i := len(b) - 1; i >= 0; i-- {
		r = r * 256
		r += int32(b[i])
	}
	return r
}
func parseInt32(b []byte) int32 {
	return (int32(b[3]) << 24) + (int32(b[2]) << 16) + (int32(b[1]) << 8) + int32(b[0])
}
