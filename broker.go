package invt

import (
	"fmt"
	"time"

	bc "github.com/apourchet/investment/lib/broadcaster"
	pb "github.com/apourchet/investment/protos"
	"golang.org/x/net/context"
)

const (
	ONLY_INSTRUMENTID = pb.InstrumentID_EURUSD
)

type DefaultBroker struct {
	broadcaster *bc.Broadcaster
	lastquote   *pb.Quote
	account     *Account
}

func NewDefaultBroker() *DefaultBroker {
	return &DefaultBroker{bc.NewBroadcaster(), nil, CreateNewAccount()}
}

func (b *DefaultBroker) GetInstrumentList(ctx context.Context, token *pb.AuthToken) (ls *pb.InstrumentList, err error) {
	ins := pb.Instrument{}
	ins.Name = ONLY_INSTRUMENTID
	ins.DisplayName = "EURUSD" // TODO map the pb.InstrumentIDs to displaynames
	ins.Pip = "0.0001"
	ins.MaxTradeUnits = 10000
	ls.Value = append(ls.Value, &ins)
	return ls, nil
}

func (b *DefaultBroker) GetPrices(ctx context.Context, il *pb.InstrumentIDList) (ls *pb.QuoteList, err error) {
	for _, iid := range il.Value {
		if iid.Id == ONLY_INSTRUMENTID {
			ls.Value = append(ls.Value, b.lastquote)
		}
	}
	return ls, err
}

func (b *DefaultBroker) StreamQuotes(iid *pb.InstrumentID, stream pb.Broker_StreamQuotesServer) error {
	if iid.Id != ONLY_INSTRUMENTID {
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

// TODO implement
func (b *DefaultBroker) GetAccounts(ctx context.Context, token *pb.AuthToken) (*pb.AccountList, error) {
	return nil, fmt.Errorf("Not Implemented")
}

// TODO implement
func (b *DefaultBroker) GetOrders(ctx context.Context, accid *pb.AccountID) (*pb.OrderList, error) {
	return nil, fmt.Errorf("Not Implemented")
}

func (b *DefaultBroker) CreateOrder(ctx context.Context, oc *pb.OrderCreation) (*pb.OrderCreationResponse, error) {
	ocr := &pb.OrderCreationResponse{}
	ocr.Time = time.Now().String()
	ocr.Order = &pb.Order{}
	ocr.Order.Instrument = oc.Instrument
	ocr.Order.Side = oc.Side
	ocr.Order.Units = oc.Units
	ocr.Order.Price = b.lastquote.Bid
	ocr.Order.Id = "1234"
	ocr.Order.Type = pb.OrderType_MARKET

	fmt.Println(b.account.Balance)
	b.account.ProcessOrder(ocr.Order)
	fmt.Println(b.account.Balance)
	return ocr, nil
}

func (b *DefaultBroker) ChangeOrder(ctx context.Context, oc *pb.OrderChange) (*pb.Order, error) {
	return nil, fmt.Errorf("Not Implemented")
}

func (b *DefaultBroker) OnQuote(q *pb.Quote) {
	fmt.Println(q)
	b.lastquote = q
	b.broadcaster.Emit(q)
}

func (b *DefaultBroker) OnEnd() {
	b.broadcaster.Emit(nil)
}
