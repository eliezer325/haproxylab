package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"sync"
)

type contextKey int

const (
	connKey contextKey = iota + 1
)

type connTracker struct {
	conns  []net.Conn
	states []http.ConnState
	mtx    *sync.Mutex
}

func (ct *connTracker) logConnState(conn net.Conn, state http.ConnState) {
	ct.mtx.Lock()
	defer ct.mtx.Unlock()

	for i := range ct.conns {
		if conn == ct.conns[i] {
			fmt.Println("conn", i, state.String())
			ct.states[i] = state
			return
		}
	}

	ct.conns = append(ct.conns, conn)
	ct.states = append(ct.states, state)
}

func (ct *connTracker) logRequest(r *http.Request) {
	value := r.Context().Value(connKey)
	if value == nil {
		panic("no conn")
	}
	conn := value.(net.Conn)

	ct.mtx.Lock()
	defer ct.mtx.Unlock()

	for i := range ct.conns {
		if conn == ct.conns[i] {
			fmt.Println("conn", i, r.Method, r.URL.String())
			return
		}
	}

	panic("request conn not found")
}

func main() {
	tracker := connTracker{mtx: new(sync.Mutex)}

	srv := &http.Server{
		ConnContext: func(ctx context.Context, conn net.Conn) context.Context {
			tracker.logConnState(conn, http.StateNew)
			return context.WithValue(ctx, connKey, conn)
		},
		ConnState: tracker.logConnState,
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		tracker.logRequest(r)
		w.WriteHeader(200)
		w.Write([]byte("ok\n"))
	})

	fmt.Println("booting")
	fmt.Println(srv.ListenAndServe())
}
