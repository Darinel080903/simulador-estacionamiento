package models

import (
	"fmt"
	"image/color"
	"log"
	"math"
	"sync"
	"time"

	"github.com/oakmound/oak/v4/alg/floatgeom"
	"github.com/oakmound/oak/v4/entities"
	"github.com/oakmound/oak/v4/render"
	"github.com/oakmound/oak/v4/scene"
)

const (
	entranceSpot = 185.00
	speed        = 10
)

type Car struct {
	area   floatgeom.Rect2
	entity *entities.Entity
	mu     sync.Mutex
}

func NewCar(ctx *scene.Context) *Car {
	area := floatgeom.NewRect2(445, -20, 465, 0)

	carRender, err := render.LoadSprite("assets/images/car.png")
	if err != nil {
		log.Fatal(err)
	}

	entity := entities.New(ctx, entities.WithRect(area), entities.WithColor(color.RGBA{51, 222, 0, 13}), entities.WithRenderable(carRender), entities.WithDrawLayers([]int{1}))

	return &Car{
		area:   area,
		entity: entity,
	}
}

func (c *Car) Enqueue(manager *CarManager) {

	for c.Y() < 145 {
		if !c.carisCollision("down", manager.GetCars()) {
			c.ShiftY(1)
			time.Sleep(speed * time.Millisecond)
		}
	}

}

func (c *Car) JoinDoor(manager *CarManager) {
	for c.Y() < entranceSpot {
		if !c.carisCollision("down", manager.GetCars()) {
			c.ShiftY(1)
			time.Sleep(speed * time.Millisecond)
		}
	}
}

func (c *Car) ExitDoor(manager *CarManager) {
	for c.Y() > 145 {
		if !c.carisCollision("up", manager.GetCars()) {
			c.ShiftY(-1)
			time.Sleep(speed * time.Millisecond)
		}
	}
}

func (c *Car) ParkZone(spot *ParkingSpot, manager *CarManager) {
	for index := 0; index < len(*spot.GetDirectionsForParking()); index++ {
		directions := *spot.GetDirectionsForParking()
		fmt.Println("Carro gira a: " + directions[index].Direction)
		fmt.Println("Se dirige a: " + fmt.Sprintf("%f", directions[index].Point))
		if directions[index].Direction == "right" {
			for c.X() < directions[index].Point {
				if !c.carisCollision("right", manager.GetCars()) {
					c.ShiftX(1)
					time.Sleep(speed * time.Millisecond)
				}
			}
		} else if directions[index].Direction == "down" {
			for c.Y() < directions[index].Point {
				if !c.carisCollision("down", manager.GetCars()) {
					c.ShiftY(1)
					time.Sleep(speed * time.Millisecond)
				}
			}
		} else if directions[index].Direction == "left" {
			for c.X() > directions[index].Point {
				if !c.carisCollision("left", manager.GetCars()) {
					c.ShiftX(-1)
					time.Sleep(speed * time.Millisecond)
				}
			}
		} else if directions[index].Direction == "up" {
			for c.Y() > directions[index].Point {
				if !c.carisCollision("up", manager.GetCars()) {
					c.ShiftY(-1)
					time.Sleep(speed * time.Millisecond)
				}
			}
		}
	}
}

func (c *Car) LeaveSlot(spot *ParkingSpot, manager *CarManager) {
	for index := 0; index < len(*spot.GetDirectionsForLeaving()); index++ {
		directions := *spot.GetDirectionsForLeaving()
		if directions[index].Direction == "left" {

			for c.X() > directions[index].Point {
				if !c.carisCollision("left", manager.GetCars()) {
					c.ShiftX(-1)
					time.Sleep(speed * time.Millisecond)
				}
			}
		} else if directions[index].Direction == "right" {
			for c.X() < directions[index].Point {
				if !c.carisCollision("right", manager.GetCars()) {
					c.ShiftX(1)
					time.Sleep(speed * time.Millisecond)
				}
			}
		} else if directions[index].Direction == "up" {
			for c.Y() > directions[index].Point {
				if !c.carisCollision("up", manager.GetCars()) {
					c.ShiftY(-1)
					time.Sleep(speed * time.Millisecond)
				}
			}
		} else if directions[index].Direction == "down" {
			for c.Y() < directions[index].Point {
				if !c.carisCollision("down", manager.GetCars()) {
					c.ShiftY(1)
					time.Sleep(speed * time.Millisecond)
				}
			}
		}
	}
}

func (c *Car) LeaveSpot(manager *CarManager) {
	spotX := c.X()
	for c.X() > spotX-30 {
		if !c.carisCollision("left", manager.GetCars()) {
			c.ShiftX(-1)
			time.Sleep(speed * time.Millisecond)
		}
		time.Sleep(20 * time.Millisecond)
	}
}

func (c *Car) GoAway(manager *CarManager) {
	for c.Y() > -20 {
		if !c.carisCollision("up", manager.GetCars()) {
			c.ShiftY(-1)
			time.Sleep(speed * time.Millisecond)
		}
	}
}

func (c *Car) ShiftY(dy float64) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.entity.ShiftY(dy)
}

func (c *Car) ShiftX(dx float64) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.entity.ShiftX(dx)
}

func (c *Car) X() float64 {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.entity.X()
}

func (c *Car) Y() float64 {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.entity.Y()
}

func (c *Car) Remove() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.entity.Destroy()
}

func (c *Car) carisCollision(direction string, cars []*Car) bool {
	const minDistance = 30.0

	for _, car := range cars {
		switch direction {
		case "left":
			if c.isWithinDistance(car, minDistance, true) && c.Y() == car.Y() && c.X() > car.X() {
				return true
			}
		case "right":
			if c.isWithinDistance(car, minDistance, true) && c.Y() == car.Y() && c.X() < car.X() {
				return true
			}
		case "up":
			if c.isWithinDistance(car, minDistance, false) && c.X() == car.X() && c.Y() > car.Y() {
				return true
			}
		case "down":
			if c.isWithinDistance(car, minDistance, false) && c.X() == car.X() && c.Y() < car.Y() {
				return true
			}
		}
	}
	return false
}

func (c *Car) CarisCollision(direction string, cars []*Car) bool {
	const minDistance = 30.0

	for _, car := range cars {
		switch direction {
		case "left":
			if c.isWithinDistance(car, minDistance, true) && c.Y() == car.Y() && c.X() > car.X() {
				return true
			}
		case "right":
			if c.isWithinDistance(car, minDistance, true) && c.Y() == car.Y() && c.X() < car.X() {
				return true
			}
		case "up":
			if c.isWithinDistance(car, minDistance, false) && c.X() == car.X() && c.Y() > car.Y() {
				return true
			}
		case "down":
			if c.isWithinDistance(car, minDistance, false) && c.X() == car.X() && c.Y() < car.Y() {
				return true
			}
		}
	}
	return false
}

func (c *Car) isWithinDistance(car *Car, distance float64, horizontal bool) bool {
	if horizontal {
		return math.Abs(c.X()-car.X()) < distance
	}
	return math.Abs(c.Y()-car.Y()) < distance
}
