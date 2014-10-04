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
