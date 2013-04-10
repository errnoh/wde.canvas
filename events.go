// Copyright 2013 errnoh. All rights reserved.
// Use of this source code is governed by a BSD-style (2-Clause)
// license that can be found in the LICENSE file.

package canvas

import (
	"encoding/json"
	"fmt"
	"github.com/skelterjohn/go.wde"
	"image"
	"strconv"
)

// Union of all possible events
type superEvent struct {
	Type   string
	Which  int
	Where  image.Point
	From   image.Point
	Width  int
	Height int
}

func (se *superEvent) Clear() {
	se.Type = ""
	se.Which = 0
	se.Where.X, se.Where.Y, se.From.X, se.From.Y = 0, 0, 0, 0
	se.Width, se.Height = 0, 0
}

var downKeys = map[string]bool{}
var lastX, lastY int
var button uint

func forwardEvent(b []byte) {
	v := super
	v.Clear()
	json.Unmarshal(b, &v)

	var key string
	var ok bool
	if v.Which != 0 {
		if v.Which > 32 && v.Which < 91 {
			key = string(v.Which)
		} else {
			key, ok = keymap[v.Which]
			if !ok {
				key = strconv.Itoa(v.Which)
			}
		}

	}

	switch v.Type {
	case "MouseMovedEvent":
		mMoved.Where.X, mMoved.Where.Y = v.Where.X, v.Where.Y
		mMoved.From.X, mMoved.From.Y = lastX, lastY
		lastX, lastY = v.Where.X, v.Where.Y
		if button > 0 {
			mDragged.MouseMovedEvent = mMoved
			mDragged.Which = wde.Button(button)
			mainwindow.events <- mDragged
			break
		}
		mainwindow.events <- mMoved
	case "KeyDownEvent":
		downKeys[key] = true
		kDown.Key = key
		kTyped.Key = key
		kTyped.Chord = wde.ConstructChord(downKeys)
		mainwindow.events <- kDown
		mainwindow.events <- kTyped
	case "KeyUpEvent":
		delete(downKeys, key)
		kUp.Key = key
		mainwindow.events <- kUp
	case "MouseDownEvent":
		mDown.Where.X, mDown.Where.Y = v.Where.X, v.Where.Y
		mDown.Which = wde.Button(v.Which)
		lastX, lastY = v.Where.X, v.Where.Y
		if v.Which < 8 {
			// Mouse wheel doesn't send "MouseUp" event
			button = button | uint(v.Which)
		}
		mainwindow.events <- mDown
	case "MouseUpEvent":
		mUp.Where.X, mUp.Where.Y = v.Where.X, v.Where.Y
		mUp.Which = wde.Button(v.Which)
		lastX, lastY = v.Where.X, v.Where.Y
		button = button & ^uint(v.Which)
		mainwindow.events <- mUp
	case "MouseExitedEvent":
		mExited.Where.X, mExited.Where.Y = v.Where.X, v.Where.Y
		mExited.From.X, mExited.From.Y = lastX, lastY
		lastX, lastY = v.Where.X, v.Where.Y
		mainwindow.events <- mExited
	case "MouseEnteredEvent":
		mEntered.Where.X, mEntered.Where.Y = v.Where.X, v.Where.Y
		mEntered.From.X, mEntered.From.Y = lastX, lastY
		lastX, lastY = v.Where.X, v.Where.Y
		mainwindow.events <- mEntered
	case "ResizeEvent":
		wResize.Width, wResize.Height = v.Width, v.Height
		mainwindow.events <- wResize
	case "CloseEvent":
		mainwindow.events <- wClose
	default:
		fmt.Printf("%+v\n", v)
	}
}

var keymap = map[int]string{
	8:   wde.KeyBackspace,
	9:   wde.KeyTab,
	13:  "Enter",
	16:  "Shift",
	17:  "Ctrl",
	18:  "Alt",
	19:  "Pause, Break",
	20:  wde.KeyCapsLock,
	27:  wde.KeyEscape,
	32:  wde.KeySpace,
	33:  "Page Up",
	34:  "Page Down",
	35:  wde.KeyEnd,
	36:  wde.KeyHome,
	37:  wde.KeyLeftArrow,
	38:  wde.KeyUpArrow,
	39:  wde.KeyRightArrow,
	40:  wde.KeyDownArrow,
	44:  "PrntScrn",
	45:  wde.KeyInsert,
	46:  wde.KeyDelete,
	112: wde.KeyF1,
	113: wde.KeyF2,
	114: wde.KeyF3,
	115: wde.KeyF4,
	116: wde.KeyF5,
	117: wde.KeyF6,
	118: wde.KeyF7,
	119: wde.KeyF8,
	120: wde.KeyF9,
	121: wde.KeyF10,
	122: wde.KeyF11,
	123: wde.KeyF12,
	144: wde.KeyNumlock,
	145: "ScrollLock",
	188: ",",
	190: ".",
	191: "/",
	192: "`",
	219: "[",
	220: "\\",
	221: "]",
	222: "'",
}

// Reuse everything, structs are only copied when sent on a channel.
var (
	super  = superEvent{}
	mMoved = wde.MouseMovedEvent{
		MouseEvent: wde.MouseEvent{
			Where: image.Point{},
		},
	}
	mDragged = wde.MouseDraggedEvent{}
	kDown    = wde.KeyDownEvent{}
	kTyped   = wde.KeyTypedEvent{
		KeyEvent: wde.KeyEvent{},
	}
	kUp   = wde.KeyUpEvent{}
	mDown = wde.MouseDownEvent{
		MouseEvent: wde.MouseEvent{
			Where: image.Point{},
		},
	}
	mUp = wde.MouseUpEvent{
		MouseEvent: wde.MouseEvent{
			Where: image.Point{},
		},
	}
	mExited = wde.MouseExitedEvent{
		MouseEvent: wde.MouseEvent{
			Where: image.Point{},
		},
	}
	mEntered = wde.MouseEnteredEvent{
		MouseEvent: wde.MouseEvent{
			Where: image.Point{},
		},
	}
	wResize = wde.ResizeEvent{}
	wClose  = wde.CloseEvent{}
)
