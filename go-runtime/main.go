package main

import (
	"fmt"
	"sync"
	"time"
	"math/rand"
)

var wg sync.WaitGroup

func SayHello(s string) {
	// Print when the goroutine starts
	fmt.Printf("[%s] %s started\n", time.Now().Format("15:04:05.000000000"), s)

	defer wg.Done() // Decrement the counter once the function finishes
	// Simulate some work with a random sleep
	time.Sleep(time.Millisecond * time.Duration(100+rand.Intn(400)))
	// Print when the goroutine finishes
	fmt.Printf("[%s] %s finished\n", time.Now().Format("15:04:05.000000000"), s)

}

func main() {
	// Increment the WaitGroup counter before launching the goroutines
	wg.Add(5)

	// Launch multiple goroutines
	go SayHello("Hello, 1")
	go SayHello("Hello, 2")
	go SayHello("Hello, 3")
	go SayHello("Hello, 4")
	go SayHello("Hello, 5")

	// Wait for all goroutines to finish
	wg.Wait() // Block until all goroutines call Done()
}
