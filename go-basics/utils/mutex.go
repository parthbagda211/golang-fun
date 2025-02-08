package utils

import (
	"sync"
)

var counter int 
var mu sync.Mutex

func IncrementCounter(wg *sync.WaitGroup) {
	defer wg.Done()
	mu.Lock()
	counter++
	mu.Unlock()
}