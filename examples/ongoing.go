package main

import (
	"fmt"
	"io"
	"os"

	"golang.org/x/net/context"

	"time"

	"log"

	"github.com/apourchet/investment"
	tl "github.com/apourchet/investment/lib/tradelogger"
	pb "github.com/apourchet/investment/protos"
	influx "github.com/influxdata/influxdb/client/v2"
)

var (
	ixClient influx.Client
)

func quickOrder(units int32, side string) *pb.OrderCreationReq {
	o := &pb.OrderCreationReq{}
	o.InstrumentId = "EURUSD"
	o.Type = invt.TYPE_MARKET
	o.Side = side
	o.Units = units
	return o
}

func getStream(broker pb.BrokerClient) pb.Broker_StreamCandleClient {
	req := &pb.StreamCandleReq{&pb.AuthToken{}, "EURUSD"}
	stream, err := broker.StreamCandle(context.Background(), req)
	if err != nil {
		fmt.Println("Error: " + err.Error())
		os.Exit(1)
	}
	return stream
}

func getInfluxClient() influx.Client {
	if ixClient != nil {
		return ixClient
	}
	c, err := influx.NewHTTPClient(influx.HTTPConfig{
		Addr:     "http://localhost:8086",
		Username: "investment",
		Password: "password",
	})
	if err != nil {
		log.Fatalln("ERROR: ", err)
	}
	ixClient = c
	return c
}

func createBatch() influx.BatchPoints {
	// Create a new point batch
	bp, err := influx.NewBatchPoints(influx.BatchPointsConfig{
		Database:  "testdb",
		Precision: "s",
	})

	if err != nil {
		log.Fatalln("Error: ", err)
	}
	return bp

}

func writePoint(candle *pb.Candle) {
	c := getInfluxClient()
	bp := createBatch()
	tags := map[string]string{"trader": "12345"}
	fields := map[string]interface{}{
		"low":  candle.Low,
		"high": candle.High,
	}
	pt, err := influx.NewPoint("candle_close", tags, fields, time.Now())

	if err != nil {
		log.Fatalln("Error: ", err)
	}

	bp.AddPoint(pt)

	// Write the batch
	c.Write(bp)

}

// Create a point and add to batch
func mine(def *invt.DefaultBroker) {
	fmt.Println("Trader started")
	broker := def.GetClient()
	stream := getStream(broker)

	steps := 0 // ema5 < ema30
	for {
		c, err := stream.Recv()
		if err == io.EOF || c == nil {
			fmt.Println("Candle stream has ended.")
			return
		}

		if steps%20 == 0 {
			req := &pb.AccountInfoReq{}
			resp, _ := broker.GetAccountInfo(context.Background(), req)
			fmt.Println(resp.Info.MarginAvail)
			writePoint(c)
		}

		if c.Close-c.Low > (c.High-c.Close)*4 {
			o := quickOrder(3000, invt.StringOfSide(invt.SIDE_BUY))
			broker.CreateOrder(context.Background(), o)
		} else if (c.Close-c.Low)*4 < c.High-c.Close {
			o := quickOrder(3000, invt.StringOfSide(invt.SIDE_SELL))
			broker.CreateOrder(context.Background(), o)
		}
		steps++
	}
}

func main() {
	go tl.StartServer(1026, "logs/")
	time.Sleep(time.Millisecond * 50)

	datafile := "examples/data/largish.csv"
	if len(os.Args) >= 2 {
		datafile = os.Args[1]
	}

	simulator := invt.NewSimulator(invt.DATAFORMAT_CANDLE, datafile, 10)
	broker := invt.NewDefaultBroker(1027)
	go simulator.SimulateDataStream(broker)

	time.Sleep(time.Millisecond * 50)
	mine(broker)
}
