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
	"time"
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
		if bounds == nil {
			bounds = append(bounds, mainwindow.Screen().Bounds())
		}
		for _, r := range bounds {
			// NOTE: Would also be possible to create rgba images of those parts and send / update only those parts.
			// 		 context.putImageData(imgData,x,y)
			<-tick
			pos, delta := w.deltaSlice(r)
			w.send(struct {
				Type string
				Pos  int
			}{"pos", pos})
			websocket.Message.Send(w.ws, delta)
		}
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

func (w *Window) deltaSlice(r image.Rectangle) (int, []uint8) {
	// XXX: buf.Pix might need a mutex, can clash with newBuffer
	start := w.buf.PixOffset(r.Min.X, r.Min.Y)
	end := w.buf.PixOffset(r.Max.X, r.Max.Y)
	if start < 0 {
		start = 0
	}
	if end > len(w.buf.Pix) {
		end = len(w.buf.Pix)
	}
	return start, w.buf.Pix[start:end]
}
