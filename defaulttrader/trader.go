package main

import (
	"fmt"
	pb "github.com/apourchet/investment/broker"
	. "github.com/apourchet/investment/lib"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"time"
)

func mystrat(tr *Trader) {
	broker := tr.Broker
	for {
		q := broker.GetQuote(QuoteRequest{"EURUSD", 0, ""})
		fmt.Println("QuoteResponse: " + q.String())
		time.Sleep(time.Second * 5)
	}
}

func main() {
	conn, _ := grpc.Dial(":8080", grpc.WithInsecure())
	defer conn.Close()

	c := pb.NewBrokerClient(conn)

	q, _ := c.GetQuote(context.Background(), &pb.QuoteID{"EURUSD"})
	fmt.Println("Quote: " + q.Name)
}
