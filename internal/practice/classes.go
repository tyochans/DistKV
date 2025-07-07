package main

import "fmt"

type Counter struct {
	val int
}

func (c Counter) IncrementValue() {
	c.val++
	fmt.Println("Inside value receiver:", c.val)
}

func (c *Counter) IncrementPointer() {
	c.val++
	fmt.Println("Inside pointer receiver:", c.val)
}

func IncrementStandalone(c Counter) {
	c.val++
	fmt.Println("Inside standalone function:", c.val)
}

func main() {
	c := Counter{val: 5}

	c.IncrementValue()
	//	fmt.Println("Inside value receiver:", c.val) :: 6
	fmt.Println("After value receiver:", c.val)
	//: 5

	c.IncrementPointer()
	// fmt.Println("Inside pointer receiver:", c.val):6
	fmt.Println("After pointer receiver:", c.val)
	//:6

	IncrementStandalone(c)
	//	fmt.Println("Inside standalone function:", c.val) : 7
	fmt.Println("After standalone function:", c.val)
	// 6 , because c COunter is a copy of c in main
}
