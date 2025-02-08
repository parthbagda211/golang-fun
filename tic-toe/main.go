package main

import "tic-toe/utils"
func Run() {
	player1 := utils.NewPlayer("Player 1", 'X')
	player2 := utils.NewPlayer("Player 2", 'O')

	game := utils.NewGame(player1, player2)
	game.Play()
}

func main(){
	Run()
}
