package builders

import "fmt"

// BuyerReport - конкретный билдер, реализовывает интерфейс iBuilder.ReportBuilder
type BuyerReport struct {
	header string
	body   string
	footer string
}

func NewBuyerRepBuilder() *BuyerReport {
	return &BuyerReport{}
}

func (b *BuyerReport) SetHeader(name string) {
	b.header = fmt.Sprintf(
		"УВАЖАЕМЫЙ, %s, ЭТО ВАШ ЧЕК!\n"+
			"===========",
		name,
	)
}

func (b *BuyerReport) SetBody() {
	b.body = GetPurchaseInfo()
}

func (b *BuyerReport) SetFooter() {
	b.footer = `===========
СПАСИБО ЗА ПОКУПКУ!
`
}

func (b *BuyerReport) GetReport() Report {
	return Report{
		header: b.header,
		body:   b.body,
		footer: b.footer,
	}
}
