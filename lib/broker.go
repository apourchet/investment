package invt

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
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

// TODO
func (bc BrokerClient) GetQuote(qr QuoteRequest) Quote {
	fmt.Println("BrokerClient Getting Quote: " + bc.BrokerURL)
	vals := url.Values{"qname": {qr.QuoteName}, "lb": {strconv.Itoa(qr.Lookback)}}
	resp, err := http.PostForm("http://"+bc.BrokerURL+"/quote", vals)
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
	qr := parseQuote(req)
	quote := bh.Broker.GetQuote(qr)
	fmt.Fprintf(rw, "%s", quote)
}

func (bh *BrokerHandler) Start() error {
	http.HandleFunc("/quote", bh.quoteHandler)
	return http.ListenAndServe(":1026", nil)
}

// TODO
func parseQuote(req *http.Request) QuoteRequest {
	return QuoteRequest{}
}
