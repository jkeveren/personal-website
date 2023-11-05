package main

import (
	"bytes"
	"fmt"
	"net/http"
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
	lineDelayMilliseconds := 100 // delay between sending each line

	h := homeHandler{}

	h.head = h.makeHead(minHeadLength)
	h.lineDelay = time.Duration(lineDelayMilliseconds) * time.Millisecond

	content, err := web.ReadFile("web/homeContent.html")
	if err != nil {
		panic(err)
	}
	stringLines := bytes.SplitAfter(content, []byte("\n"))
	for _, stringLine := range stringLines {
		h.lines = append(h.lines, []byte(stringLine))
	}

	return h
}

// Trickles content to client
func (h homeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		w.Header().Add("Location", "/")
		w.WriteHeader(http.StatusTemporaryRedirect)
		return
	}

	flusher, ok := w.(http.Flusher)
	if !ok {
		w.WriteHeader(500)
		w.Write([]byte("server error"))
		fmt.Println("Could not get flusher")
		return
	}

	headers := w.Header()
	headers.Add("Content-Type", "text/html")
	// Caching does not allow trickle
	headers.Add("Cache-Control", "no-store")
	w.Write(h.head)

	doneChan := r.Context().Done()
	for _, line := range h.lines {
		select {
		case <-time.After(h.lineDelay):
			w.Write(line)
			// Allows line by line trickle. Otherwise net/http will buffer more content before sending.
			flusher.Flush()
		case <-doneChan:
			// Keeps tests fast by not wating for full response and doesnt send full response if client disconnects.
			break
		}
	}
}

func (h homeHandler) makeHead(minLength int) []byte {
	emptyBaseLength := len(h.makeParametricHead(0))
	paddingLength := minLength - emptyBaseLength
	return []byte(h.makeParametricHead(paddingLength))
}

// Returns HTML head with optional padding
func (h homeHandler) makeParametricHead(paddingLength int) string {
	title := "James Keveren"
	description := "Software Engineer"

	var padding string
	for i := 0; i < paddingLength; i++ {
		padding += "!"
	}

	// first line includes garbage used for black box test
	return `<!-- homem98y2r8 -->
<!DOCTYPE html>
<html style="font-size:16px;background:#000;color:#fff">
<head>
	<meta charset="utf-8">
	<title>` + title + `</title>
	<meta name=viewport content="initial-scale=0.6"/>
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
</head>
<body>
	<!-- this is a padding comment to ensure that browsers start rendering content as it arrives. Email me if you'd like more info. ` + padding + " -->\n"
}
