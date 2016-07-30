package invt

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
)

type Quote interface {
	String() string
}

type SimpleQuote string

type BrokerHandler struct {
	Broker Broker
}

type BrokerClient struct {
	BrokerURL string
}

type Broker interface {
	GetQuote(string, int) Quote
}

func (sq SimpleQuote) String() string {
	return string(sq)
}

func NewBrokerClient(brokerURL string) Broker {
	return BrokerClient{brokerURL}
}

// TODO
func (bc BrokerClient) GetQuote(qname string, lb int) Quote {
	fmt.Println("BrokerClient Getting Quote: " + bc.BrokerURL)
	vals := url.Values{"qname": {qname}, "lb": {strconv.Itoa(lb)}}
	resp, err := http.PostForm("http://"+bc.BrokerURL+"/quote", vals)
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()

	qbytes, _ := ioutil.ReadAll(resp.Body)
	qstring := string(qbytes)
	return SimpleQuote(qstring)
}

// TODO
func (bh *BrokerHandler) quoteHandler(rw http.ResponseWriter, req *http.Request) {
	qname, lb := parseQuote(req)
	quote := bh.Broker.GetQuote(qname, lb)
	fmt.Fprintf(rw, "%s", quote)
}

func (bh *BrokerHandler) Start() error {
	http.HandleFunc("/quote", bh.quoteHandler)
	return http.ListenAndServe(":1026", nil)
}

// TODO
func parseQuote(req *http.Request) (string, int) {
	return "EURUSD", 0
}
