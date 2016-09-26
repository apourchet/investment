package oanda

import (
	"time"
)

type AccountListShort struct {
	Accounts []AccountShort
}

type AccountResp struct {
	Account Account
}

type AccountShort struct {
	Id   string
	Tags []string
}

// Account represents an Oanda account.
// TODO positions and ordrs...
type Account struct {
	NAV                         string    `json:"NAV"`
	Alias                       string    `json:"alias"`
	Balance                     string    `json:"balance"`
	CreateByUserID              int       `json:"createdByUserID"`
	CreatedTime                 time.Time `json:"createdTime"`
	Currency                    string    `json:"currency"`
	Hedgeing                    bool      `json:"hedging"`
	AccountId                   string    `json:"id"`
	LastTransactionID           string    `json:"lastTransactionID"`
	MarginAvailable             string    `json:"marginAvailable"`
	MarginCallMarginUsed        string    `json:"marginCallMarginUsed"`
	MarginCallPercent           string    `json:"marginCallPercent"`
	MarginCloseoutMarginUseds   string    `json:"marginCloseoutMarginUseds"`
	MarginCloseoutNAV           string    `json:"marginCloseoutNAV"`
	MarginCloseoutPercent       string    `json:"marginCloseoutPercent"`
	MarginCloseoutPositionValue string    `json:"marginCloseoutPositionValue"`
	MarginCloseoutUnrealizedPL  string    `json:"marginCloseoutUnrealizedPL"`
	MarginRate                  string    `json:"marginRate"`
	MarginUsed                  string    `json:"marginUsed"`
	OpenPositionCount           int       `json:"openPositionCount"`
	OpenTradeCount              int       `json:"openTradeCount"`
	orders                      []string
	pendingOrderCount           int
	pl                          string
	positionValue               string
	positions                   []string
	resettablePL                string
	trades                      []string
	UnrealizedPL                string `json:"unrealizePL"`
	withdrawalLimit             string
}
