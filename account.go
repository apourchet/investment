package invt

type Account struct {
	Id            string
	Name          string
	Currency      string
	Balance       float64
	RealizedPl    float64
	OpenPositions map[string]*OpenPosition
	MarginRate    float64
	OpenTrades    interface{}
	OpenOrders    interface{}
}

func NewAccount(balance float64) *Account {
	a := &Account{}
	a.Balance = balance
	a.MarginRate = 0.02
	a.OpenPositions = make(map[string]*OpenPosition)
	return a
}

func (a *Account) MarginUsed() float64 {
	// marginrate * exposure
	// TODO
	return a.Balance
}

func (a *Account) MarginAvailable(qc *QuoteContext) float64 {
	// balance - marginrate * exposure
	return a.Balance - a.MarginRate*a.Exposure(qc)
}

func (a *Account) Exposure(qc *QuoteContext) float64 {
	exposure := 0.
	for _, o := range a.OpenPositions {
		if o.Side == SIDE_BUY {
			exposure += qc.Get(o.InstrumentId).Ask * o.FloatUnits()
		} else {
			exposure += qc.Get(o.InstrumentId).Bid * o.FloatUnits()
		}
	}
	return exposure
}

func (a *Account) UnrealizedPl(qc *QuoteContext) float64 {
	pl := 0.
	for _, o := range a.OpenPositions {
		if o.Side == SIDE_BUY {
			pl += o.Price - qc.Get(o.InstrumentId).Bid
		} else {
			pl += o.Price - qc.Get(o.InstrumentId).Ask
		}
	}
	return pl
}
