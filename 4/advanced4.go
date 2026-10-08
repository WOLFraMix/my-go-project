package main

import (
	"fmt"
	"math"
)

type Shape interface {
	Area() float64
}

type Rectangle struct {
	Width  float64
	Height float64
}

type Circle struct {
	Radius float64
}

func (r Rectangle) Area() float64 {
	return r.Width * r.Height
}

func (c Circle) Area() float64 {
	return math.Pi * c.Radius * c.Radius
}

func main() {
	s := []Shape{Rectangle{Width: 3, Height: 4}, Circle{Radius: 5}}
	for _, v := range s {
		fmt.Printf("Площадь: %0.2f\n", v.Area())
	}
}
