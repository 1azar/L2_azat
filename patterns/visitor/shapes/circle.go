package shapes

// Circle Конкретный элемент
// Недостатком посетителя является нарушение инкапсуляции, т.е. определенные поля должны быть экспортируемыми
type Circle struct {
	Radius int
}

// Accept - функция для работы посетителя. Она универсальна для любого посетителя и вносится в код объекта лишь единожды.
func (c *Circle) Accept(v Visitor) {
	v.VisitForCircle(c)
}

// GetType условный метод, нужен для того, чтобы посетитель идентифицировал класс объекта. тк у каждого класса свое поведение
func (c *Circle) GetType() string {
	return "circle"
}
