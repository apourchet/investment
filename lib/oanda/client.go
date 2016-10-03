package oanda

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"
)

func NewDefaultHttpClient() *http.Client {
	return &http.Client{
		Transport: &http.Transport{
			Dial: (&net.Dialer{
				Timeout:   30 * time.Second,
				KeepAlive: 30 * time.Second,
			}).Dial,
			TLSHandshakeTimeout: 20 * time.Second,
		},
	}
}

type OandaError struct {
	Code     int    `json:"code"`
	Message  string `json:"message"`
	MoreInfo string `json:"moreInfo"`
}

func (o *OandaError) Error() string {
	return o.Message
}

func CloseResponse(rc io.ReadCloser) {
	io.Copy(ioutil.Discard, rc)
	rc.Close()
}

type Token struct {
	token string
}

func (t *Token) AddTokenToRequest(r *http.Request) {
	r.Header.Set("Authorization", fmt.Sprintf("Bearer %s", t.token))
}

func (o *OandaBroker) doOandaRequest(method, urlStr string, data url.Values) (*http.Response, error) {
	var body io.Reader
	if len(data) > 0 {
		body = strings.NewReader(data.Encode())
	}
	req, err := http.NewRequest(method, urlStr, body)
	if err != nil {
		return nil, err
	}
	setOandaUrl(req, o.practice)
	o.token.AddTokenToRequest(req)
	return o.cl.Do(req)
}

func setOandaUrl(req *http.Request, practice bool) {
	req.URL.Scheme = "https"
	if req.URL.Host == "" {
		if practice {
			req.URL.Host = fmt.Sprintf("api-%s.oanda.com", "fxpractice")
		} else {
			req.URL.Host = fmt.Sprintf("api-%s.oanda.com", "fxtrade")
		}
	}
}

func handleOandaResponse(resp *http.Response, value interface{}) error {
	defer CloseResponse(resp.Body)

	dec := json.NewDecoder(resp.Body)
	if resp.StatusCode < 400 {
		return dec.Decode(value)
	}

	oandaError := OandaError{}
	err := dec.Decode(&oandaError)
	if err != nil {
		return err
	}
	return &oandaError
}
