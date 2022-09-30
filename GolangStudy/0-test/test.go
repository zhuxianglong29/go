package main

import "fmt"

func main() {
	var red, green, blue uint16 = 1, 11, 34
	fmt.Print(red, green, blue)
	fmt.Printf("%x %x %v", red, green, blue)
}
