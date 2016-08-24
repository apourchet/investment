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
	Format       DataFormat
	Filename     string
	StepMillis   int
	Simulatables []Simulatable
}

type DataFormat int

type Strat func(pb.BrokerClient, pb.Broker_StreamPricesClient)

type Simulatable interface {
	Start() error // Assumed to be blocking
	OnData([]string, DataFormat)
	OnEnd()
}

const (
	DATAFORMAT_QUOTE  = DataFormat(iota)
	DATAFORMAT_CANDLE = DataFormat(iota)
)

func NewSimulator(format DataFormat, filename string, stepmillis int) *Simulator {
	s := &Simulator{}
	s.Format = format
	s.Filename = filename
	s.StepMillis = stepmillis
	s.Simulatables = make([]Simulatable, 0)
	return s
}

func (s *Simulator) AddSimulatable(sim Simulatable) {
	s.Simulatables = append(s.Simulatables, sim)
}

func (s *Simulator) StartSimulation() error {
	for _, sim := range s.Simulatables {
		go sim.Start()
	}
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
				for _, sim := range s.Simulatables {
					sim.OnEnd()
				}
				break
			}
			for _, sim := range s.Simulatables {
				sim.OnData(record, s.Format)
			}
			time.Sleep(time.Millisecond * time.Duration(s.StepMillis))
		}
	}()
	return nil
}

// Deprecated
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
