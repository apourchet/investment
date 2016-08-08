package invt

import "fmt"

type OpenPosition struct {
	InstrumentId string
	Units        int32
	Price        float64
	Side         string
}

func (pos *OpenPosition) FloatUnits() float64 {
	return float64(pos.Units)
}

func (pos *OpenPosition) Value() float64 {
	return pos.FloatUnits() * pos.Price
}

func (pos *OpenPosition) SwitchSide() {
	if pos.Side == "buy" {
		pos.Side = "sell"
	} else {
		pos.Side = "buy"
	}
}

func (pos *OpenPosition) SplitPosition(units int32) *OpenPosition {
	if units > pos.Units {
		panic("Cannot split position with this many units")
	}
	pos.Units -= units

	return &OpenPosition{pos.InstrumentId, units, pos.Price, pos.Side}
}

func (pos *OpenPosition) String() string {
	return fmt.Sprintf("u: %d\np: %f\ns: %d", pos.Units, pos.Price, pos.Side)
}
