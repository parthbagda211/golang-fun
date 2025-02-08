package utils

import "fmt"

func Example() {
	defer func() {
        if r := recover(); r != nil {
            fmt.Println("example 4", r)
        }
    }()

	defer fmt.Println("example 1")
	fmt.Println("example 2")
	panic("example 3")
}