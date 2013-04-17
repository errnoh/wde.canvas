// Copyright 2013 errnoh. All rights reserved.
// Use of this source code is governed by a BSD-style (2-Clause)
// license that can be found in the LICENSE file.

package canvas

import (
	"code.google.com/p/go.net/websocket"
	"github.com/skelterjohn/go.wde"
	"image"
	"image/color"
	"image/draw"
	"net/http"
	"time"
)

type RGBA struct {
	*image.RGBA
}

func (img RGBA) Set(x, y int, c color.Color) {
	// XXX: draw.Draw is not using image.RGBA but canvas.RGBA and Set()
	// Because of this it's probably still faster to do manual buffering and CopyRGBA it
	mainwindow.deltaUpdate(x, y)
	img.RGBA.Set(x, y, c)
}

func (img RGBA) SetRGBA(x, y int, c color.RGBA) {
	mainwindow.deltaUpdate(x, y)
	img.RGBA.SetRGBA(x, y, c)
}

func (img RGBA) CopyRGBA(src *image.RGBA, bounds image.Rectangle) {
	mainwindow.deltaUpdate(bounds.Min.X, bounds.Min.Y)
	mainwindow.deltaUpdate(bounds.Max.X-1, bounds.Max.Y-1)
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
	delta  image.Rectangle // TODO: new connection == delta whole screen

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
		mainwindow.delta = mainwindow.Screen().Bounds()
		mainwindow.delta.Max.X--
		mainwindow.delta.Max.Y--

		http.HandleFunc("/", rootHandler)
		http.Handle("/socket", websocket.Handler(socketHandler))

		go http.ListenAndServe(listenAddr, nil)
		mainwindow.started = true
	}

	return mainwindow, nil
}

var fpscount int

func (w *Window) FPS() <-chan int {
	c := make(chan int)
	ticker := time.Tick(time.Second)
	go func() {
		for {
			fpscount = 0
			<-ticker
			c <- fpscount
		}
	}()
	return c
}

func (w *Window) Close() (err error) {
	w.ws.Close()

	w.end <- struct{}{}
	return nil
}

func (w *Window) EventChan() (events <-chan interface{}) {
	return w.events
}

func (w *Window) FlushImage(bounds ...image.Rectangle) {
	if w.ws != nil {
		<-tick
		websocket.Message.Send(w.ws, w.deltaSlice())
		w.delta.Min.X, w.delta.Min.Y, w.delta.Max.X, w.delta.Max.Y = -1, -1, -1, -1
		fpscount++
	} else {
		// NOTE: Prevent locking in busy wait when there's no client available
		time.Sleep(time.Millisecond * 500)
	}
}

func (w *Window) Screen() (im wde.Image) {
	return w.buf
}

func (w *Window) LockSize(lock bool) {
	var data struct {
		Type string
		Bool bool
	}
	tempvars.LockSize = lock
	data.Type, data.Bool = "locksize", lock
	w.send(data)
}

// NOTE: Not yet supported
func (w *Window) SetIcon(icon image.Image) {}

// NOTE: Not yet supported
func (w *Window) SetIconName(name string) {}

func (w *Window) SetSize(width, height int) {
	if tempvars.LockSize {
		return
	}

	var data struct {
		Type          string
		Width, Height int
	}
	data.Type, data.Width, data.Height = "resize", width, height
	tempvars.Width, tempvars.Height = width, height
	w.send(data)

	w.newBuffer(width, height)
}

func (w *Window) SetTitle(title string) {
	var data struct {
		Type  string
		Title string
	}
	data.Type, data.Title = "title", title
	tempvars.Title = title
	w.send(data)
}

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

func (w *Window) send(data interface{}) {
	if w.ws != nil {
		websocket.JSON.Send(w.ws, data)
	}
}

func (w *Window) deltaUpdate(x, y int) {
	if w.delta.Min.X == -1 || w.delta.Max.X == -1 {
		w.delta.Min.X, w.delta.Min.Y, w.delta.Max.X, w.delta.Max.Y = x, y, x, y
		return
	}
	// NOTE: Doesn't check if x,y is inside w.Bounds rectangle.
	// NOTE: Could also support ... syntax with image.Point
	// Also consider using only the row and skip column to reduce overhead
	if y <= w.delta.Min.Y {
		if x < w.delta.Min.X {
			w.delta.Min.Y = y
			w.delta.Min.X = x
		}
	}
	if y >= w.delta.Max.Y {
		w.delta.Max.Y = y
		w.delta.Max.X = x
		if x > w.delta.Max.X {
			w.delta.Max.Y = y
			w.delta.Max.X = x
		}
	}
}

func (w *Window) deltaSlice() []uint8 {

	start := w.buf.PixOffset(w.delta.Min.X, w.delta.Min.Y)
	end := w.buf.PixOffset(w.delta.Max.X, w.delta.Max.Y)
	if end > len(w.buf.Pix) {
		end = len(w.buf.Pix)
	}
	start = start // debug
	return w.buf.Pix[0:end]
}
