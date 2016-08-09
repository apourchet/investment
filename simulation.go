package invt

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"time"

	"golang.org/x/net/context"
	"google.golang.org/grpc"

	pb "github.com/apourchet/investment/protos"
)

type Strat func(broker pb.BrokerClient, stream pb.Broker_StreamPricesClient)

type Simulatable interface {
	Start() error
	ParseQuote([]string) *Quote
	OnQuote(*Quote)
	OnEnd()
}

func SimulateDataStream(s Simulatable, datafile string, milliStep int) error {
	in, err := os.Open(datafile)
	if err != nil {
		fmt.Println("Could not open data file: " + err.Error())
		return err
	}

	fmt.Println("Simulating: " + datafile)
	reader := csv.NewReader(in)
	for {
		record, err := reader.Read()
		if err == io.EOF {
			fmt.Println("Simulation Ended")
			s.OnEnd()
			break
		}
		q := s.ParseQuote(record)

		s.OnQuote(q)
		time.Sleep(time.Millisecond * time.Duration(milliStep))
	}
	return nil
}

func SimulateTradingScenario(s Simulatable, strat Strat, datafile string) error {
	go s.Start()
	time.Sleep(time.Millisecond * 50)

	conn, err := grpc.Dial(":8080", grpc.WithInsecure())
	if err != nil {
		return err
	}
	defer conn.Close()

	broker := pb.NewBrokerClient(conn)

	req := &pb.StreamPricesReq{&pb.AuthToken{}, "EURUSD"}
	stream, err := broker.StreamPrices(context.Background(), req)
	if err != nil {
		return err
	}

	go strat(broker, stream)

	return SimulateDataStream(s, datafile, 1)
}
