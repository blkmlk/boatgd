package can

import (
	"github.com/brutella/can"
	"encoding/binary"
)

const (
	PgnGps	= 129025
)

const (
	gpsResolution = 0.0000001
)

type Pgn can.Frame

func (p *Pgn) Id() (pgn int) {
	MS := (p.ID >> 24) & 0x03
	PF := (p.ID >> 16) & 0xFF
	PS := (p.ID >> 8) & 0xFF

	if PF > 239 {
		pgn = int((MS << 16) | (PF << 8) | PS)
	} else {
		pgn = int((MS << 16) | (PF << 8))
	}

	return
}

func (p *Pgn) GetGpsLatLon() (lat float64, lon float64) {
	lat = float64(int32(binary.LittleEndian.Uint32(p.Data[0:4]))) * gpsResolution
	lon = float64(int32(binary.LittleEndian.Uint32(p.Data[4:]))) * gpsResolution

	return
}