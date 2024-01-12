package enum

type Currency string

const (
	CurrencyUSD Currency = "USD"
	CurrencyEUR Currency = "EUR"
	CurrencyGBP Currency = "GBP"
)

func (e Currency) String() string {
	return string(e)
}

func DecodeCurrency(code string) Currency {
	switch code {
	case "USD":
		return CurrencyUSD
	case "EUR":
		return CurrencyEUR
	case "GBP":
		return CurrencyGBP
	default:
		return ""
	}
}
