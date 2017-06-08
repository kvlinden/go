package main

import (
	"fmt"
)

func main() {
	c := make(chan int)
	go func() { c <- 3 }()
	fmt.Println(<-c)
}
