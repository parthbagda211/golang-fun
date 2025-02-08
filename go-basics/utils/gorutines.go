package utils

import (
	"fmt"
	"time"
)

func SayHello(name string) {
	time.Sleep(200 * time.Millisecond)
	fmt.Println("Hello", name)
}