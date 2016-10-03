package oanda

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"

	"golang.org/x/net/context"

	"github.com/apourchet/investment/protos"
	"google.golang.org/grpc"
)

type DeleteTrade struct {
	TransactionId uint64
	Price         float64
	Instrument    string
	Profit        float64
	Side          string
	Time          time.Time
}

type Order struct {
	Instrument string
	Units      int
	Side       string
	Type       string
	Expiry     time.Time
	Price      float64
	StopLoss   float64
	TakeProfit float64

	TrailingStop   float64
	TrailingAmount float64
	UpperBound     float64
	LowerBound     float64

	TradeTime time.Time
}

type OrderResponse struct {
	Instrument    string
	Time          time.Time
	Price         float64
	TradeOpened   Trade
	TradesClosed  []int
	TradesReduced map[int]int
}

type Trade struct {
	Id           uint64
	Instrument   string
	Units        int
	Side         string
	TakeProfit   float64
	StopLoss     float64
	TrailingStop float64
}

var _ protos.BrokerClient = &OandaBroker{}

type OandaBroker struct {
	token    Token
	acc      Account
	cl       *http.Client
	practice bool
}

func NewPracticeBroker() *OandaBroker {
	return &OandaBroker{
		token:    Token{os.Getenv("OANDA_TOKEN")},
		acc:      Account{},
		cl:       http.DefaultClient,
		practice: true,
	}
}

// TODO finish implementing the interface
func (o *OandaBroker) GetInstrumentList(ctx context.Context, in *protos.InstrumentListReq, opts ...grpc.CallOption) (*protos.InstrumentListResp, error) {
	return nil, nil
}

func (o *OandaBroker) GetLastCandle(ctx context.Context, in *protos.LastCandleReq, opts ...grpc.CallOption) (*protos.LastCandleResp, error) {
	return nil, nil
}

func (o *OandaBroker) StreamPrices(ctx context.Context, in *protos.StreamPricesReq, opts ...grpc.CallOption) (protos.Broker_StreamPricesClient, error) {
	return nil, nil
}

func (o *OandaBroker) StreamCandle(ctx context.Context, in *protos.StreamCandleReq, opts ...grpc.CallOption) (protos.Broker_StreamCandleClient, error) {
	return nil, nil
}

func (o *OandaBroker) GetOrders(ctx context.Context, in *protos.OrderListReq, opts ...grpc.CallOption) (*protos.OrderListResp, error) {
	return nil, nil
}

func (o *OandaBroker) CreateOrder(ctx context.Context, in *protos.OrderCreationReq, opts ...grpc.CallOption) (*protos.OrderCreationResp, error) {
	// TODO add stoploss takeprofit trailing stop to OrderCreationReq
	order, err := o.NewMarketOrder(int(in.Units), in.InstrumentId, in.Side, in.Type, 0.0, 0.0, 0.0)
	if err != nil {
		return nil, err
	}
	return &protos.OrderCreationResp{
		InstrumentId: order.Instrument,
		Time:         order.Time.String(),
		Price:        order.Price,
		Id:           strconv.FormatUint(order.TradeOpened.Id, 10),
	}, nil
}

func (o *OandaBroker) NewMarketOrder(units int, instrument, side, orderType string, takeProfit, stopLoss, trailingStop float64) (*OrderResponse, error) {
	data := url.Values{
		"type":       {orderType},
		"side":       {side},
		"units":      {strconv.Itoa(units)},
		"instrument": {instrument},
	}
	if takeProfit != 0 {
		data.Set("takeProfit", strconv.FormatFloat(takeProfit, 'f', -1, 64))
	}
	if stopLoss != 0 {
		data.Set("stopLoss", strconv.FormatFloat(stopLoss, 'f', -1, 64))
	}
	if trailingStop != 0 {
		data.Set("trailingStop", strconv.FormatFloat(trailingStop, 'f', -1, 64))
	}

	tr := &OrderResponse{}
	resp, err := o.doOandaRequest("POST", fmt.Sprintf("/v3/accounts/%d/orders", o.acc), data)
	if err != nil {
		return nil, err
	}
	err = handleOandaResponse(resp, tr)
	return tr, err
}

func (o *OandaBroker) GetTrade(tradeId string) (*Trade, error) {
	t := &Trade{}
	urlStr := fmt.Sprintf("/v3/accounts/%d/trades/%d", o.acc.AccountId, tradeId)
	resp, err := o.doOandaRequest("GET", urlStr, nil)
	if err != nil {
		return nil, err
	}
	err = handleOandaResponse(resp, t)
	return t, err
}

func (o *OandaBroker) CloseTrade(ctx context.Context, in *protos.CloseTradeReq, opts ...grpc.CallOption) (*protos.CloseTradeResp, error) {
	dt := &protos.CloseTradeResp{}
	urlStr := fmt.Sprintf("/v3/accounts/%d/trades/%d", o.acc.AccountId, in.Id)
	resp, err := o.doOandaRequest("DELETE", urlStr, nil)
	if err != nil {
		return nil, err
	}
	err = handleOandaResponse(resp, dt)
	return dt, err
}

// TODO incomplete
func (o *OandaBroker) ModifyTrade(tradeId string) (*Trade, error) {
	data := url.Values{}
	t := &Trade{}
	urlStr := fmt.Sprintf("/v3/accounts/%d/trades/%d", o.acc.AccountId, tradeId)
	resp, err := o.doOandaRequest("PATCH", urlStr, data)
	if err != nil {
		return nil, err
	}
	err = handleOandaResponse(resp, t)
	return t, err
}

// Accounts returns a list with all the known accounts
func (o *OandaBroker) listAccounts() ([]AccountShort, error) {
	accounts := &AccountListShort{make([]AccountShort, 0)}
	resp, err := o.doOandaRequest("GET", "/v3/accounts", nil)
	if err != nil {
		return nil, err
	}
	err = handleOandaResponse(resp, accounts)
	return accounts.Accounts, err
}

// Accounts returns a list with all the known accounts
func (o *OandaBroker) getAccount(id string) (*Account, error) {
	accounts := &AccountResp{}
	resp, err := o.doOandaRequest("GET", fmt.Sprintf("/v3/accounts/%s", id), nil)
	if err != nil {
		return nil, err
	}
	err = handleOandaResponse(resp, accounts)
	return &accounts.Account, err
}

func (o *OandaBroker) GetAccounts(ctx context.Context, in *protos.AccountListReq, opts ...grpc.CallOption) (*protos.AccountListResp, error) {
	accounts, err := o.listAccounts()
	if err != nil {
		return nil, err
	}
	resp := &protos.AccountListResp{
		Accounts: make([]*protos.Account, len(accounts)),
	}
	for i, acc := range accounts {
		// TODO get all info
		fullAccount, err := o.getAccount(acc.Id)
		fmt.Println(fullAccount, err)
		resp.Accounts[i] = &protos.Account{
			Id:       fullAccount.AccountId,
			Name:     fullAccount.Alias,
			Currency: fullAccount.Currency,
		}

	}
	return resp, nil
}

func (o *OandaBroker) GetAccountInfo(ctx context.Context, in *protos.AccountInfoReq, opts ...grpc.CallOption) (*protos.AccountInfoResp, error) {
	account, err := o.getAccount(in.AccountId)
	if err != nil {
		return nil, err
	}
	resp := &protos.AccountInfoResp{
		Info: &protos.AccountInfo{
			Id:       account.AccountId,
			Name:     account.Alias,
			Currency: account.Currency,
			//MarginRate:   account.MarginRate,
			//Balance: account.Balance,
			//UnrealizedPl: account.UnrealizedPl,
			//MarginUsed:   account.MarginUsed,
		},
	}
	return resp, nil
}
