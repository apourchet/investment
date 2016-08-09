package invt

import (
	"net"

	"google.golang.org/grpc"

	"fmt"

	"strconv"

	bc "github.com/apourchet/investment/lib/broadcaster"
	pb "github.com/apourchet/investment/protos"
	"golang.org/x/net/context"
)

const (
	ONLY_INSTRUMENTID = "EURUSD"
)

type DefaultBroker struct {
	broadcaster *bc.Broadcaster
	lastquote   *Quote
	account     *Account
}

func NewDefaultBroker() *DefaultBroker {
	return &DefaultBroker{bc.NewBroadcaster(), nil, NewAccount(10000)}
}

func (b *DefaultBroker) GetInstrumentList(context.Context, *pb.InstrumentListReq) (*pb.InstrumentListResp, error) {
	return nil, nil
}

func (b *DefaultBroker) GetPrices(context.Context, *pb.PriceListReq) (*pb.PriceListResp, error) {
	return nil, nil
}

func (b *DefaultBroker) StreamPrices(req *pb.StreamPricesReq, stream pb.Broker_StreamPricesServer) error {
	if req.InstrumentId != ONLY_INSTRUMENTID {
		return fmt.Errorf("Unsupported InstrumentID. Only support " + "EURUSD")
	}
	cb := make(chan interface{})
	rid := b.broadcaster.Register(cb)
	for qdata := range cb {
		if qdata == nil {
			stream.Send(nil)
			b.broadcaster.Deregister(rid)
			return nil
		} else {
			q := qdata.(*pb.Quote)
			err := stream.Send(q)
			if err != nil {
				b.broadcaster.Deregister(rid)
				return err
			}
		}
	}
	return nil
}

func (b *DefaultBroker) GetAccounts(context.Context, *pb.AccountListReq) (*pb.AccountListResp, error) {
	return nil, nil
}

func (b *DefaultBroker) GetAccountInfo(context.Context, *pb.AccountInfoReq) (*pb.AccountInfoResp, error) {
	return nil, nil
}

func (b *DefaultBroker) GetOrders(context.Context, *pb.OrderListReq) (*pb.OrderListResp, error) {
	return nil, nil
}

func (b *DefaultBroker) CreateOrder(ctx context.Context, req *pb.OrderCreationReq) (*pb.OrderCreationResp, error) {
	if req.Type == "market" {
		TradeQuote(b.account, b.lastquote, req.Units, ParseSide(req.Side))
	}
	return nil, nil
}

func (b *DefaultBroker) OnQuote(q *Quote) {
	b.lastquote = q
	b.broadcaster.Emit(q)
}

func (b *DefaultBroker) OnEnd() {
	b.broadcaster.Emit(nil)
}

func (b *DefaultBroker) Start() error {
	lis, err := net.Listen("tcp", ":8080")
	if err != nil {
		return err
	}

	server := grpc.NewServer()
	pb.RegisterBrokerServer(server, b)

	server.Serve(lis)
	return nil
}

func (b *DefaultBroker) ParseQuote(record []string) *Quote {
	q := &Quote{}
	q.InstrumentId = "EURUSD"
	v, err := strconv.ParseFloat(record[2], 64)
	q.Bid = v
	if err != nil {
		return nil
	}

	q.Ask, err = strconv.ParseFloat(record[4], 64)
	if err != nil {
		return nil
	}
	// TODO
	// q.Time = date.ParseDate(record[0])
	return q
}
