package data

import (
	"errors"
	"math"
)

// Currency represents a currency with its 3-letter code, symbol (if available), and full name
type Currency struct {
	Code               string
	Symbol             string
	Name               string
	SmallestUnitFactor int // The factor to convert to the smallest unit
}

var currencies = []Currency{
	{"USD", "$", "United States Dollar", 100},   // Smallest unit is cent (1 USD = 100 cents)
	{"EUR", "€", "Euro", 100},                   // Smallest unit is cent (1 EUR = 100 cents)
	{"JPY", "¥", "Japanese Yen", 1},             // Smallest unit is 1 yen
	{"GBP", "£", "British Pound Sterling", 100}, // Smallest unit is penny (1 GBP = 100 pence)
	{"AUD", "A$", "Australian Dollar", 100},     // Smallest unit is cent (1 AUD = 100 cents)
	{"CAD", "C$", "Canadian Dollar", 100},       // Smallest unit is cent (1 CAD = 100 cents)
	{"CHF", "Fr", "Swiss Franc", 100},           // Smallest unit is centime (1 CHF = 100 centimes)
	{"CNY", "¥", "Chinese Yuan", 100},           // Smallest unit is fen (1 CNY = 100 fen)
	{"SEK", "kr", "Swedish Krona", 100},         // Smallest unit is öre (1 SEK = 100 öre)
	{"NZD", "NZ$", "New Zealand Dollar", 100},   // Smallest unit is cent (1 NZD = 100 cents)
	{"KRW", "₩", "South Korean Won", 1},         // Smallest unit is 1 won
	{"SGD", "S$", "Singapore Dollar", 100},      // Smallest unit is cent (1 SGD = 100 cents)
	{"NOK", "kr", "Norwegian Krone", 100},       // Smallest unit is øre (1 NOK = 100 øre)
	{"MXN", "Mex$", "Mexican Peso", 100},        // Smallest unit is centavo (1 MXN = 100 centavos)
	{"INR", "₹", "Indian Rupee", 100},           // Smallest unit is paise (1 INR = 100 paise)
	{"HKD", "HK$", "Hong Kong Dollar", 100},     // Smallest unit is cent (1 HKD = 100 cents)
	{"BRL", "R$", "Brazilian Real", 100},        // Smallest unit is centavo (1 BRL = 100 centavos)
	{"ZAR", "R", "South African Rand", 100},     // Smallest unit is cent (1 ZAR = 100 cents)
	{"TRY", "₺", "Turkish Lira", 100},           // Smallest unit is kuruş (1 TRY = 100 kuruş)
	{"RON", "L", "Romanian Leu", 100},           // Smallest unit is ban (1 RON = 100 bani)
}

// InSmallestCurrencyUnit converts an amount in float64 to the smallest unit of the given currency
func InSmallestCurrencyUnit(code string, amount float64) (int64, error) {
	code = uppercaseTrim(code)

	var currencyFound bool
	var smallestUnitFactor int

	for _, curr := range currencies {
		if curr.Code == code {
			currencyFound = true
			smallestUnitFactor = curr.SmallestUnitFactor
			break
		}
	}

	if !currencyFound {
		return 0, errors.New("currency not found")
	}

	smallestUnit := int64(math.Round(amount * float64(smallestUnitFactor)))
	return smallestUnit, nil
}

func uppercaseTrim(s string) string {
	return string([]rune(s)[:3]) // Truncate to 3 characters and convert to uppercase
}
