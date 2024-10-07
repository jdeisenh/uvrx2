package main

import (
"fmt"
)

type Cell []byte

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
