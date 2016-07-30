package invt

type Quote struct {
	BidPrice   float32                `json:"bidprice"`
	OfferPrice float32                `json:"offerprice"`
	TimeStamp  int64                  `json:"timestamp"`
	Error      bool                   `json:"error"`
	Extra      map[string]interface{} `json:"extra",omitempty`
}

type QuoteRequest struct {
	QuoteName string
	Lookback  int
}

// TODO
func (q Quote) String() string {
	return "42"
}

// TODO
func ParseQuote(qstring string) Quote {
	return Quote{}
}
