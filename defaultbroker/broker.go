package main

import (
	"fmt"
	bc "github.com/apourchet/investment/lib/broadcaster"
	pb "github.com/apourchet/investment/protos"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"net"
	"time"
)

const (
	ONLY_QUOTEID = "EURUSD"
)

type defaultBroker struct {
	broadcaster *bc.Broadcaster
}

func (b *defaultBroker) GetQuote(ctx context.Context, qid *pb.QuoteID) (*pb.Quote, error) {
	q := &pb.Quote{}
	q.Name = qid.Name
	return q, nil
}

func (b *defaultBroker) StreamQuotes(qid *pb.QuoteID, stream pb.Broker_StreamQuotesServer) error {
	if qid.Name != ONLY_QUOTEID {
		return fmt.Errorf("We only support EURUSD as currency.")
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

func (b *defaultBroker) simulate() {
	for {
		q := pb.Quote{}
		q.Name = "EURUSD"
		q.Bid = 42
		q.Ask = 43

		b.broadcaster.Emit(q)
		time.Sleep(time.Second * 2)
	}
}

func main() {
	broker := defaultBroker{bc.NewBroadcaster()}
	go broker.simulate()

	lis, _ := net.Listen("tcp", ":8080")
	server := grpc.NewServer()
	pb.RegisterBrokerServer(server, &broker)
	server.Serve(lis)
}
