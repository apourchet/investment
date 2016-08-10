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

type serverLogger struct {
	Dirname string
	Loggers map[string]*localLogger
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

func newServerLogger(dirname string) *serverLogger {
	l := &serverLogger{dirname, make(map[string]*localLogger)}
	return l
}

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
	fmt.Fprintf(l.fh, "(%s) [%s]: %s\n",
		item.Date.Format(time.UnixDate), item.Tag, item.Message)
	return nil
}

func (c *clientLogger) Log(item Loggable) error {
	i := item.ToLogItem()
	if i == nil {
		i = &Item{time.Now(), "ERROR", "Could not convert to log item."}
	}
	v := url.Values{
		"id":      {c.LogId},
		"date":    {i.Date.Format(time.UnixDate)},
		"tag":     {i.Tag},
		"message": {i.Message}}

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
		fmt.Println("See logs at: " + c.Url + "/logs?id=" + c.LogId)
		return c.Log(item)
	}
	return nil
}

func (l *serverLogger) LogHandler() func(http.ResponseWriter, *http.Request) {
	return func(rw http.ResponseWriter, req *http.Request) {
		if req.Method == http.MethodGet {
			fmt.Fprintf(rw, "OK")
			// TODO
		} else if req.Method == http.MethodPost {
			item := parseItemFromRequest(req)

			id, ok := req.PostForm["id"]
			if !ok {
				id = []string{DEFAULT_LOGGERID}
			}
			lid := id[0]
			if lid == DEFAULT_LOGGERID {
				// TODO generate new id
				lid = "1234"
				fmt.Fprintf(rw, "%s", lid)
				return
			}

			if ll, ok := l.Loggers[lid]; ok {
				ll.Log(item)
				return
			}

			ll, err := NewLocalLogger(l.Dirname + "/" + lid)
			if err != nil {
				// TODO error response
				return
			}
			l.Loggers[lid] = ll
			ll.Log(item)
		}
	}
}

func parseItemFromRequest(req *http.Request) *Item {
	req.ParseForm()
	m := req.PostForm

	item := &Item{}
	if d, ok := m["date"]; ok {
		item.Date, _ = time.Parse(time.UnixDate, d[0])
	}
	if t, ok := m["tag"]; ok {
		item.Tag = t[0]
	}
	if msg, ok := m["tag"]; ok {
		item.Message = msg[0]
	}
	return item
}

func (i *Item) ToLogItem() *Item {
	return i
}

func StartServer(port int, dirname string) {
	l := newServerLogger(dirname)
	http.HandleFunc("/log", l.LogHandler())
	err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	if err != nil {
		fmt.Println("Error:", err.Error())
		os.Exit(1)
	}
}
