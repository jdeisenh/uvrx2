package main


type descriptor struct {
	name    string
	idx     uint16
	sub     uint8
	getname bool
}

var interestingdata = []struct {
	device      string
	descriptors []descriptor
}{
	{
		"HKR #1 Büro",
		[]descriptor{
			{"Betrieb", 10299, 2, true},
			{"Betriebsart 1", 11535, 2, true},
			{"Vorlaufsoll 1", 11521, 2, true},
			{"Vorlauf", 11043, 2, true},
			{"Mischer 1", 11529, 2, true},
			{"Absenk temp", 10256, 2, true},
			{"Normal temp", 10257, 2, true},
			{"Offset", 10269, 2, true},
		},
	}, {
		"HKR #2 Fussboden",
		[]descriptor{

			{"Modus 3", 10299, 3, true},
			{"Betriebsart 2", 11535, 3, true},
			{"Vorlaufsoll 2", 11521, 3, true},
			{"Vorlauf", 11043, 3, true},
			{"Mischer 1", 11529, 3, true},
			{"Absenk temp", 10256, 3, true},
			{"Normal temp", 10257, 3, true},
			{"Offset", 10269, 3, true},
		},
	}, {
		"HKR #3 Werkstatt",
		[]descriptor{
			{"Modus 3", 10299, 4, true},
			{"Betriebsart 3", 11535, 4, true},
			{"Vorlaufsoll 3", 11521, 4, true},
			{"Vorlauf", 11043, 4, true},
			{"Absenk temp", 10256, 3, true},
			{"Normal temp", 10257, 3, true},
			{"Offset", 10269, 4, true},
		},
	}, {
		"Anforderung Heizung",
		[]descriptor{
			{"Anforderung Heizung", 11019, 5, true},
			{"Anforderung Soll", 11031, 5, true},
			{"Anforderung", 11521, 5, true},
		},
	}, {
		"Eingänge",
		[]descriptor{
			// Eingänge
			{"T. Kessel VL", 8272, 0, false},
			{"T. Aussen", 8272, 1, false},
			{"T. HK VL 1", 8272, 2, false},
			{"T. HK VL 2", 8272, 3, false},
			{"T. HK VL 3", 8272, 4, false},
		},
	}, {
		"Mischer 1",
		[]descriptor{
			{"Mischer 1 Laufzeit", 8348, 0, true},
		},
	}, {
		"Pumpe #1",
		[]descriptor{
			{"Brenner Zustand", 8400, 4, true},
			{"Brenner Laufzeit", 8402, 4, true},
		},
	}, {
		"Pumpe #2",
		[]descriptor{
			{"Brenner Zustand", 8400, 3, true},
			{"Brenner Laufzeit", 8402, 3, true},
		},
	}, {
		"Pumpe #3",
		[]descriptor{
			{"Brenner Zustand", 8400, 2, true},
			{"Brenner Laufzeit", 8402, 2, true},
		},
	}, {
		"Brenner",
		[]descriptor{
			{"Brenner Zustand", 8400, 5, true},
			{"Brenner Laufzeit", 8402, 5, true},
		},
	}, {
		"Zeitprogramm",
		[]descriptor{
			{"Ein", 10312, 1, true},
			{"Aus", 10347, 1, true},
		},
	}, {
		"System",
		[]descriptor{
			{"Reglerstart", 9360, 0, false},
			{"Sensorfehler Eingang", 9361, 0, false},
			{"Meldung", 9369, 0, false},
			{"Warnung", 9370, 0, false},
			{"Störung", 9371, 0, false},
			{"Fehler", 9374, 0, false},
		},
	},
}

