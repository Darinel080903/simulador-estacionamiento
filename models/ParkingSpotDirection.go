package models

type ParkingSpotDirection struct {
	Direction string
	Point     float64
}

func newParkingSpotDirection(direction string, point float64) *ParkingSpotDirection {
	return &ParkingSpotDirection{
		Direction: direction,
		Point:     point,
	}
}
