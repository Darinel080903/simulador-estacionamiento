package models

import (
	"github.com/oakmound/oak/v4/alg/floatgeom"
)

type ParkingSpot struct {
	area                 *floatgeom.Rect2
	directionsForParking *[]ParkingSpotDirection
	directionsForLeaving *[]ParkingSpotDirection
	number               int
	isAvailable          bool
}

func NewParkingSpot(x, y, x2, y2 float64, column, number int) *ParkingSpot {
	directionsForParking := getDirectionForParking(x, y, column)
	directionsForLeaving := getDirectionsForLeaving()
	area := floatgeom.NewRect2(x, y, x2, y2)

	return &ParkingSpot{
		area:                 &area,
		directionsForParking: directionsForParking,
		directionsForLeaving: directionsForLeaving,
		number:               number,
		isAvailable:          true,
	}
}

func getDirectionForParking(x, y float64, column int) *[]ParkingSpotDirection {
	var directions []ParkingSpotDirection

	if column == 1 {
		directions = append(directions, *newParkingSpotDirection("left", 445))
	} else if column == 2 {
		directions = append(directions, *newParkingSpotDirection("left", 355))
	} else if column == 3 {
		directions = append(directions, *newParkingSpotDirection("left", 265))
	} else if column == 4 {
		directions = append(directions, *newParkingSpotDirection("left", 175))
	} else if column == 5 {
		directions = append(directions, *newParkingSpotDirection("left", 85))
	}

	directions = append(directions, *newParkingSpotDirection("down", y+5))
	directions = append(directions, *newParkingSpotDirection("left", x+5))

	return &directions
}

func getDirectionsForLeaving() *[]ParkingSpotDirection {
	var directions []ParkingSpotDirection

	directions = append(directions, *newParkingSpotDirection("down", 380))
	directions = append(directions, *newParkingSpotDirection("right", 475))
	directions = append(directions, *newParkingSpotDirection("up", 185))

	return &directions
}

func (p *ParkingSpot) GetArea() *floatgeom.Rect2 {
	return p.area
}

func (p *ParkingSpot) GetNumber() int {
	return p.number
}

func (p *ParkingSpot) GetDirectionsForParking() *[]ParkingSpotDirection {
	return p.directionsForParking
}

func (p *ParkingSpot) GetDirectionsForLeaving() *[]ParkingSpotDirection {
	return p.directionsForLeaving
}

func (p *ParkingSpot) GetIsAvailable() bool {

	return p.isAvailable
}

func (p *ParkingSpot) SetIsAvailable(isAvailable bool) {
	p.isAvailable = isAvailable
}
