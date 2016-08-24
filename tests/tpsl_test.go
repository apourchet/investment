package test

import (
	. "github.com/apourchet/investment"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestBasicTP(t *testing.T) {
	acc := NewAccount(0)
	tm := NewTradingManager(acc)
	Buy(acc, "EURUSD", 10, 1.)
	tm.OnQuote(&Quote{"EURUSD", time.Now(), 1., 1.})
	tm.TakeProfit("EURUSD", 2., 1.)
	tm.OnQuote(&Quote{"EURUSD", time.Now(), 1.5, 1.})
	tm.OnQuote(&Quote{"EURUSD", time.Now(), 1.9, 2.2}) // STILL DONT SELL
	tm.OnQuote(&Quote{"EURUSD", time.Now(), 2.1, 2.2})
	assert.InEpsilon(t, acc.Balance, 11., 0.00001)
}

func TestBasicSL(t *testing.T) {
	acc := NewAccount(0)
	tm := NewTradingManager(acc)
	Buy(acc, "EURUSD", 10, 1.)
	tm.OnQuote(&Quote{"EURUSD", time.Now(), 1., 1.})
	tm.StopLoss("EURUSD", 0.5, 1.)
	tm.OnQuote(&Quote{"EURUSD", time.Now(), 1.5, 1.})
	tm.OnQuote(&Quote{"EURUSD", time.Now(), 1.9, 2.1})
	tm.OnQuote(&Quote{"EURUSD", time.Now(), 0.6, 0.6})
	tm.OnQuote(&Quote{"EURUSD", time.Now(), 0.5, 0.6})
	assert.InEpsilon(t, acc.Balance, -5., 0.00001)
}

func TestBasicSLSell(t *testing.T) {
	acc := NewAccount(0)
	tm := NewTradingManager(acc)
	Sell(acc, "EURUSD", 10, 1.)
	tm.OnQuote(&Quote{"EURUSD", time.Now(), 1., 1.})
	tm.StopLoss("EURUSD", 2., 1.)
	tm.OnQuote(&Quote{"EURUSD", time.Now(), 1.5, 1.})
	tm.OnQuote(&Quote{"EURUSD", time.Now(), 1.9, 2.2})
	assert.InEpsilon(t, acc.Balance, -12., 0.00001)
}

func TestBasicTPSell(t *testing.T) {
	acc := NewAccount(0)
	tm := NewTradingManager(acc)
	Sell(acc, "EURUSD", 10, 1.)
	tm.OnQuote(&Quote{"EURUSD", time.Now(), 1., 1.})
	tm.TakeProfit("EURUSD", 0.5, 1.)
	tm.OnQuote(&Quote{"EURUSD", time.Now(), 1.5, 1.})
	tm.OnQuote(&Quote{"EURUSD", time.Now(), 1.9, 2.1})
	tm.OnQuote(&Quote{"EURUSD", time.Now(), 0.6, 0.6})
	tm.OnQuote(&Quote{"EURUSD", time.Now(), 0.5, 0.6})
	tm.OnQuote(&Quote{"EURUSD", time.Now(), 0.5, 0.5})
	assert.InEpsilon(t, acc.Balance, 5., 0.0001)
}
