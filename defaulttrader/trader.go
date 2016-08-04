package main

import (
	"fmt"
	pb "github.com/apourchet/investment/protos"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"io"
)

func main() {
	conn, _ := grpc.Dial(":8080", grpc.WithInsecure())
	defer conn.Close()

	broker := pb.NewBrokerClient(conn)
	iid := &pb.InstrumentID{"EURUSD"}
	stream, _ := broker.StreamQuotes(context.Background(), iid)

	for {
		q, err := stream.Recv()
		if err == io.EOF || q == nil {
			// Done reading
			return
		}
		fmt.Println(q)
	}

	// for {
	// 	q, _ := broker.GetQuote(context.Background(), &pb.QuoteID{"EURUSD"})
	// 	fmt.Println("Quote: " + q.String())
	// 	time.Sleep(time.Second * 2)
	// }
}
