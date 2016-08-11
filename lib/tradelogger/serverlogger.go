package tradelogger

import (
	"crypto/rand"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"encoding/hex"

	bc "github.com/apourchet/investment/lib/broadcaster"
)

type serverLogger struct {
	Dirname      string
	Loggers      map[string]*localLogger
	Broadcasters map[string]*bc.Broadcaster
}

const (
	RAND_ID_SIZE = 10
)

func newServerLogger(dirname string) *serverLogger {
	l := &serverLogger{dirname, make(map[string]*localLogger), make(map[string]*bc.Broadcaster)}
	return l
}

func (l *serverLogger) LogHandler() func(http.ResponseWriter, *http.Request) {
	return func(rw http.ResponseWriter, req *http.Request) {
		req.ParseForm()
		id, ok := req.PostForm["lid"]
		if !ok {
			id = []string{DEFAULT_LOGGERID}
		}
		lid := id[0]

		if req.Method == http.MethodPost {
			item := parseItemFromRequest(req)

			if lid == DEFAULT_LOGGERID {
				lid = getNewLogId()
				fmt.Fprintf(rw, "%s", lid)
				return
			}

			if ll, ok := l.Loggers[lid]; ok {
				ll.Log(item)
				l.getOrCreateBroadcaster(lid).Emit(item)
				return
			}

			ll, err := NewLocalLogger(l.Dirname + "/" + lid)
			if err != nil {
				// TODO error response
				fmt.Println("COULD NOT CREATE NEW LOCAL LOGGER")
				return
			}
			l.Loggers[lid] = ll
			ll.Log(item)
			l.getOrCreateBroadcaster(lid).Emit(item)
		} else {
			// Serve that webpage

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
			http.Error(rw, "Could not read the file: "+err.Error(), http.StatusInternalServerError)
			return
		}
		io.Copy(rw, fh)
	}
}

func (l *serverLogger) GetLive() func(http.ResponseWriter, *http.Request) {
	return func(rw http.ResponseWriter, req *http.Request) {
		f, cn := checkStreamable(rw)
		if cn == nil {
			fmt.Fprintln(rw, `Not streamable`)
			return
		}

		setHeaders(rw)

		id, ok := req.PostForm["lid"]
		if !ok {
			id = []string{DEFAULT_LOGGERID}
		}
		lid := id[0]

		cb := make(chan interface{})
		rid := l.getOrCreateBroadcaster(lid).Register(cb)
		defer l.getOrCreateBroadcaster(lid).Deregister(rid)
		for data := range cb {
			item := data.(Item)
			fmt.Fprintf(rw, "data: %s\n\n", item.String())
			f.Flush()
		}
	}
}

func checkStreamable(rw http.ResponseWriter) (http.Flusher, http.CloseNotifier) {
	f, ok := rw.(http.Flusher)
	if !ok {
		http.Error(rw, "cannot stream", http.StatusInternalServerError)
		return nil, nil
	}

	cn, ok := rw.(http.CloseNotifier)
	if !ok {
		http.Error(rw, "cannot stream", http.StatusInternalServerError)
		return nil, nil
	}
	return f, cn
}

func setHeaders(rw http.ResponseWriter) {
	rw.Header().Set("Content-Type", "text/event-stream")
	rw.Header().Set("Cache-Control", "no-cache")
	rw.Header().Set("Connection", "keep-alive")
	rw.Header().Set("Access-Control-Allow-Origin", "*")
}

func (l *serverLogger) getOrCreateBroadcaster(lid string) *bc.Broadcaster {
	if _, ok := l.Broadcasters[lid]; !ok {
		l.Broadcasters[lid] = bc.NewBroadcaster()
	}
	return l.Broadcasters[lid]
}

func getNewLogId() string {
	b := make([]byte, RAND_ID_SIZE)
	rand.Read(b)
	return hex.EncodeToString(b)
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
