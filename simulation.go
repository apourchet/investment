package invt

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"time"
)

type Simulatable interface {
	ParseQuote([]string) *Quote
	OnQuote(*Quote)
	OnEnd()
}

func SimulateDataStream(s Simulatable, datafile string, milliStep int) {
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
			s.OnEnd()
			break
		}
		q := s.ParseQuote(record)

		s.OnQuote(q)
		time.Sleep(time.Millisecond * time.Duration(milliStep))
	}
}
