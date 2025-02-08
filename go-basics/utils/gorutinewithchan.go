package utils

import ( 
	"time"
)

func CalculateSum(a int, b int, c chan int) {
	time.Sleep(200 * time.Millisecond)
	c <- a + b
}