package main

import "fmt"

type Point struct{ X, Y float64 }

func (p Point) add() {
	fmt.Println("ok")

}
func (p Point) ad() {
	fmt.Println("good")

}

func main() {
	var p Point
	p.add()
	p.ad()
}
