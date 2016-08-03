package main

import (
	pb "github.com/apourchet/investment/broker"
	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
	"net"
)

type DefaultBroker struct {
}

func (b *DefaultBroker) GetQuote(ctx context.Context, qid *pb.QuoteID) (*pb.Quote, error) {
	q := &pb.Quote{}
	q.Name = "EURUSD"
	return q, nil
}

func (b *DefaultBroker) StreamQuotes(qid *pb.QuoteID, stream pb.Broker_StreamQuotesServer) error {
	return nil
}

func main() {
	lis, _ := net.Listen("tcp", ":8080")

	server := grpc.NewServer()
	pb.RegisterBrokerServer(server, &DefaultBroker{})
	server.Serve(lis)
}
