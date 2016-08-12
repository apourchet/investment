package tradelogger

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"time"
)

type Logger interface {
	Log(Loggable) error
}

type localLogger struct {
	Filename string
	fh       *os.File
}

type clientLogger struct {
	LogId string
	Url   string
}

type Item struct {
	Date    time.Time `json:"date"`
	Tag     string    `json:"tag"`
	Message string    `json:"message"`
}

type Loggable interface {
	ToLogItem() *Item
}

const (
	DEFAULT_LOGGERID = "unnamed"
)

func NewLocalLogger(filename string) (*localLogger, error) {
	if _, err := os.Stat(filename); err == nil {
		err = os.Rename(filename, filename+".old")
	}
	fh, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY, 0667)
	if err != nil {
		return nil, err
	}
	return &localLogger{filename, fh}, nil
}

func NewLoggerClient(url string) *clientLogger {
	return &clientLogger{DEFAULT_LOGGERID, url}
}

func (l *localLogger) Log(data Loggable) error {
	item := data.ToLogItem()
	fmt.Fprintln(l.fh, item.String())
	return nil
}

func (c *clientLogger) Log(item Loggable) error {
	i := item.ToLogItem()
	if i == nil {
		i = &Item{time.Now(), "ERROR", "Could not convert to log item."}
	}
	v := url.Values{
		"lid":  {c.LogId},
		"date": {i.Date.Format(time.UnixDate)},
		"tag":  {i.Tag},
		"msg":  {i.Message}}

	resp, err := http.PostForm(c.Url+"/log", v)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	idData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if c.LogId == DEFAULT_LOGGERID {
		// TODO error response case
		c.LogId = string(idData)
		fmt.Println("See logs at: " + c.Url + "/log?lid=" + c.LogId)
		fmt.Println("See logs at: " + c.Url + "/raw?lid=" + c.LogId)
		return c.Log(i)
	}
	return nil
}

func (i *Item) ToLogItem() *Item {
	return i
}

func (i *Item) String() string {
	return fmt.Sprintf("(%s) [%s]: %s\n", i.Date.Format(time.UnixDate), i.Tag, i.Message)
}

func StartServer(port int, dirname string) {
	l := newServerLogger(dirname)
	http.HandleFunc("/log", l.LogHandler())
	http.HandleFunc("/raw", l.GetRaw())
	http.HandleFunc("/live", l.GetLive())
	err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	if err != nil {
		fmt.Println("Error:", err.Error())
		os.Exit(1)
	}
}
