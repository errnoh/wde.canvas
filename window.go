// Copyright 2013 errnoh. All rights reserved.
// Use of this source code is governed by a BSD-style (2-Clause)
// license that can be found in the LICENSE file.

package canvas

import (
	"code.google.com/p/go.net/websocket"
	"github.com/skelterjohn/go.wde"
	"image"
	"image/draw"
	"net/http"
)

type RGBA struct {
	*image.RGBA
}

func (img RGBA) CopyRGBA(src *image.RGBA, bounds image.Rectangle) {
	draw.Draw(img.RGBA, bounds, src, image.ZP, draw.Src)
}

func (img RGBA) AsRGBA() *image.RGBA {
	return img.RGBA
}

type Window struct {
	ws     *websocket.Conn
	buf    *RGBA
	events chan interface{}
	end    chan struct{}

	started bool
}

func NewWindow(width, height int) (w wde.Window, err error) {
	if !mainwindow.started {
		tempvars.Address = listenAddr
		tempvars.Width = width
		tempvars.Height = height

		mainwindow.buf = &RGBA{image.NewRGBA(image.Rect(0, 0, width, height))}
		mainwindow.events = make(chan interface{})
		mainwindow.end = make(chan struct{})

		http.HandleFunc("/", rootHandler)
		http.Handle("/socket", websocket.Handler(socketHandler))

		go http.ListenAndServe(listenAddr, nil)
		mainwindow.started = true
	}

	return mainwindow, nil
}

func (w *Window) Close() (err error) {
	w.end <- struct{}{}
	return nil
}

func (w *Window) EventChan() (events <-chan interface{}) {
	return w.events
}

// XXX: Flushing too fast might cause pipes to get stuck.
//
// Workaround: Limit FPS.
//
// FIX: Possibly send JSON response after each frame to signal ready
func (w *Window) FlushImage(bounds ...image.Rectangle) {
	if w.ws != nil {
		websocket.Message.Send(w.ws, w.buf.Pix)
	}
}

func (w *Window) Screen() (im wde.Image) {
	return w.buf
}

// NOTE: Not yet supported
func (w *Window) LockSize(lock bool) {}

// NOTE: Not yet supported
func (w *Window) SetIcon(icon image.Image) {}

// NOTE: Not yet supported
func (w *Window) SetIconName(name string) {}

func (w *Window) SetSize(width, height int) {
	var data struct {
		Type          string
		Width, Height int
	}
	data.Type, data.Width, data.Height = "resize", width, height
	websocket.JSON.Send(w.ws, data)

	w.newBuffer(width, height)
}

// NOTE: Not yet supported
func (w *Window) SetTitle(title string) {}

// NOTE: Doesn't do anything yet.
func (w *Window) Show() {}

// Returns width and height of the canvas.
func (w *Window) Size() (width, height int) {
	return w.Screen().Bounds().Dx(), w.Screen().Bounds().Dy()
}

// If the Window is resized we need new buffer since the size changes.
// To minimize losses as much of the old screen is copied to new buffer as possible.
func (w *Window) newBuffer(width, height int) {
	oldbuf := mainwindow.buf
	mainwindow.buf = &RGBA{image.NewRGBA(image.Rect(0, 0, width, height))}
	mainwindow.Screen().CopyRGBA(oldbuf.RGBA, oldbuf.Bounds())
}
