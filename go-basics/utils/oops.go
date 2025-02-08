package utils

import "fmt"

type Person struct {
	FirstName string
	LastName  string
	Age 	  int
}

func (p Person) FullName() string {
	return p.FirstName + " " + p.LastName
}

func (p Person) Greet() {
	fmt.Println("Hello", p.FullName())
}

func (p Person) DisplayAge() {
	fmt.Println(p.FullName(), "is", p.Age, "years old")
}
