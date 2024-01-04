package builders

import "fmt"

// SellerReport - конкретный билдер, реализовывает интерфейс iBuilder.ReportBuilder
type SellerReport struct {
	header string
	body   string
	footer string
}

func NewSellerRepBuilder() *SellerReport {
	return &SellerReport{}
}

func (s *SellerReport) SetHeader(name string) {
	s.header = fmt.Sprintf(
		"%s совершил покупку\n"+
			"===========",
		name,
	)
}

func (s *SellerReport) SetBody() {
	s.body = GetPurchaseInfo()
}

func (s *SellerReport) SetFooter() {
	s.footer = `===========
<информация для продавца>
`
}

func (s *SellerReport) GetReport() Report {
	return Report{
		header: s.header,
		body:   s.body,
		footer: s.footer,
	}
}
