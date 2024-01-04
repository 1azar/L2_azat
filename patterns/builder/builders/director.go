package builders

// Director
// Директор знает как правильно вызывать методы билдера, он наш помощник
type Director struct {
	builder ReportBuilder
}

func NewDirector(b ReportBuilder) *Director {
	return &Director{
		builder: b,
	}
}

// SetBuilder присваиваем билдер директору.
// "Даем рабочего, которым он будет командовать"
func (d *Director) SetBuilder(b ReportBuilder) {
	d.builder = b
}

func (d *Director) BuildReport() Report {
	d.builder.SetHeader("AbstractBuyer") // имя может быть получено с помощью определенной сложной логики
	d.builder.SetBody()
	d.builder.SetFooter()
	return d.builder.GetReport()
}
