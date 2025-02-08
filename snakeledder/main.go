package main

import (
    "fmt"
    "snakeledder/components"
)

func Run() {
    gameManager := components.GetGameManager()

    player1 := []string{"p1", "p2", "p3"}
    gameManager.StartNewGame(player1)
    fmt.Println("Games started. Check game output above.")
}

func main() {
    Run()
}