package invt

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

type BrokerHandler struct {
	Broker Broker
}

type BrokerClient struct {
	BrokerURL string
}

type Broker interface {
	GetQuote(QuoteRequest) Quote
}

func NewBrokerClient(brokerURL string) Broker {
	return BrokerClient{brokerURL}
}

// TODO error handling
func (bc BrokerClient) GetQuote(qr QuoteRequest) Quote {
	qrReader := strings.NewReader(qr.String())
	resp, err := http.Post("http://"+bc.BrokerURL+"/quote", "text", qrReader)
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()

	qbytes, _ := ioutil.ReadAll(resp.Body)
	qstring := string(qbytes)
	return ParseQuote(qstring)
}

// TODO
func (bh *BrokerHandler) quoteHandler(rw http.ResponseWriter, req *http.Request) {
	qr := parseQuoteRequest(req)
	quote := bh.Broker.GetQuote(qr)
	fmt.Fprintf(rw, "%s", quote)
}

func (bh *BrokerHandler) Start() error {
	http.HandleFunc("/quote", bh.quoteHandler)
	return http.ListenAndServe(":1026", nil)
}

// TODO
func parseQuoteRequest(req *http.Request) QuoteRequest {
	return QuoteRequest{}
}
