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
	ONLY_QUOTEID = "EURUSD"
)

type defaultBroker struct {
	broadcaster *bc.Broadcaster
	lastquote   *pb.Quote
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

		b.lastquote = &q
		b.broadcaster.Emit(q)
		time.Sleep(time.Millisecond * 10)
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
