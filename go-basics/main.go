package main

import (
	"fmt"
	
	// "go-basics/utils"
	_ "go-basics/utils"
	"sync"
)

// if conditions
var counter int 
var mu sync.Mutex
var wg sync.WaitGroup

func IncrementCounter(wg *sync.WaitGroup) {
	defer wg.Done()
	mu.Lock()
	counter++
	mu.Unlock()
}

func SayHell(s string) {
	defer wg.Done()
	fmt.Println(s)
	wg.Add(1)
	
}

func main() {
	// fmt.Println("Hello, World!")
	// utils.Solve()
	// res := utils.Sum(2, 3)
	// fmt.Println(res)
	// c := utils.Circle{Radius: 5}
	// fmt.Println(c.Area())
	// res, err := utils.Divide(10, 0)
	// if err != nil {
	// 	fmt.Println(err)
	// } else {
	// 	fmt.Println(res)
	// }
	// // utils.WebServer()

	// result := func(a  ,b int) int {
	// 	return a-b
	// }(5,1)
	// fmt.Println(result)

	// summ := utils.Sum(1, 2, 4, 5)
    // fmt.Println("Total Sum:", summ)

	// defer fmt.Println("main 1")
    // utils.Example()
	// fmt.Println("main 2")

	// res := utils.Factorial(5)
	// fmt.Println(res)

	// person := utils.Person{FirstName: "John", LastName: "Doe", Age: 25}
	// person.FullName()
	// person.Greet()
	// person.DisplayAge()

	// go utils.SayHello("parth")
	// go utils.SayHello("bagda")

	// ch := make(chan int)
	// go utils.CalculateSum(2, 3, ch)
	// go utils.CalculateSum(5, 7, ch)

	// result1 := <-ch


	// fmt.Println(result1)

    // var wg sync.WaitGroup
	// for i := 0; i < 1000; i++ {
	// 	wg.Add(1)
	// 	go IncrementCounter(&wg)
	// }
	// wg.Wait()
	// fmt.Println("Counter Value:", counter)

	// arr := []int{1,2,3}

	// arr = append(arr,4)
	// fmt.Println(arr)
   

	go SayHell("HELL")
	





}
