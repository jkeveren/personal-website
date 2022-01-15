package main

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"time"
)

// Implements http.Handler
type trickler struct {
	head []byte
	lines [][]byte
	lineDelay time.Duration
}

// Trickles content to client
func (t trickler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// print raw request for live analytics. This website does not need to be scalable.
	dump, _ := httputil.DumpRequest(r, false)
	fmt.Printf("Request from: %s\n%s", r.RemoteAddr, dump)

	flusher, ok := w.(http.Flusher)
	if !ok {
		w.WriteHeader(500)
		w.Write([]byte("server error"))
		fmt.Println("Could not get flusher")
		return
	}

	headers := w.Header()
	headers.Add("content-type", "text/html")
	// Caching does not allow trickle
	headers.Add("cache-control", "no-store")
	w.Write(t.head)

	for i, line := range t.lines {
		if i != 0 {
			time.Sleep(t.lineDelay * time.Millisecond)
		}
		w.Write(line)
		// Allows line by line trickle. Otherwise http will buffer more content before sending.
		flusher.Flush()
	}
}
