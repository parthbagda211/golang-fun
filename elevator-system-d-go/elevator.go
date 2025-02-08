package main

import (
	"fmt"
	"sync"
)

type Elevator struct {
	ID int
	Capacity int 
	CurrentFloor int 
	CurrentDirection Directions
	CurentLoad int 
	ElevatorPanel *ElevatorPanel
    Destinations []int 
	sync.Mutex
}

func NewElevator(id int) *Elevator {
	return &Elevator{ID: id,Capacity: 10,CurrentFloor: 1,CurrentDirection: Still,CurentLoad: 0,ElevatorPanel: NewElevatorPanel(id)}
}

func (e *Elevator) AddDestination(destinationFloor int) {
	e.Lock()
	e.ElevatorPanel.AddDestinationFloor(destinationFloor)
	e.Destinations = append(e.Destinations, destinationFloor)
	fmt.Printf("Elevator %d received destination floor %d\n", e.ID, destinationFloor)
	e.Unlock()
}

func (e *Elevator) RemoveDestinationFloor(destinationFloor int ) {
	e.Lock()
	for i,floor := range e.Destinations {
		if floor==destinationFloor {
			e.Destinations = append(e.Destinations[:i],e.Destinations[i+1:]... )
			e.ElevatorPanel.RemoveDestinationFloor(destinationFloor)
			break
		}
	}
	e.Unlock()
}

func (e *Elevator) UpdateCurrentFloor(newFloor int ){
	e.Lock()
	e.CurrentFloor = newFloor
	e.Unlock()
}

func (e * Elevator) UpdateCurrentDirection(newDirection Directions){
	e.Lock()
    e.CurrentDirection = newDirection
	e.Unlock()
}

func (e * Elevator) FarthestDestination() int {
	maxFloor :=0

	for _,floor := range e.Destinations {
		if floor > maxFloor {
			maxFloor = floor
		}
	}
	return maxFloor
}

func (e * Elevator) NearestDestination() int {
	minFloor :=100

	for _,floor := range e.Destinations {
		if floor < minFloor {
			minFloor = floor
		}
	}
	return minFloor
}