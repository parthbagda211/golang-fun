package main

import (
	"fmt"
	"sort"
)


type ElevatorManager struct {
	Building *Building
}


func NewElevatorManager(building *Building) *ElevatorManager {
	return &ElevatorManager{Building: building}
}

func (em * ElevatorManager) OperateAllElevators() {
	for _,elevator := range em.Building.Elevators {
		// why here i am using go routine , bcoz elevator can be oprate paralleliy for each building 
		go em.OperateElevator(elevator)
	}
}


func (em * ElevatorManager) OperateElevator(elevator *Elevator) {
	for {
		elevator.Lock()
		if len(elevator.Destinations) == 0 {
			elevator.CurrentDirection = Still
			elevator.Unlock()
			continue
		}
        
		sort.Ints(elevator.Destinations)
		fmt.Printf("elevator %d is starting from %d and going to %s\n", elevator.ID, elevator.CurrentFloor, elevator.CurrentDirection)
        
		if elevator.CurrentDirection == Up {
			em.MoveElevatorUp(elevator)
		} else if elevator.CurrentDirection == Down {
			em.MoveElevatorDown(elevator)
		} else {
			em.DecideDirection(elevator)
		}
		elevator.Unlock()
	}
}

func (em * ElevatorManager) MoveElevatorUp(elevator * Elevator) {
	for i:=0; i<len(elevator.Destinations);i++ {
		destination := elevator.Destinations[i]

		if destination >= elevator.CurrentFloor {
			fmt.Printf("elevator %d moving up to floor %d\n", elevator.ID, destination)
			elevator.UpdateCurrentFloor(destination)
			elevator.RemoveDestinationFloor(destination)
			i--
		}
	}

	if len(elevator.Destinations) == 0 {
		elevator.UpdateCurrentDirection(Still)
	} else {
		elevator.UpdateCurrentDirection(Down)
	}
}


func (em *ElevatorManager) MoveElevatorDown(elevator * Elevator) {
	for i:= len(elevator.Destinations)-1; i>=0 ; i-- {
		destination := elevator.Destinations[i]

		if destination <= elevator.CurrentFloor {
			fmt.Printf("elevator %d moving down to floor %d\n", elevator.ID, destination)
			elevator.UpdateCurrentFloor(destination)
			elevator.RemoveDestinationFloor(destination)
		}
	}

	if len(elevator.Destinations) == 0 {
		elevator.UpdateCurrentDirection(Still)
	} else {
		elevator.UpdateCurrentDirection(Up)
	}
}

func (em * ElevatorManager) AssignElevator(floor int , direction Directions) (bestElevator *Elevator) {
	bestElevator = em.FindClosestElevator(floor,direction)

	if bestElevator != nil {
		bestElevator.AddDestination(floor)
		fmt.Printf("elevator %d assigned to floor %d with direction %s\n",bestElevator.ID,floor,direction)

	}
	return bestElevator
}

func (em * ElevatorManager) DecideDirection(elevator * Elevator) {
	currentFloor := elevator.CurrentFloor

	if len(elevator.Destinations) == 0 {
		return
	}

	nearestDestination := elevator.Destinations[0]
	if nearestDestination > currentFloor {
		elevator.UpdateCurrentDirection(Up)
		em.MoveElevatorUp(elevator)
	} else {
		elevator.UpdateCurrentDirection(Down)
		em.MoveElevatorDown(elevator)
	}
}


func (em * ElevatorManager) FindClosestElevator(floor int , direction Directions) *Elevator {
	var closestElevator *Elevator
	minDist :=int(1e9)

	for _,elevator := range em.Building.Elevators {
		elevator.Lock()
		distance  := em.CalculateDistance(elevator,floor,direction)

		if distance<minDist {
			minDist = distance
			closestElevator = elevator
		}

		elevator.Unlock()
	}

	return closestElevator
}

func (em * ElevatorManager) CalculateDistance(elevator * Elevator, floor int , direction Directions) int {
	currentFloor := elevator.CurrentFloor
	currentDirection := elevator.CurrentDirection

	if currentDirection == Still || (currentDirection == direction && ((direction == Up && floor > currentFloor)) || (direction == Down && floor < currentFloor)){
		return abs(floor-currentFloor)
	}

	if (currentDirection == Up && direction == Down ) || (currentDirection == Down && direction == Up) {
		if currentDirection == Up {
			return abs(elevator.FarthestDestination()-currentFloor) + abs(elevator.FarthestDestination()-floor)
		} else {
			return abs(elevator.NearestDestination()-currentFloor) + abs(elevator.NearestDestination()-floor)
		}
	}

	return 100
}

func abs(x int ) int {
	if x < 0 {
		return -x
	}
	return x
}