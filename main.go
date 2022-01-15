package main

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"os"
	"bytes"
	"time"
)

// Minimum number of bytes to send before text will render as it arrives.
var minBaseLength = 1023 // 1023 for Firefox, 511 for Edge, 2 for Chrome
var lineDelayMilliseconds = 20 // delay between sending each line
var lineDelay time.Duration
var base []byte
var lines [][]byte

var content []byte

func main() {
	lineDelay = time.Duration(lineDelayMilliseconds)
	content, err := os.ReadFile("./content.html")
	if err != nil {
		panic(err)
	}

	emptyBaseLength := len(makeBase(""))
	paddingLength := minBaseLength - emptyBaseLength
	var padding string
	for i := 0; i < paddingLength; i++ {
		padding += "a"
	}
	base = []byte(makeBase(padding))

	stringLines := bytes.SplitAfter(content, []byte("\n"))
	for _, stringLine := range stringLines {
		lines = append(lines, []byte(stringLine))
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "5000"
	}
	fmt.Println("Starting HTTP Server on Port " + port + ". Configure using PORT environment variable.")
	panic(http.ListenAndServe(":"+port, handler{}))
}

type handler struct{}

func (handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	dump, _ := httputil.DumpRequest(r, false)
	fmt.Printf("Request from: %s\n%s", r.RemoteAddr, dump)

	flusher, ok := w.(http.Flusher)
	if !ok {
		w.WriteHeader(500)
		w.Write([]byte("server error"))
		panic("Could not get flusher")
		return
	}

	headers := w.Header()
	headers.Add("content-type", "text/html")
	headers.Add("cache-control", "no-store")
	w.Write(base)

	for i, line := range lines {
		if i != 0 {
			time.Sleep(lineDelay * time.Millisecond)
		}
		w.Write(line)
		flusher.Flush()
	}
}

func makeBase(padding string) string {
	title := "James Keveren"
	description := "Software Developer"

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
<!-- this is a padding comment to ensure that browsers start rendering content as it arrives. Email me if you'd like more info. ` + padding + "h -->\n"
}
