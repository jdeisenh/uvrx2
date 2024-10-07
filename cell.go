package uvrx2

import (
	"fmt"
	"strings"

	log "github.com/sirupsen/logrus"
)

type Cell []byte

func scaleInt(mode byte, val int32) float64 {
	switch mode {
	case 1, 7, 8:
		return float64(val) / 10
	default:
		return float64(val)
	}
}

func printNumberType(sort byte, val int32) string {
	switch sort {
	case 1, 7:
		return fmt.Sprintf("%g°", float32(val)/10)
	case 4:
		return fmt.Sprintf("%02d:%02d:%02d", val/3600, val/60%60, val%60)
	case 8:
		return fmt.Sprintf("%g%%", float32(val)/10)
	case 0x2b:
		switch val {
		case 0:
			return "Aus"
		case 1:
			return "Ein"
		default:
			return fmt.Sprintf("bool %d ??", val)
		}

	case 0x2c:
		switch val {
		case 0:
			return "Nein"
		case 1:
			return "Ja"
		default:
			return fmt.Sprintf("0x2c %d ??", val)
		}
	case 0x30:
		switch val {
		case 0:
			return "Nicht aktiv"
		case 1:
			return "Normal"
		case 2:
			return "Abgesenkt"
		case 3:
			return "Standby"
		case 4:
			return "Frostzschutz"
		case 6:
			return "Urlaub"
		case 7:
			return "Feiertag"
		case 8:
			return "Party"
		case 9:
			return "Störung"
		case 10:
			return "Wartung"
		case 11:
			return "Externe VL_Solltemp"
		default:
			return fmt.Sprintf("Unbekannt: %d", val)
		}
	case 0x3c:
		return fmt.Sprintf("%02d:%02d", val/60, val%60)
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
	// Decode by type
	switch b[0] {
	case 31:
		return printStringType(b[1], printableASCIIString(b[2:]))
	case 0xd0, 0x50, 0xc0, 0xb0:
		return printNumberType(b[1], parseInt(b[2:]))
	default:
		return fmt.Sprintf("%X %v", b[0], b[1:])
	}
}

// First byte represents the tpe
func (b Cell) TypeCode() byte {
	return b[0]
}

func (b Cell) Float64() float64 {
	switch b.TypeCode() {
	case 0xd0, 0x50, 0xc0, 0xb0:
		return scaleInt(b[1], parseInt(b[2:]))
	default:
		log.Warnf("Not convertable to double")
		return 0.0
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
