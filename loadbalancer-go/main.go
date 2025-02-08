package main

import (
	"fmt"
	"time"
)

func SayHello() {
	fmt.Println("boom")

}

func main() {

	go SayHello()
	time.Sleep(1 * time.Minute)
}
