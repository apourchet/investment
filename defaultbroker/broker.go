package main

import (
	"encoding/csv"
	"fmt"
	bc "github.com/apourchet/investment/lib/broadcaster"
	pb "github.com/apourchet/investment/protos"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"io"
	"net"
	"os"
	"strconv"
	"time"
)

const (
	ONLY_INSTRUMENTID = "EURUSD"
)

type defaultBroker struct {
	broadcaster *bc.Broadcaster
	lastquote   *pb.Quote
}

func (b *defaultBroker) GetInstrumentList(ctx context.Context, token *pb.AuthToken) (ls *pb.InstrumentList, err error) {
	ins := pb.Instrument{}
	ins.Name = ONLY_INSTRUMENTID
	ins.DisplayName = ins.Name
	ins.Pip = "0.0001"
	ins.MaxTradeUnits = 10000
	ls.Value = append(ls.Value, &ins)
	return ls, nil
}

func (b *defaultBroker) GetPrices(ctx context.Context, il *pb.InstrumentIDList) (ls *pb.QuoteList, err error) {
	for _, iid := range il.Value {
		if iid.Value == ONLY_INSTRUMENTID {
			ls.Value = append(ls.Value, b.lastquote)
		}
	}
	return ls, err
}

func (b *defaultBroker) StreamQuotes(iid *pb.InstrumentID, stream pb.Broker_StreamQuotesServer) error {
	if iid.Value != ONLY_INSTRUMENTID {
		return fmt.Errorf("Unsupported InstrumentID. Only support " + ONLY_INSTRUMENTID)
	}
	cb := make(chan interface{}, 10)
	rid := b.broadcaster.Register(cb)
	for qdata := range cb {
		q := qdata.(pb.Quote)
		err := stream.Send(&q)
		if err != nil {
			b.broadcaster.Deregister(rid)
			return err
		}
	}
	return nil
}

func (b *defaultBroker) GetAccounts(ctx context.Context, token *pb.AuthToken) (*pb.AccountList, error) {
	return nil, fmt.Errorf("Not Implemented")
}

func (b *defaultBroker) GetOrders(ctx context.Context, accid *pb.AccountID) (*pb.OrderList, error) {
	return nil, fmt.Errorf("Not Implemented")
}

func (b *defaultBroker) CreateOrder(ctx context.Context, oc *pb.OrderCreation) (*pb.OrderCreationResponse, error) {
	return nil, fmt.Errorf("Not Implemented")
}

func (b *defaultBroker) ChangeOrder(ctx context.Context, oc *pb.OrderChange) (*pb.Order, error) {
	return nil, fmt.Errorf("Not Implemented")
}

func (b *defaultBroker) simulate(datafile string) {
	in, err := os.Open(datafile)
	if err != nil {
		fmt.Println("Could not open data file.")
		os.Exit(1)
	}

	reader := csv.NewReader(in)
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		fmt.Println(record)
		q := pb.Quote{}
		q.Name = "EURUSD"
		q.Bid, err = strconv.ParseFloat(record[2], 64)
		q.Ask, err = strconv.ParseFloat(record[4], 64)

		// TODO
		// q.Time = date.ParseDate(record[0])

		b.lastquote = &q
		b.broadcaster.Emit(q)
		fmt.Println(q)
		time.Sleep(time.Millisecond * 1000)
	}
}

func main() {
	broker := defaultBroker{bc.NewBroadcaster(), nil}
	go broker.simulate("data/DAT_MT_EURUSD_M1_2006.csv")

	lis, _ := net.Listen("tcp", ":8080")
	server := grpc.NewServer()
	pb.RegisterBrokerServer(server, &broker)
	server.Serve(lis)
}
