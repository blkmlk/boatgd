package can

import (
	"github.com/brutella/can"
	"log"
	"../db"
	"time"
	"math"
	"encoding/json"
)

//Set your CAN interface
const CANInterface = "vcan"

const Timeout = 15 * time.Minute
const MaxMeters = 10

var lastWrite = time.Now()
var lastLatitude float64 = 0
var lastLongitude float64 = 0
var lastMeters float64 = 0

type GpsData struct {
	Id				int 	`json:"id"`
	Description		string	`json:"description"`
	Latitude		float64	`json:"latitude"`
	Longitude		float64	`json:"longitude"`
}

func RunCANReader() {
	bus, err := can.NewBusForInterfaceWithName(CANInterface)

	if err != nil {
		log.Fatal("CAN:", err.Error())
	}

	bus.SubscribeFunc(handleFrame)
}

func handleFrame(frame can.Frame) {
	pgn := Pgn(frame)

	switch pgn.Id() {
	case PgnGps: {
		lat, long := pgn.GetGpsLatLon()

		lastMeters += getMeters(lat, lastLatitude, long, lastLongitude)

		if lastMeters >= MaxMeters || time.Since(lastWrite) >= Timeout {
			lastMeters = 0
		} else {
			break
		}

		data := GpsData{
			Id:          pgn.Id(),
			Description: "GPS Data",
			Latitude:    lat,
			Longitude:   long,
		}

		blob, err := json.Marshal(&data)

		if err != nil {
			break
		}

		lastWrite = time.Now()

		now := lastWrite.Unix()

		db.GetDB().WriteData(pgn.Id(), blob, now)
	}
	}
}

func getMeters(lat1, lat2, long1, long2 float64) float64 {
	defer func() {
		lastLatitude = lat1
		lastLongitude = long1
	}()

	if lastLatitude == 0 || lastLongitude == 0 {
		return 0
	}

	latMid := float64((lat1 + lat2) / 2.0)

	mPerDegLat :=  111132.954 - 559.822 * math.Cos(2.0 * latMid) + 1.175 * math.Cos(4.0 + latMid)
	mPerDegLong := (math.Pi/180) * 6367449 * math.Cos(latMid)

	deltaLat := math.Abs(lat1 - lat2)
	deltaLong := math.Abs(long1 - long2)

	return math.Sqrt(math.Pow(deltaLat * mPerDegLat, 2) + math.Pow(deltaLong * mPerDegLong, 2))
}