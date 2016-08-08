package invt

import (
	"net"

	"google.golang.org/grpc"

	"fmt"

	bc "github.com/apourchet/investment/lib/broadcaster"
	pb "github.com/apourchet/investment/protos"
	"golang.org/x/net/context"
)

const (
	ONLY_INSTRUMENTID = "EURUSD"
)

type DefaultBroker struct {
	broadcaster *bc.Broadcaster
	lastquote   *pb.Quote
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
		return fmt.Errorf("Unsupported InstrumentID. Only support " + "EURUSD") // TODO
	}
	cb := make(chan interface{}, 10)
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

func (b *DefaultBroker) CreateOrder(context.Context, *pb.OrderCreationReq) (*pb.OrderCreationResp, error) {
	return nil, nil
}

func (b *DefaultBroker) OnQuote(q *pb.Quote) {
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

//func (b *DefaultBroker) GetPrices(ctx context.Context, il *pb.InstrumentIDList) (ls *pb.QuoteList, err error) {
//	for _, iid := range il.Value {
//		if iid.Id == ONLY_INSTRUMENTID {
//			ls.Value = append(ls.Value, b.lastquote)
//		}
//	}
//	return ls, err
//}
//func (b *DefaultBroker) CreateOrder(ctx context.Context, oc *pb.OrderCreation) (*pb.OrderCreationResponse, error) {
//	ocr := &pb.OrderCreationResponse{}
//	ocr.Time = time.Now().String()
//	ocr.Order = &pb.Order{}
//	ocr.Order.Instrument = oc.Instrument
//	ocr.Order.Side = oc.Side
//	ocr.Order.Units = oc.Units
//	ocr.Order.Price = b.lastquote.Bid
//	ocr.Order.Id = "1234"
//	ocr.Order.Type = pb.OrderType_MARKET
//
//	b.account.ProcessOrder(ocr.Order)
//	fmt.Printf("Balance: %f\n", b.account.Balance)
//	return ocr, nil
//}
