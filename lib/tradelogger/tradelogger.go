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

func NewLogger(name, filename string) (*serverLogger, error) {
	if _, err := os.Stat(filename); err == nil {
		err = os.Rename(filename, filename+".old")
	}
	fh, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY, 0667)
	if err != nil {
		return nil, err
	}
	return &serverLogger{name, filename, fh}, nil
}

func NewLoggerClient(url string) *clientLogger {
	return &clientLogger{url}
}

func (l *serverLogger) Log(item Loggable) error {
	i := item.ToLogItem()
	fmt.Fprintf(l.fh, "(%s) [%s]: %s\n", i.Date.Format(time.UnixDate), i.Tag, i.Message)
	return nil
}

func (c *clientLogger) Log(item Loggable) error {
	i := item.ToLogItem()
	if i == nil {
		i = &Item{time.Now(), "ERROR", "Could not convert to log item."}
	}
	v := url.Values{"date": {i.Date.Format(time.UnixDate)},
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

func (l *serverLogger) LogHandler() func(http.ResponseWriter, *http.Request) {
	return func(rw http.ResponseWriter, req *http.Request) {
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
			l.Log(item)
		}
	}
}

func StartServer(port int, name, filename string) {
	l, err := NewLogger(name, filename)
	http.HandleFunc("/log", l.LogHandler())
	err = http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	if err != nil {
		fmt.Println("Error:", err.Error())
		os.Exit(1)
	}
}
