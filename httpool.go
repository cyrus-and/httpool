// Package httpool provides wrappers for standard HTTP handlers that impose a
// limit to the number of Goroutines spawned to perform the actual work.
//
// The maximum number of active CPUs is not changed; it is left to the
// GOMAXPROCS environment variable.
package httpool

import (
	"errors"
	"net/http"
)

// Wrap accepts an ordinary http.Handler and produces another handler that uses
// a self-balancing pool of goroutine of the given size to serve the incoming
// requests. A pool size less than 1 will cause panic.
func Wrap(handler http.Handler, poolSize int) http.Handler {
	// require valid pool size
	if poolSize < 1 {
		panic(errors.New("the pool size must be at least 1"))
	}
	var token struct{}
	pool := make(chan struct{}, poolSize)
	// Return the wrapped handler.  Uses the buffered channel as a semaphore
	// to restrict the number of concurrent executions of the underlying
	// handler.
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		pool <- token           // wait for a slot to be available
		handler.ServeHTTP(w, r) // allow this request to proceed
		<-pool                  // release the slot
	})
}

// WrapFunc it is a helper that behaves like Wrap but accepts an
// http.HandlerFunc object instead.
func WrapFunc(function http.HandlerFunc, poolSize int) http.Handler {
	return Wrap(function, poolSize)
}
