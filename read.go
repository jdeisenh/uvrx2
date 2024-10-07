package uvrx2

import (
	"fmt"
	"strings"

	"github.com/brutella/can"
	"github.com/brutella/canopen"
	"github.com/brutella/canopen/sdo"
	"github.com/brutella/uvr"
)

func ReadFromIndex(idx canopen.ObjectIndex, nodeID uint8, bus *can.Bus) (interface{}, error) {
	upload := sdo.Upload{
		ObjectIndex:   idx,
		RequestCobID:  uint16(uvr.SSDOClientToServer2) + uint16(nodeID),
		ResponseCobID: uint16(uvr.SSDOServerToClient2) + uint16(nodeID),
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

// Parse list of 16bit Unicode characters
func printableASCIIString(b []byte) string {

	//fmt.Println(cpd.DecodeUTF16le(string(b)))
	var ascii strings.Builder
	var elem uint16
	for i, b := range b {
		if i%2 == 0 {
			elem = uint16(b)
			continue
		}
		elem = elem | uint16(b)<<8
		if elem == 0xfffc {
			// Skip strange control characters
			continue
		}
		ascii.WriteRune(rune(elem))
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
