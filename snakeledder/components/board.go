package components

type Board struct{
	Size int 
	Snakes []*Snake
	Ladder []*Ladder
}

func NewBoard() *Board{
	board := &Board{
		Size: 100,
		Snakes: []*Snake{},
		Ladder: []*Ladder{},
	}
	board.initBoard()
	return board
}

func (b*Board) initBoard(){
   b.Snakes = append(b.Snakes, NewSnake(16,6),NewSnake(48,26), NewSnake(64, 60), NewSnake(93, 73))
   b.Ladder = append(b.Ladder, NewLeadder(1, 38), NewLeadder(4, 14), NewLeadder(9, 31), NewLeadder(21, 42),
   NewLeadder(28, 84), NewLeadder(51, 67), NewLeadder(80, 99))
  
}

func (b * Board) GetNewPos(position int ) int {
	for _,snake := range b.Snakes {
		if snake.Start == position {
			return snake.End
		}
	}
	for _,ladder := range b.Ladder {
		if ladder.Start == position {
			return ladder.End
		}
	}

	return position
}