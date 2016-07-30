package invt

type Trader struct {
	Broker   Broker
	Strategy Strategy
	Margin   int
}
