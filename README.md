httpool
=======

Package `httpool` provides wrappers for standard HTTP handlers that impose a
limit to the number of Goroutines spawned to perform the actual work.

The maximum number of active CPUs is not changed; it is left to the `GOMAXPROCS`
environment variable.

Installation
------------

    go get github.com/cyrus-and/httpool

Example
-------

The following example shows how a regular `http.HandlerFunc` can be wrapped to
use at most 100 Goroutines and 4 CPUs to execute the handler.

```go
package main

import (
	"github.com/cyrus-and/httpool"
	"log"
	"net/http"
	"runtime"
)

func MyHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello, world!"))
}

func main() {
	runtime.GOMAXPROCS(4)
	h := httpool.WrapFunc(MyHandler, 100)
	log.Fatal(http.ListenAndServe(":8080", h))
}
```

Documentation
-------------

It can be found on [GoDoc](https://godoc.org/github.com/cyrus-and/httpool) or
directly in the [source file](httpool.go).
