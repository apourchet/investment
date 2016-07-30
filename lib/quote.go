package invt

import (
	"encoding/json"
)

type Quote struct {
	BidPrice   float32                `json:"bidprice"`
	OfferPrice float32                `json:"offerprice"`
	TimeStamp  int64                  `json:"timestamp"`
	Error      string                 `json:"error"`
	Extra      map[string]interface{} `json:"extra,omitempty"`
}

type QuoteRequest struct {
	Currency string `json:"currency"`
	Lookback int    `json:"lookback"`
	Error    string `json:"error"`
}

func (q Quote) String() string {
	qBytes, _ := json.Marshal(q)
	return string(qBytes)
}

func (qr QuoteRequest) String() string {
	qrBytes, _ := json.Marshal(qr)
	return string(qrBytes)
}

func ParseQuote(qString string) Quote {
	q := Quote{}
	qBytes := []byte(qString)
	err := json.Unmarshal(qBytes, &q)
	if err != nil {
		q.Error = err.Error()
	}
	return q
}

func ParseQuoteRequest(qrString string) QuoteRequest {
	qr := QuoteRequest{}
	qrBytes := []byte(qrString)
	err := json.Unmarshal(qrBytes, &qr)
	if err != nil {
		qr.Error = err.Error()
	}
	return qr
}
