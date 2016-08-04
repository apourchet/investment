package main

import (
	"fmt"
	"io"
	"net"
	"os"
	"time"

	"golang.org/x/net/context"

	"google.golang.org/grpc"

	"github.com/apourchet/investment"
	pb "github.com/apourchet/investment/protos"
)

func SimpleCSVStrategy(broker pb.BrokerClient, stream pb.Broker_StreamQuotesClient) {
	for {
		q, err := stream.Recv()
		if err == io.EOF || q == nil {
			return
		}
		if int(q.Bid) == 5 {
			// Sell a bunch
			o := &pb.OrderCreation{}
			o.Instrument = "EURUSD"
			o.Type = pb.OrderType_MARKET
			o.Side = pb.OrderSide_SELL
			o.Units = 100
			broker.CreateOrder(context.Background(), o)
		}
		if int(q.Bid) == 1 {
			// Buy a bunch
			o := &pb.OrderCreation{}
			o.Instrument = "EURUSD"
			o.Type = pb.OrderType_MARKET
			o.Side = pb.OrderSide_BUY
			o.Units = 100
			broker.CreateOrder(context.Background(), o)
		}
	}
}

func startBroker(datafile string) {
	milliStep := 200
	broker := invt.NewDefaultBroker()
	go invt.Simulate(broker, datafile, milliStep)

	lis, err := net.Listen("tcp", ":8080")
	exitOnError(err)

	server := grpc.NewServer()
	pb.RegisterBrokerServer(server, broker)
	err = server.Serve(lis)
	exitOnError(err)
}

func startTrader() {
	conn, err := grpc.Dial(":8080", grpc.WithInsecure())
	exitOnError(err)
	defer conn.Close()

	broker := pb.NewBrokerClient(conn)
	iid := &pb.InstrumentID{"EURUSD"}
	stream, err := broker.StreamQuotes(context.Background(), iid)
	exitOnError(err)

	SimpleCSVStrategy(broker, stream)
}

func exitOnError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s", err.Error())
	}
}

func main() {
	go startBroker("./data/downup.csv")
	time.Sleep(time.Millisecond * 50)
	startTrader()
}
