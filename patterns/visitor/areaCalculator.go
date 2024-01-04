package main

import (
	"L2_azat/patterns/visitor/shapes"
	"fmt"
	"math"
)

// AreaCalculator - Конкретный посетитель определяющий площадь фигур
type AreaCalculator struct {
}

// Поведение для объектов типа square.Square
func (a AreaCalculator) VisitForSquare(square *shapes.Square) {
	fmt.Printf("Площадь для объекта квадрата %v : %d\n", square, square.A*square.A)
}

// Поведение для объектов типа circle.Circle
func (a AreaCalculator) VisitForCircle(circle *shapes.Circle) {
	fmt.Printf("Площадь для объекта круга %v : %.2f\n", circle, math.Pi*math.Pow(float64(circle.Radius), 2))
}

// Поведение для объектов типа rectangle.Rectangle
func (a AreaCalculator) VisitForRectangle(rectangle *shapes.Rectangle) {
	fmt.Printf("Площадь для объекта прямоугольника %v : %d\n", rectangle, rectangle.B*rectangle.L)
}
