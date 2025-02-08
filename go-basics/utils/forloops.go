package utils

import "fmt"

// for loop
func Solve() {
	var x int = 10
	var y int = 2

	for i:=0; i<x;i++ {
		fmt.Println("y*{i}",y*i)
	}
}