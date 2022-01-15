package main

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httputil"
	"os"
	"time"
)

// Implements http.Handler
type homeHandler struct {
	head      []byte
	lines     [][]byte
	lineDelay time.Duration
}

func newHomeHandler() homeHandler {
	// Minimum number of bytes to send before text will render as it arrives.
	minHeadLength := 1023       // 1023 for Firefox, 511 for Edge, 2 for Chrome
	lineDelayMilliseconds := 20 // delay between sending each line

	t := homeHandler{}

	t.head = t.makeHead(minHeadLength)
	t.lineDelay = time.Duration(lineDelayMilliseconds)

	content, err := os.ReadFile("./content.html")
	if err != nil {
		panic(err)
	}
	stringLines := bytes.SplitAfter(content, []byte("\n"))
	for _, stringLine := range stringLines {
		t.lines = append(t.lines, []byte(stringLine))
	}

	return t
}

// Trickles content to client
func (t homeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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

func (t homeHandler) makeHead(minLength int) []byte {
	emptyBaseLength := len(t.makeParametricHead(0))
	paddingLength := minLength - emptyBaseLength
	return []byte(t.makeParametricHead(paddingLength))
}

// Returns HTML head with optional padding
func (t homeHandler) makeParametricHead(paddingLength int) string {
	title := "James Keveren"
	description := "Software Engineer"

	var padding string
	for i := 0; i < paddingLength; i++ {
		padding += "!"
	}

	return `<meta charset="utf-8">
<html style="font-size:16px;background:#000;color:#fff">
<title>` + title + `</title>
<meta name=viewport content=width=device-width,user-scalable=no />
<meta name=title content="` + title + `">
<meta name=description content="` + description + `">
<link rel=icon type=image/png href=data:image/png>
<style>
	a{color:#ff0}
	pre{margin:0}
	span{font-weight: normal}
</style>
<script async src=https://www.googletagmanager.com/gtag/js?id=UA-107575308-2></script>
<base target="_blank">
<!-- this is a padding comment to ensure that browsers start rendering content as it arrives. Email me if you'd like more info. ` + padding + " -->\n"
}
