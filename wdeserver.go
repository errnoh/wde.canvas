// Copyright 2013 errnoh. All rights reserved.
// Use of this source code is governed by a BSD-style (2-Clause)
// license that can be found in the LICENSE file.

package canvas

import (
	"code.google.com/p/go.net/websocket"
	"github.com/skelterjohn/go.wde"
	"net/http"
)

const listenAddr = "localhost:12500"

var (
	mainwindow = new(Window)
	tempvars   = new(templateVars)
	tick       = make(chan struct{})
)

func init() {
	wde.BackendNewWindow = NewWindow
	ch := make(chan struct{}, 1)
	wde.BackendRun = func() {
		<-ch
	}
	wde.BackendStop = func() {
		ch <- struct{}{}
	}
}

func socketHandler(ws *websocket.Conn) {
	mainwindow.ws = ws
	go listen(ws)
	<-mainwindow.end
}

func listen(ws *websocket.Conn) {
	// NOTE: Could use websocket.JSON.Receive
	var buf = make([]byte, 512)
	for {
		n, err := ws.Read(buf)
		if err != nil {
			println(err.Error())
			ws.Close()
			return
		}
		if string(buf[:n]) == "ok" {
			tick <- struct{}{}
			continue
		}
		forwardEvent(buf[:n])
	}
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	rootTemplate.Execute(w, tempvars)
}
