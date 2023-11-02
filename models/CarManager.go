package models

import "sync"

type CarManager struct {
	Cars  []*Car
	Mutex sync.Mutex
}

func NewCarManager() *CarManager {
	return &CarManager{
		Cars: make([]*Car, 0),
	}
}

func (cm *CarManager) AddCar(car *Car) {
	cm.Mutex.Lock()
	defer cm.Mutex.Unlock()
	cm.Cars = append(cm.Cars, car)
}

func (cm *CarManager) RemoveCarFromScene(car *Car) {
	cm.Mutex.Lock()
	defer cm.Mutex.Unlock()
	for i, c := range cm.Cars {
		if c == car {
			cm.Cars = append(cm.Cars[:i], cm.Cars[i+1:]...)
			break
		}
	}
}

func (cm *CarManager) GetCars() []*Car {
	cm.Mutex.Lock()
	defer cm.Mutex.Unlock()
	carsCopy := make([]*Car, len(cm.Cars))
	copy(carsCopy, cm.Cars)
	return carsCopy
}

