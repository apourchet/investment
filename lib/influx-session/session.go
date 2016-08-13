package ix_session

import (
	"crypto/rand"
	"log"
	"sync"

	"encoding/hex"

	"time"

	structs "github.com/fatih/structs"
	ix "github.com/influxdata/influxdb/client/v2"
)

type Session struct {
	Id       string
	Address  string
	Username string
	Password string
	Database string
	client   ix.Client
}

const (
	RAND_ID_SIZE       = 10
	DEFAULT_ADDRESS    = "http://localhost:8086"
	DEFAULT_BATCH_SIZE = 1   // TODO
	DEFAULT_PRECISION  = "s" // TODO
)

var (
	once sync.Once
)

// Default address is localhost:8086
func NewSession(address string, username string, password string, database string) *Session {
	s := &Session{}
	s.Id = getNewId()
	s.Address = address
	s.Username = username
	s.Password = password
	s.Database = database
	log.Printf("New Influx-Session %s", s.Id)
	return s
}

func (s *Session) WritePoint(measurement string, input interface{}, date time.Time) error {
	once.Do(s.getInfluxClient)
	bp, err := ix.NewBatchPoints(ix.BatchPointsConfig{
		Database:  s.Database,
		Precision: "s",
	})
	if err != nil {
		return err
	}
	tags := map[string]string{"session_id": s.Id}
	fields := structs.Map(input)
	pt, err := ix.NewPoint(measurement, tags, fields, date)

	if err != nil {
		return err
	}

	bp.AddPoint(pt)
	return s.client.Write(bp)
}

func (s *Session) getInfluxClient() {
	c, err := ix.NewHTTPClient(ix.HTTPConfig{
		Addr:     s.Address,
		Username: s.Username,
		Password: s.Password,
	})
	if err != nil {
		log.Fatalln("Influx Error: ", err)
	}
	s.client = c
}

func getNewId() string {
	b := make([]byte, RAND_ID_SIZE)
	_, err := rand.Read(b)
	if err != nil {
		return "123456"
	}
	return hex.EncodeToString(b)
}
