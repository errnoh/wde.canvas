wde.canvas
==========

HTML5 Canvas backend for [go.wde](https://github.com/skelterjohn/go.wde)

Usage
-----

When using go.wde, just import:

    _ "github.com/errnoh/wde.canvas"

And then connect to localhost:12500 with your browser.

Note:
-----

- Remember not to import multiple backends. (e.g. go.wde/init)
- Some functions (Resize for example) are not yet implemented.
- Browser needs websocket and canvas support.
