package shapes

// Rectangle Конкретный элемент
// Недостатком посетителя является нарушение инкапсуляции, т.е. определенные поля должны быть экспортируемыми
type Rectangle struct {
	L int
	B int
}

// Accept - функция для работы посетителя. Она универсальна для любого посетителя и вносится в код объекта лишь единожды.
func (r *Rectangle) Accept(v Visitor) {
	v.VisitForRectangle(r)
}

// GetType условный метод, нужен для того, чтобы посетитель идентифицировал класс объекта. тк у каждого класса свое поведение
func (r *Rectangle) GetType() string {
	return "rectangle"
}
