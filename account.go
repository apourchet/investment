package invt

type Account struct {
	Id              string
	Name            string
	Currency        string
	Balance         float64
	UnrealizedPl    float64
	RealizedPl      float64
	OpenPositions   map[string]*OpenPosition
	MarginRate      float64
	MarginUsed      float64
	MarginAvailable float64
	OpenTrades      interface{}
	OpenOrders      interface{}
}

func NewAccount(balance float64) *Account {
	a := &Account{}
	a.Balance = balance
	a.MarginRate = 0.02
	a.OpenPositions = make(map[string]*OpenPosition)
	return a
}
