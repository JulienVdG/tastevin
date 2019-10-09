// Copyright 2019 Splitted-Desktop Systems. All rights reserved
// Copyright 2019 Julien Viard de Galbert
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gotestweb

import (
	"bytes"
	"io"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

var (
	verbDebug = func(string, ...interface{}) {}
	//logDebug  = fmt.Printf
	logDebug = func(string, ...interface{}) {}
)

// HandleLive adds a websocket http.Handle for /live serving the test json written in the returned WriteCloser.
func HandleLive() io.WriteCloser {
	l := newLive()

	http.HandleFunc("/live", func(w http.ResponseWriter, r *http.Request) {
		ws, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			logDebug("%v\n", err)
			return
		}
		logDebug("Client subscribed %v\n", ws.RemoteAddr())
		go writer(ws, l)
		reader(ws)
		logDebug("Client unsubscribed %v\n", ws.RemoteAddr())
	})

	return l
}

type live struct {
	content   bytes.Buffer
	contentMu sync.RWMutex
	clients   map[*websocket.Conn]chan int
	clientsMu sync.RWMutex
	done      bool
}

func newLive() *live {
	var l live
	l.clients = make(map[*websocket.Conn]chan int)
	return &l
}

// Write buffers the test input and sent it to the websocket clients
func (l *live) Write(b []byte) (int, error) {
	verbDebug("write %d %q\n", len(b), b[len(b)-1])
	l.contentMu.Lock()
	n, err := l.content.Write(b)
	l.contentMu.Unlock()
	if err != nil {
		return n, err
	}
	l.clientsMu.RLock()
	defer l.clientsMu.RUnlock()
	for _, ch := range l.clients {
		ch <- l.content.Len()
	}
	return n, err
}

// Close mark the end of the test run
func (l *live) Close() error {
	l.done = true
	l.clientsMu.RLock()
	defer l.clientsMu.RUnlock()
	for _, ch := range l.clients {
		ch <- l.content.Len()
	}
	return nil
}

func (l *live) register(ws *websocket.Conn) chan int {
	logDebug("register %v\n", ws.RemoteAddr())
	ch := make(chan int)
	go func() {
		ch <- l.content.Len()
	}()
	l.clientsMu.Lock()
	l.clients[ws] = ch
	l.clientsMu.Unlock()
	return ch
}

func (l *live) unregister(ws *websocket.Conn) {
	logDebug("unregister %v\n", ws.RemoteAddr())
	l.clientsMu.Lock()
	defer l.clientsMu.Unlock()
	ch := l.clients[ws]
	close(ch)
	delete(l.clients, ws)
}

const (
	// Time allowed to write the file to the client.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the client.
	pongWait = 60 * time.Second

	// Send pings to client with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10
)

var upgrader = websocket.Upgrader{}

func reader(ws *websocket.Conn) {
	defer ws.Close()
	ws.SetReadLimit(512)
	ws.SetReadDeadline(time.Now().Add(pongWait))
	ws.SetPongHandler(func(string) error { ws.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		_, _, err := ws.ReadMessage()
		if err != nil {
			verbDebug("ReadMessage: %v\n", err)
			break
		}
	}
}

func writer(ws *websocket.Conn, l *live) {
	pingTicker := time.NewTicker(pingPeriod)
	defer func() {
		pingTicker.Stop()
		ws.Close()
	}()
	//ws.WriteMessage(websocket.TextMessage, []byte("Connected !"))
	ch := l.register(ws)
	defer l.unregister(ws)
	var last int
	for {
		select {
		case end := <-ch:
			if end > last {
				l.contentMu.RLock()
				line := l.content.Bytes()[last:end]
				last = end
				verbDebug("lines %d %q\n", len(line), line[len(line)-1])
				err := ws.WriteMessage(websocket.TextMessage, line)
				l.contentMu.RUnlock()
				if err != nil {
					verbDebug("WriteMessage: %v\n", err)
					return
				}
			}
			if l.done {
				return
			}

		case <-pingTicker.C:
			ws.SetWriteDeadline(time.Now().Add(writeWait))
			if err := ws.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
				verbDebug("Write PingMessage: %v\n", err)
				return
			}
		}
	}
}
