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
	LockSize      bool
}

var rootTemplate = template.Must(template.New("root").Parse(`<!DOCTYPE html>
<html>
<head>
<meta charset="utf-8" />

<body>
<canvas id="wdecanvas" width="{{.Width}}", height="{{.Height}}"></canvas>
</body>

<script type="text/javascript">

websocket = new WebSocket("ws://{{.Address}}/socket");
websocket.binaryType = 'arraybuffer';
websocket.onmessage = onMessage;
websocket.onclose = onClose;

var element = document.getElementById("wdecanvas");
var c = element.getContext("2d");

var locked = {{.LockSize}};
var bufsize = element.width * element.height * 4;

function onMessage(e) {
    start = Date.now();
    if (e.data.byteLength == null) {
        // Data is not bytearray
        var data = JSON.parse(e.data);
        handle(data);
        return;
    }
    var byteArray = new Uint8ClampedArray(e.data);
    flush(byteArray);
    websocket.send("ok");
    console.log("done: " + (Date.now()-start));
}

function handle(data) {
    switch (data.Type) {
    case "resize":
        setSize(data.Width, data.Height);
    case "title":
        title(data.Title);
    case "icon":
    case "iconname":
    case "locksize":
        lockSize(data.Bool);
    }
}

function onClose(e) {
    close();
    console.log("closed");
}

function setSize(width, height) {
    if (!locked) {
        element.width = width;
        element.height = height;
        var bufsize = element.width * element.height * 4;
    }
}

function size() {
    return [element.width, element.height];
}

function lockSize(value) {
    console.log(locked, value);
    locked = value;
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

function getMousePos(canvas, evt) {
    var rect = canvas.getBoundingClientRect();
    return {
        X: evt.clientX - rect.left,
        Y: evt.clientY - rect.top
    };
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
    ev.Where = getMousePos(element, e);
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
    ev.Which = 1 << e.button;
    ev.Where = getMousePos(element, e);
    websocket.send(JSON.stringify(ev));
}

// MouseUpEvent
wdecanvas.onmouseup = function(e) {
    e.preventDefault();
    var ev = new Object();
    ev.Type = "MouseUpEvent";
    ev.Which = 1 << e.button;
    ev.Where = getMousePos(element, e);
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
    ev.Where = getMousePos(element, e);
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
    ev.Where = getMousePos(element, e);
    // ev.From?
    websocket.send(JSON.stringify(ev));
}

// MouseExitedEvent
wdecanvas.onmouseout = function(e) {
    e.preventDefault();
    var ev = new Object();
    ev.Type = "MouseExitedEvent";
    ev.Where = getMousePos(element, e);
    // ev.From?
    websocket.send(JSON.stringify(ev));
}

</script>
</html>
`))
