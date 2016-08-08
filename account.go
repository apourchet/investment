package invt

import (
	"fmt"

	pb "github.com/apourchet/investment/protos"
)

type Account struct {
	Id              string
	Name            string
	Currency        string
	Balance         float64
	UnrealizedPl    float64
	RealizedPl      float64
	OpenPositions   map[pb.InstrumentID_ID]*OpenPosition
	MarginRate      float64
	MarginUsed      float64
	MarginAvailable float64
	OpenTrades      interface{}
	OpenOrders      interface{}
}

type OpenPosition struct {
	Instrument pb.InstrumentID_ID
	Units      int32
	Price      float64
	Side       pb.OrderSide
}

func CreateNewAccount() *Account {
	a := &Account{}
	a.Balance = 10000
	a.OpenPositions = make(map[pb.InstrumentID_ID]*OpenPosition)
	return a
}

func (pos *OpenPosition) FloatUnits() float64 {
	return float64(pos.Units)
}

func (pos *OpenPosition) Value() float64 {
	return pos.FloatUnits() * pos.Price
}

func (pos *OpenPosition) SwitchSide() {
	if pos.Side == pb.OrderSide_BUY {
		pos.Side = pb.OrderSide_SELL
	} else {
		pos.Side = pb.OrderSide_BUY
	}
}

func (pos *OpenPosition) SplitPosition(units int32) *OpenPosition {
	if units > pos.Units {
		panic("Cannot split position with this many units")
	}
	pos.Units -= units

	return &OpenPosition{pos.Instrument, units, pos.Price, pos.Side}
}

func (pos *OpenPosition) String() string {
	return fmt.Sprintf("u: %d\np: %f\ns: %d", pos.Units, pos.Price, pos.Side)
}

func (a *Account) ClosePosition(pos *OpenPosition, price float64) {
	fmt.Println("Closing Position: \n" + pos.String())
	a.Balance += pos.Value() // Gain value of position
	if pos.Side == pb.OrderSide_BUY {
		a.Balance += pos.FloatUnits() * (price - pos.Price) // Gain delta
	} else {
		a.Balance += pos.FloatUnits() * (pos.Price - price) // Gain delta
	}
	fmt.Println("New Balance: ", a.Balance)
}

func (a *Account) MergePositions(from, to *OpenPosition) {
	if from.Side == to.Side {
		fmt.Println("Widening position")
		a.Balance -= from.Value()
		totalUnits := from.Units + to.Units
		totalValue := from.Value() + to.Value()
		avgPrice := totalValue - float64(totalUnits)

		to.Price = avgPrice
		to.Units = totalUnits
	} else {
		if from.Units == to.Units {
			a.ClosePosition(to, from.Price)
			delete(a.OpenPositions, to.Instrument)
		} else if to.Units > from.Units {
			fmt.Println("Reducing position")
			toclose := to.SplitPosition(from.Units)
			a.ClosePosition(toclose, from.Price)
		} else if from.Units > to.Units {
			fmt.Println("Flipping position")
			a.ClosePosition(to, from.Price)
			delete(a.OpenPositions, to.Instrument)
			from.Units -= to.Units
			a.OpenNewPosition(from)
		}
	}
}

func (a *Account) OpenNewPosition(pos *OpenPosition) {
	fmt.Println("Opening new position")
	a.Balance -= pos.Value()
	a.OpenPositions[pos.Instrument] = pos
	fmt.Println("Balance: ", a.Balance)
}

func (a *Account) ProcessOrder(o *pb.Order) {
	if o.Type == pb.OrderType_MARKET {
		fmt.Println("Market Order")
		pos := &OpenPosition{o.Instrument, o.Units, o.Price, o.Side}
		if other, ok := a.OpenPositions[pos.Instrument]; ok {
			a.MergePositions(pos, other)
		} else {
			a.OpenNewPosition(pos)
		}
	}
}
