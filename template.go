// Copyright 2013 errnoh. All rights reserved.
// Use of this source code is governed by a BSD-style (2-Clause)
// license that can be found in the LICENSE file.

package canvas

import (
	"html/template"
)

type templateVars struct {
	Address       string
	Width, Height int
}

var rootTemplate = template.Must(template.New("root").Parse(`<!DOCTYPE html>
<html>
<head>
<meta charset="utf-8" />

<body>
<canvas id="wdecanvas" width="{{.Width}}", height="{{.Height}}"></canvas>
</body>

<script type="text/javascript">

function onMessage(e) {
    start = Date.now();
    var byteArray = new Uint8ClampedArray(e.data);
    flush(byteArray);
    console.log("done: " + (Date.now()-start));
}

function onClose(e) {
    console.log("closed");
}

websocket = new WebSocket("ws://{{.Address}}/socket");
websocket.binaryType = 'arraybuffer';
websocket.onmessage = onMessage;
websocket.onclose = onClose;

var element = document.getElementById("wdecanvas");
var c = element.getContext("2d");

var locked = false;

function setSize(width, height) {
    if (!locked) {
        canvas1.width = width;
        canvas1.height = height;
    }
}

function size() {
    return [canvas1.width, canvas1.height];
}

function lockSize(bool) {
    locked = bool;
}

function title(s) {
    document.title = s;
}

function flush(arr) {
    imageData = c.createImageData(element.width, element.height); // Luoko aina uuden?
    //imageData = c.getImageData(0, 0, element.width, element.height);
    imageData.data.set(arr, 0);
    c.putImageData(imageData, 0, 0);
}

function close() {

}

wdecanvas.onclick = function(e) {
    e.preventDefault();
}

wdecanvas.oncontextmenu = function(e) {
    e.preventDefault();
}

// MouseDownEvent
wdecanvas.onmousewheel = function(e) {
    e.preventDefault();
    var ev = new Object();
    ev.Type = "MouseDownEvent";
    ev.Where = new Object();
    ev.Where.X = e.offsetX;
    ev.Where.Y = e.offsetY;
    if (e.wheelDeltaY >= 0) {
        ev.Which = 8;
    } else {
        ev.Which = 16;
    }
    websocket.send(JSON.stringify(ev));
}

// MouseDownEvent
wdecanvas.onmousedown = function(e) {
    e.preventDefault();
    var ev = new Object();
    ev.Type = "MouseDownEvent";
    ev.Where = new Object();
    ev.Which = 1 << e.button;
    ev.Where.X = e.offsetX;
    ev.Where.Y = e.offsetY;   
    websocket.send(JSON.stringify(ev));
}

// MouseUpEvent
wdecanvas.onmouseup = function(e) {
    e.preventDefault();
    var ev = new Object();
    ev.Type = "MouseUpEvent";
    ev.Where = new Object();
    ev.Which = 1 << e.button;
    ev.Where.X = e.offsetX;
    ev.Where.Y = e.offsetY;   
    websocket.send(JSON.stringify(ev));
}

// Addeventlistener?
// ResizeEvent
window.onresize = function(e) {
    var ev = new Object();
    ev.Type = "ResizeEvent";
    ev.Width = window.innerWidth;
    ev.Height = window.innerHeight;
    websocket.send(JSON.stringify(ev));
}

// MouseMovedEvent
wdecanvas.onmousemove = function(e) {
    e.preventDefault();
    var ev = new Object();
    ev.Type = "MouseMovedEvent";
    ev.Where = new Object();
    ev.Where.X = e.offsetX;
    ev.Where.Y = e.offsetY;   
    websocket.send(JSON.stringify(ev));
}

// KeyUpEvent
window.onkeyup = function(e) {
    e.preventDefault();
    var ev = new Object();
    ev.Type = "KeyUpEvent";
    ev.Which = e.which; 
    websocket.send(JSON.stringify(ev));
}

// KeyDownEvent
window.onkeydown = function(e) {
    e.preventDefault();
    var ev = new Object();
    ev.Type = "KeyDownEvent";
    ev.Which = e.which;
    websocket.send(JSON.stringify(ev));
}

window.onkeypress = function(e) {
    e.preventDefault();
}

// MouseEnteredEvent
wdecanvas.onmouseover = function(e) {
    e.preventDefault();
    var ev = new Object();
    ev.Type = "MouseEnteredEvent";
    ev.Where = new Object();
    ev.Where.X = e.offsetX;
    ev.Where.Y = e.offsetY; 
    // ev.From?
    websocket.send(JSON.stringify(ev));
}

// MouseExitedEvent
wdecanvas.onmouseout = function(e) {
    e.preventDefault();
    var ev = new Object();
    ev.Type = "MouseExitedEvent";
    ev.Where = new Object();
    ev.Where.X = e.offsetX;
    ev.Where.Y = e.offsetY;
    // ev.From?
    websocket.send(JSON.stringify(ev));
}

</script>
</html>
`))
