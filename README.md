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
use at most 100 Goroutines to execute the handler.

```go
package main

import (
	"github.com/cyrus-and/httpool"
	"log"
	"net/http"
)

func MyHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello, world!"))
}

func main() {
	h := httpool.WrapFunc(MyHandler, 100)
	log.Fatal(http.ListenAndServe(":8080", h))
}
```