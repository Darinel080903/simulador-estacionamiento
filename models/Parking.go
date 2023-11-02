package models
import (
	"sync"
)


type CarQueue struct {
	cars []Car
}

func NewCarQueue() *CarQueue {
	return &CarQueue{
		cars: make([]Car, 0),
	}
}




type ParkingZone struct {
	spots         []*ParkingSpot
	queueCars     *CarQueue
	mu            sync.Mutex
	availableCond *sync.Cond
}


func NewParkingSlot(spots []*ParkingSpot) *ParkingZone {
	p := &ParkingZone{
		spots:     spots,
		queueCars: NewCarQueue(),
	}
	p.availableCond = sync.NewCond(&p.mu)
	return p
}


func (p *ParkingZone) GetSpots() []*ParkingSpot {
	return p.spots
}


func (p *ParkingZone) GetParkingSpotAvailable() *ParkingSpot {
	p.mu.Lock()
	defer p.mu.Unlock()

	for {
		for _, spot := range p.spots {
			if spot.GetIsAvailable() {
				spot.SetIsAvailable(false)
				return spot
			}
		}
		p.availableCond.Wait()
	}
}


func (p *ParkingZone) ReleaseParkingSpot(spot *ParkingSpot) {
	p.mu.Lock()
	defer p.mu.Unlock()

	spot.SetIsAvailable(true)
	p.availableCond.Signal()
}


func (p *ParkingZone) GetQueueCars() *CarQueue {
	return p.queueCars
}