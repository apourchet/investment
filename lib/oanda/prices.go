package oanda

import (
	"fmt"
	"golang.org/x/net/context"
	"net/url"
	"time"

	"github.com/apourchet/investment/protos"
	"google.golang.org/grpc"
)

type Price struct {
	Ti     time.Time
	Bid    float64
	Ask    float64
	Status string
}

func (o *OandaBroker) GetPrices(ctx context.Context, in *protos.PriceListReq, opts ...grpc.CallOption) (*protos.PriceListResp, error) {
	vals := url.Values{}
	for i, instrumentId := range in.InstrumentId {
		if i == 0 {
			vals.Set("instruments", instrumentId)
		} else {
			vals.Add("instruments", instrumentId)
		}
	}
	prices := map[string][]*protos.Quote{}
	resp, err := o.doOandaRequest("GET", fmt.Sprintf("/v1/prices?%s", vals.Encode()), nil)
	if err != nil {
		return nil, err
	}
	err = handleOandaResponse(resp, &prices)
	return &protos.PriceListResp{prices["prices"]}, err

}
