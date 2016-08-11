package tradelogger

import (
	"fmt"
	"io"
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
		fmt.Println("See logs at: " + c.Url + "/anim?lid=" + c.LogId)
		fmt.Println("See logs at: " + c.Url + "/live?lid=" + c.LogId)
		return c.Log(i)
	}
	return nil
}

func (l *serverLogger) LogHandler() func(http.ResponseWriter, *http.Request) {
	return func(rw http.ResponseWriter, req *http.Request) {
		if req.Method == http.MethodPost {
			item := parseItemFromRequest(req)

			id, ok := req.PostForm["lid"]
			if !ok {
				id = []string{DEFAULT_LOGGERID}
			}
			lid := id[0]
			if lid == DEFAULT_LOGGERID {
				lid = getNewLogId()
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
		} else {
			// Serve
		}
	}
}

// Logs are too big to serve as a chunk
func (l *serverLogger) GetRaw() func(http.ResponseWriter, *http.Request) {
	return func(rw http.ResponseWriter, req *http.Request) {
		lid := req.URL.Query().Get("lid")
		filename := l.Dirname + "/" + lid
		fh, err := os.OpenFile(filename, os.O_RDONLY, 0677)
		if err != nil {
			http.Error(rw, "Could not read the file", http.StatusInternalServerError)
			return
		}
		io.Copy(rw, fh)
	}
}

func (l *serverLogger) GetAnimated() func(http.ResponseWriter, *http.Request) {
	return func(rw http.ResponseWriter, req *http.Request) {

	}
}

func (l *serverLogger) GetLive() func(http.ResponseWriter, *http.Request) {
	return func(rw http.ResponseWriter, req *http.Request) {

	}
}

func getNewLogId() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
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
	if msg, ok := m["msg"]; ok {
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
	http.HandleFunc("/raw", l.GetRaw())
	http.HandleFunc("/anim", l.GetAnimated())
	http.HandleFunc("/live", l.GetLive())
	err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	if err != nil {
		fmt.Println("Error:", err.Error())
		os.Exit(1)
	}
}
