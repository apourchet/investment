package invt

func Buy(a *Account, instrumentId string, units int32, price float64) {
	do(a, instrumentId, units, price, "buy")
}

func Sell(a *Account, instrumentId string, units int32, price float64) {
	do(a, instrumentId, units, price, "sell")
}

func do(a *Account, instrumentId string, units int32, price float64, side string) {
	pos := &OpenPosition{instrumentId, units, price, side}
	if other, ok := a.OpenPositions[pos.InstrumentId]; ok {
		mergePositions(a, pos, other)
	} else {
		openNewPosition(a, pos)
	}
}

func closePosition(a *Account, pos *OpenPosition, price float64) {
	a.Balance += pos.Value() // Gain value of position
	if pos.Side == "buy" {
		a.Balance += pos.FloatUnits() * (price - pos.Price) // Gain delta
	} else {
		a.Balance += pos.FloatUnits() * (pos.Price - price) // Gain delta
	}
}

func mergePositions(a *Account, from, to *OpenPosition) {
	if from.Side == to.Side {
		a.Balance -= from.Value()
		totalUnits := from.Units + to.Units
		totalValue := from.Value() + to.Value()
		avgPrice := totalValue - float64(totalUnits)

		to.Price = avgPrice
		to.Units = totalUnits
	} else {
		if from.Units == to.Units {
			closePosition(a, to, from.Price)
			delete(a.OpenPositions, to.InstrumentId)
		} else if to.Units > from.Units {
			toclose := to.SplitPosition(from.Units)
			closePosition(a, toclose, from.Price)
		} else if from.Units > to.Units {
			closePosition(a, to, from.Price)
			delete(a.OpenPositions, to.InstrumentId)
			from.Units -= to.Units
			openNewPosition(a, from)
		}
	}
}

func openNewPosition(a *Account, pos *OpenPosition) {
	a.Balance -= pos.Value()
	a.OpenPositions[pos.InstrumentId] = pos
}
