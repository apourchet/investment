package invt

import (
	"fmt"
	"net/http"
)

type Quote interface {
	String() string
}

type BrokerHandler struct {
	Broker Broker
}

type Broker interface {
	Start()
	GetQuote(string, int64) Quote
}

func (bh *BrokerHandler) quoteHandler(rw http.ResponseWriter, req *http.Request) {
	qname, lb := parseQuote(req)
	quote := bh.Broker.GetQuote(qname, lb)
	fmt.Fprintf(rw, "%s", quote)
}

// func (bh *BrokerHandler) quoteStreamHandler(rw http.ResponseWriter, req *http.Request) {
// 	f, cn := checkStreamable(rw)
// 	if cn == nil {
// 		fmt.Fprintf(rw, `{"success":0}`)
// 		return
// 	}
//
// 	setStreamHeaders(rw)
//
// 	quote, _ := parseQuote(req)
// 	ticker := time.NewTicker(time.Second * 1)
// 	for {
// 		select {
// 		case <-cn.CloseNotify():
// 			fmt.Fprintf(rw, `{"success":0}`)
// 			return
// 		case <-ticker.C:
// 			fmt.Fprintf(rw, "data:%s\n\n", bh.Broker.GetQuote(quote, 0))
// 			f.Flush()
// 		}
// 	}
// }

func (bh *BrokerHandler) Start() error {
	http.HandleFunc("/quote", bh.quoteHandler)
	// http.HandleFunc("/quote_stream", bh.quoteStreamHandler)
	bh.Broker.Start()
	return http.ListenAndServe(":1026", nil)
}

func parseQuote(req *http.Request) (string, int64) {
	return "EURUSD", 0
}

func setStreamHeaders(rw http.ResponseWriter) {
	rw.Header().Set("Content-Type", "text/event-stream")
	rw.Header().Set("Cache-Control", "no-cache")
	rw.Header().Set("Connection", "keep-alive")
	rw.Header().Set("Access-Control-Allow-Origin", "*")
}

func checkStreamable(rw http.ResponseWriter) (http.Flusher, http.CloseNotifier) {
	f, ok := rw.(http.Flusher)
	if !ok {
		http.Error(rw, "cannot stream", http.StatusInternalServerError)
		return nil, nil
	}

	cn, ok := rw.(http.CloseNotifier)
	if !ok {
		http.Error(rw, "cannot stream", http.StatusInternalServerError)
		return nil, nil
	}
	return f, cn
}
