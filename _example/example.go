// Copyright 2013 errnoh. All rights reserved.
// Use of this source code is governed by a BSD-style (2-Clause)
// license that can be found in the LICENSE file.

package main

import (
	_ "github.com/errnoh/wde.canvas"
	"fmt"
	"github.com/skelterjohn/go.wde"
	"image"
	"image/color"
	"image/draw"
	"math/rand"
	"time"
)

const (
	WIDTH, HEIGHT = 600, 600
)

var (
	mousex, mousey int
	posx, posy     int
	radius         int
	r, g, b, a     uint8

	dw   wde.Window
	done = make(chan struct{})
)

func main() {
	go run()
	wde.Run()
}

func run() {
	var (
		err error
	)
	radius = 10
	r, g, b, a = 0, 0, 170, 255

	dw, err = wde.NewWindow(WIDTH, HEIGHT)
	if err != nil {
		fmt.Println(err)
		return
	}
	dw.Show()

	events := dw.EventChan()
	go listen(events)
	render()
}

func listen(c <-chan interface{}) {
loop:
	for ev := range c {
		switch e := ev.(type) {
		case wde.MouseDownEvent:
			if e.Which == 8 {
				radius++
			} else if e.Which == 16 {
				radius--
			} else {
				randomize()
			}
		case wde.MouseMovedEvent:
			mousex, mousey = e.Where.X, e.Where.Y
		case wde.CloseEvent:
			fmt.Println("closed")
			dw.Close()
			break loop
		}
	}
	done <- struct{}{}
}

func render() {
	for {
		draw.DrawMask(dw.Screen(), dw.Screen().Bounds(), &image.Uniform{color.RGBA{r, g, b, a}}, image.ZP, &circle{image.Point{mousex, mousey}, radius}, image.ZP, draw.Over)
		dw.FlushImage()
		select {
		case <-time.After(time.Second / 25):
		case <-done:
			return
		}
	}
}

func randomize() {
	r, g, b = uint8(rand.Uint32()), uint8(rand.Uint32()), uint8(rand.Uint32())
}

type circle struct {
	p image.Point
	r int
}

func (c *circle) ColorModel() color.Model {
	return color.AlphaModel
}

func (c *circle) Bounds() image.Rectangle {
	return image.Rect(c.p.X-c.r, c.p.Y-c.r, c.p.X+c.r, c.p.Y+c.r)
}

func (c *circle) At(x, y int) color.Color {
	xx, yy, rr := float64(x-c.p.X)+0.5, float64(y-c.p.Y)+0.5, float64(c.r)
	if xx*xx+yy*yy < rr*rr {
		return color.Alpha{255}
	}
	return color.Alpha{0}
}
