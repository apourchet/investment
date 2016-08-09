package tradelogger

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"time"

	"strings"

	"github.com/apourchet/investment/lib/utils"
)

type Logger interface {
	Log(Loggable) error
}

type serverLogger struct {
	Name     string
	Filename string
	fh       *os.File
}

type clientLogger struct {
	Url string
}

type Item struct {
	Date    time.Time `json:"date"`
	Tag     string    `json:"tag"`
	Message string    `json:"message"`
}

type Loggable interface {
	ToLogItem() *Item
}

var (
	mainLogger *serverLogger
)

func NewLogger(name, filename string) Logger {
	fh, err := os.Open(filename)
	if err != nil {
		return nil
	}
	return &serverLogger{name, filename, fh}
}

func NewLoggerClient(url string) Logger {
	return &clientLogger{url}
}

func (l *serverLogger) Log(item Loggable) error {
	i := item.ToLogItem()
	fmt.Fprintf(l.fh, "(%s) [%s]: %s", i.Date.String(), i.Tag, i.Message)
	return nil
}

func (c *clientLogger) Log(item Loggable) error {
	i := item.ToLogItem()
	v := url.Values{"date": {i.Date.String()},
		"tag":     {i.Tag},
		"message": {i.Message}}

	resp, err := http.PostForm(c.Url+"/log", v)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	_, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	return nil

}

func (i *Item) ToLogItem() *Item {
	return i
}

func LogHandler(rw *http.ResponseWriter, req http.Request) {
	if req.Method == http.MethodGet {
		// TODO
	} else if req.Method == http.MethodPost {
		req.ParseForm()
		m := req.PostForm
		dateStr := m["date"][0]
		dateArr := strings.Split(dateStr, ",")
		date, err := utils.ParseDate(dateArr)
		if err != nil {
			fmt.Fprintf(rw, err.Error())
			return
		}
		item := &Item{date, m["tag"][0], m["message"][0]}
		mainLogger.Log(item)
	}
}

func StartServer(port int, name, filename string) {
	mainLogger = NewLogger(name, filename)
	http.HandleFunc("/log", LogHandler)
	err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	if err != nil {
		fmt.Println("Error:", err.Error())
		os.Exit(1)
	}
}
