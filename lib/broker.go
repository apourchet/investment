package invt

type Data interface{}

type Broker interface {
	GetQuote(string) Data
}

func StartBroker(b Broker) {

}
