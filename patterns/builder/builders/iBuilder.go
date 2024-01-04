package builders

type BuilderTypes string

const (
	Buyer  BuilderTypes = "buyer"
	Seller              = "seller"
)

// ReportBuilder - поведение для билдеров
type ReportBuilder interface {
	SetHeader(name string)
	SetBody()
	SetFooter()
	GetReport() Report
}

func GetBuilder(builderType BuilderTypes) ReportBuilder {
	switch builderType {
	case Buyer:
		return NewBuyerRepBuilder()
	case Seller:
		return NewSellerRepBuilder()
	default:
		panic("unknown builder type")
	}
}
