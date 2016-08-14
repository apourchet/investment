package invt

import (
	"encoding/csv"
	"io"
	"os"
	"time"

	"log"

	pb "github.com/apourchet/investment/protos"
)

type Simulator struct {
	Format     DataFormat
	Filename   string
	StepMillis int
}

type DataFormat int

type Strat func(pb.BrokerClient, pb.Broker_StreamPricesClient)

type Simulatable interface {
	Start() error
	OnData([]string, DataFormat)
	OnEnd()
}

const (
	DATAFORMAT_QUOTE  = DataFormat(iota)
	DATAFORMAT_CANDLE = DataFormat(iota)
)

func NewSimulator(format DataFormat, filename string, stepmillis int) *Simulator {
	return &Simulator{format, filename, stepmillis}
}

func (s *Simulator) SimulateDataStream(sim Simulatable) error {
	go func() {
		in, err := os.Open(s.Filename)
		if err != nil {
			log.Fatalln("Could not open data file:", err)
		}

		log.Println("Simulating: " + s.Filename)
		reader := csv.NewReader(in)
		for {
			record, err := reader.Read()
			if err == io.EOF {
				log.Println("Simulation Ended")
				sim.OnEnd()
				break
			}
			sim.OnData(record, s.Format)
			time.Sleep(time.Millisecond * time.Duration(s.StepMillis))
		}
	}()
	return sim.Start()
}
