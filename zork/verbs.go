package zork

import "encoding/binary"

func VersionSub() bool {
	Print("ZORK I: The Great Underground Empire")
	NewLine()
	Print("Copyright (c) 1981, 1982, 1983 Infocom, Inc. ")
	NewLine()
	if DollarZero[1]&8 > 0 {
		Print("Licensed to Tandy Corporation.")
		NewLine()
	}
	Print("ZORK is a registered trademark of Infocom, Inc.")
	NewLine()
	Print("Revision ")
	num := binary.LittleEndian.Uint32(DollarZero[4:8])
	PrintNumber(num & 2047)
	Print(" / Serial number ")
	for offset := 18; offset <= 23; offset++ {
		PrintChar(DollarZero[offset])
	}
	NewLine()
	return true
}

func LookSub() bool {
	if !DescribeRoom(true) {
		return false
	}
	return DescribeObjects(true)
}
