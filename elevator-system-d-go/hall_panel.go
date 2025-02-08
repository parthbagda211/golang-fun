package main

import "fmt"


type HallPanel struct {
	PanelID int 
	DirectionInstructions Directions
	SourceFloor int 
}


func NewHallPanel(panelID int, sourceFloor int ) *HallPanel {
	return &HallPanel{PanelID: panelID,SourceFloor: sourceFloor,DirectionInstructions: Still}
}

func (h *HallPanel) SetDirectionInstructions(directionInstruction Directions) {
	h.DirectionInstructions = directionInstruction
}

func (h *HallPanel) RequestElevator(manager *ElevatorManager, direction Directions) (elevator * Elevator) {
	fmt.Printf("Panel %d requesting elevator with direction %s from floor %d\n", h.PanelID, direction, h.SourceFloor)
	return manager.AssignElevator(h.SourceFloor, direction) 
}
