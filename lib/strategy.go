package invt

type Strategy struct {
	QuoteRequest chan string
	BuyOrder     chan string
	SellOrder    chan string
}
