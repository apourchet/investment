package invt

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strconv"
	"time"

	pb "github.com/apourchet/investment/protos"
)

type Simulatable interface {
	pb.BrokerServer
	OnQuote(*Quote)
	OnEnd()
}

func SimulateDataStream(b Simulatable, datafile string, milliStep int) {
	in, err := os.Open(datafile)
	if err != nil {
		fmt.Println("Could not open data file: " + err.Error())
		os.Exit(1)
	}

	fmt.Println("Simulating: " + datafile)
	reader := csv.NewReader(in)
	for {
		record, err := reader.Read()
		if err == io.EOF {
			fmt.Println("Simulation Ended")
			b.OnEnd()
			break
		}
		q := &Quote{}
		q.InstrumentId = "EURUSD"
		q.Bid, err = strconv.ParseFloat(record[2], 64)
		q.Ask, err = strconv.ParseFloat(record[4], 64)
		// TODO
		// q.Time = date.ParseDate(record[0])

		b.OnQuote(q)
		time.Sleep(time.Millisecond * time.Duration(milliStep))
	}
}
