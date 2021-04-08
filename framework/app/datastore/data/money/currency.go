package money

type Currency struct {
	Name          string
	Symbol        CurrencySymbol
	ExchangeRates []*ExchangeRate
}

type SymbolLocation bool

const (
	Before SymbolLocation = iota
	After
)

type CurrencySymbol struct {
	Rune     rune
	Location SymbolLocation
}
