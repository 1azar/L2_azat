package shapes

// Shape - интерфейс, который реализуют все фигуры
type Shape interface {
	GetType() string
	Accept(visitor Visitor)
}
