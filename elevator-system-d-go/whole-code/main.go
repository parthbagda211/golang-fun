package wholecode

import (
	"fmt"
	"sort"
	"sync"
)

// Directions type and constants
type Directions string

const (
	Up    Directions = "Up"
	Down  Directions = "Down"
	Still Directions = "Still"
)

// Building has Floors and Elevators
type Building struct {
	Floors    []*Floor
	Elevators []*Elevator
}

func NewBuilding() *Building {
	building := &Building{Floors: make([]*Floor, 0)}

	for i := 1; i <= 15; i++ {
		floor := NewFloor(i)
		building.Floors = append(building.Floors, floor)
	}

	for i := 1; i <= 3; i++ {
		elevator := NewElevator(i)
		building.Elevators = append(building.Elevators, elevator)
	}

	return building
}

// ElevatorManager controls all elevators in the building
type ElevatorManager struct {
	Building *Building
}

func NewElevatorManager(building *Building) *ElevatorManager {
	return &ElevatorManager{Building: building}
}

func (em *ElevatorManager) OperateAllElevators() {
	for _, elevator := range em.Building.Elevators {
		// Operate elevators concurrently
		go em.OperateElevator(elevator)
	}
}

func (em *ElevatorManager) OperateElevator(elevator *Elevator) {
	for {
		elevator.Lock()
		if len(elevator.Destinations) == 0 {
			elevator.CurrentDirection = Still
			elevator.Unlock()
			continue
		}

		sort.Ints(elevator.Destinations)
		fmt.Printf("Elevator %d is starting from floor %d and going %s\n", elevator.ID, elevator.CurrentFloor, elevator.CurrentDirection)

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

func (em *ElevatorManager) MoveElevatorUp(elevator *Elevator) {
	for i := 0; i < len(elevator.Destinations); i++ {
		destination := elevator.Destinations[i]

		if destination >= elevator.CurrentFloor {
			fmt.Printf("Elevator %d moving up to floor %d\n", elevator.ID, destination)
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

func (em *ElevatorManager) MoveElevatorDown(elevator *Elevator) {
	for i := len(elevator.Destinations) - 1; i >= 0; i-- {
		destination := elevator.Destinations[i]

		if destination <= elevator.CurrentFloor {
			fmt.Printf("Elevator %d moving down to floor %d\n", elevator.ID, destination)
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

func (em *ElevatorManager) AssignElevator(floor int, direction Directions) *Elevator {
	bestElevator := em.FindClosestElevator(floor, direction)

	if bestElevator != nil {
		bestElevator.AddDestination(floor)
		fmt.Printf("Elevator %d assigned to floor %d with direction %s\n", bestElevator.ID, floor, direction)
	}
	return bestElevator
}

func (em *ElevatorManager) DecideDirection(elevator *Elevator) {
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

func (em *ElevatorManager) FindClosestElevator(floor int, direction Directions) *Elevator {
	var closestElevator *Elevator
	minDist := int(1e9)

	for _, elevator := range em.Building.Elevators {
		elevator.Lock()
		distance := em.CalculateDistance(elevator, floor, direction)

		if distance < minDist {
			minDist = distance
			closestElevator = elevator
		}

		elevator.Unlock()
	}

	return closestElevator
}

func (em *ElevatorManager) CalculateDistance(elevator *Elevator, floor int, direction Directions) int {
	currentFloor := elevator.CurrentFloor
	currentDirection := elevator.CurrentDirection

	if currentDirection == Still || (currentDirection == direction && ((direction == Up && floor > currentFloor)) || (direction == Down && floor < currentFloor)) {
		return abs(floor - currentFloor)
	}

	if (currentDirection == Up && direction == Down) || (currentDirection == Down && direction == Up) {
		if currentDirection == Up {
			return abs(elevator.FarthestDestination()-currentFloor) + abs(elevator.FarthestDestination()-floor)
		} else {
			return abs(elevator.NearestDestination()-currentFloor) + abs(elevator.NearestDestination()-floor)
		}
	}

	return 100
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

// ElevatorPanel represents the button panel inside each elevator
type ElevatorPanel struct {
	PanelID      int
	FloorButtons [15]bool
}

func NewElevatorPanel(panelID int) *ElevatorPanel {
	return &ElevatorPanel{PanelID: panelID, FloorButtons: [15]bool{}}
}

func (ep *ElevatorPanel) AddDestinationFloor(floor int) {
	ep.FloorButtons[floor] = true
}

func (ep *ElevatorPanel) RemoveDestinationFloor(floor int) {
	ep.FloorButtons[floor] = false
}

// Elevator represents a single elevator
type Elevator struct {
	ID             int
	Capacity       int
	CurrentFloor   int
	CurrentDirection Directions
	CurentLoad     int
	ElevatorPanel  *ElevatorPanel
	Destinations   []int
	sync.Mutex
}

func NewElevator(id int) *Elevator {
	return &Elevator{ID: id, Capacity: 10, CurrentFloor: 1, CurrentDirection: Still, CurentLoad: 0, ElevatorPanel: NewElevatorPanel(id)}
}

func (e *Elevator) AddDestination(destinationFloor int) {
	e.Lock()
	e.ElevatorPanel.AddDestinationFloor(destinationFloor)
	e.Destinations = append(e.Destinations, destinationFloor)
	fmt.Printf("Elevator %d received destination floor %d\n", e.ID, destinationFloor)
	e.Unlock()
}

func (e *Elevator) RemoveDestinationFloor(destinationFloor int) {
	e.Lock()
	for i, floor := range e.Destinations {
		if floor == destinationFloor {
			e.Destinations = append(e.Destinations[:i], e.Destinations[i+1:]...)
			e.ElevatorPanel.RemoveDestinationFloor(destinationFloor)
			break
		}
	}
	e.Unlock()
}

func (e *Elevator) UpdateCurrentFloor(newFloor int) {
	e.Lock()
	e.CurrentFloor = newFloor
	e.Unlock()
}

func (e *Elevator) UpdateCurrentDirection(newDirection Directions) {
	e.Lock()
	e.CurrentDirection = newDirection
	e.Unlock()
}

func (e *Elevator) FarthestDestination() int {
	maxFloor := 0

	for _, floor := range e.Destinations {
		if floor > maxFloor {
			maxFloor = floor
		}
	}
	return maxFloor
}

func (e *Elevator) NearestDestination() int {
	minFloor := 100

	for _, floor := range e.Destinations {
		if floor < minFloor {
			minFloor = floor
		}
	}
	return minFloor
}

// Floor represents a floor in the building
type Floor struct {
	Number    int
	HallPanels []*HallPanel
}

func NewFloor(number int) *Floor {
	floor := &Floor{Number: number, HallPanels: make([]*HallPanel, 0)}

	for i := 1; i <= 3; i++ {
		hallPanel := NewHallPanel(i, number)
		floor.HallPanels = append(floor.HallPanels, hallPanel)
	}
	return floor
}

// HallPanel represents a hall panel on each floor for elevator requests
type HallPanel struct {
	PanelID             int
	DirectionInstructions Directions
	SourceFloor         int
}

func NewHallPanel(panelID int, sourceFloor int) *HallPanel {
	return &HallPanel{PanelID: panelID, SourceFloor: sourceFloor, DirectionInstructions: Still}
}

func (h *HallPanel) SetDirectionInstructions(directionInstruction Directions) {
	h.DirectionInstructions = directionInstruction
}

func (h *HallPanel) RequestElevator(manager *ElevatorManager, direction Directions) *Elevator {
	fmt.Printf("Panel %d requesting elevator with direction %s from floor %d\n", h.PanelID, direction, h.SourceFloor)
	return manager.AssignElevator(h.SourceFloor, direction)
}

// Main function to simulate elevator system
func main() {
	building := NewBuilding()
	manager := NewElevatorManager(building)

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		e := building.Floors[1].HallPanels[1].RequestElevator(manager, Up)
		e.AddDestination(6)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		elevator := building.Floors[8].HallPanels[2].RequestElevator(manager, Down)
		elevator.AddDestination(7)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		thirdElevator := building.Floors[3].HallPanels[0].RequestElevator(manager, Up)
		thirdElevator.AddDestination(12)
	}()

	wg.Wait()

	go manager.OperateAllElevators()

	select {}
}
