package shapes

// Visitor - интерфейс посетителя, определяет методы для объектов, которым добавляет поведение
type Visitor interface {
	VisitForSquare(*Square)
	VisitForCircle(*Circle)
	VisitForRectangle(*Rectangle)
}
