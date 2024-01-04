package builders

import "fmt"

// Report собираемый билдерами объект
type Report struct {
	header string
	body   string
	footer string
}

func (r Report) PrintReport() {
	fmt.Print(r.header)
	fmt.Print(r.body)
	fmt.Print(r.footer)
}
