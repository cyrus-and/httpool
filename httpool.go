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

type task struct {
	w   http.ResponseWriter
	r   *http.Request
	end chan bool
}

// Wrap accepts an ordinary http.Handler and produces another handler that uses
// a self-balancing pool of goroutine of the given size to serve the incoming
// requests. A pool size less than 1 will cause panic.
func Wrap(handler http.Handler, poolSize int) http.Handler {
	// require valid (and useful) pool size
	if poolSize < 1 {
		panic(errors.New("the pool size must be at least 1"))
	}
	// create poolSize goroutines
	taskChans := make([]chan *task, poolSize)
	readyChan := make(chan int)
	for id := 0; id < poolSize; id++ {
		id := id
		taskChans[id] = make(chan *task)
		go func() {
			// initially mark this worker as ready
			readyChan <- id
			for task := range taskChans[id] {
				// for each task serve the associated request
				handler.ServeHTTP(task.w, task.r)
				// finish this task then finally again mark this worker as
				// ready; the order is important so the associated goroutine
				// created by the http Go library can terminate before a new
				// task is accepted for this worker
				task.end <- true
				readyChan <- id
			}
		}()
	}
	// return the wrapped handler: fetch the next available worker id, send the
	// task to that channel and awaits for its completion
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		readyId := <-readyChan
		task := &task{w, r, make(chan bool)}
		taskChans[readyId] <- task
		<-task.end
	})
}

// WrapFunc it is a helper that behaves like Wrap but accepts an
// http.HandlerFunc object instead.
func WrapFunc(function http.HandlerFunc, poolSize int) http.Handler {
	return Wrap(function, poolSize)
}
