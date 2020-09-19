package main

import (
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

// Minimum number of bytes to send before text will render as it arrives.
var minBaseLength = 4095 // 4095 for IE, 1023 for Firefox, 511 for Edge, 2 for Chrome
var base []byte
var url = "https://james.keve.ren"
var lines [][]byte

var content = `James Keveren
<a target=_blank href=` + url + `>james.keve.ren</a>
<a target=_blank href=mailto:james@keve.ren>james@keve.ren</a>

I write software and make things.

Links:
- <a target=_blank href=https://github.com/jkeveren>GitHub</a>
- <a target=_blank href=https://instagram.com/jameskeveren>Instagram</a>
- <a target=_blank href=https://www.thingiverse.com/jkeveren/designs>Thingiverse</a>`

func main() {
	emptyBaseLength := len(makeBase(""))
	paddingLength := minBaseLength - emptyBaseLength
	var padding string
	for i := 0; i < paddingLength; i++ {
		padding += "a"
	}
	base = []byte(makeBase(padding))

	content := strings.Replace(content, "\n", "<br>", -1)
	stringLines := strings.SplitAfter(content, "<br>")
	for _, stringLine := range stringLines {
		lines = append(lines, []byte(stringLine))
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "50000"
	}
	log.Print("Starting HTTP Server on Port " + port + ". Configure using PORT environment variable.")
	log.Panic(http.ListenAndServe(":"+port, handler{}))
}

func makeBase(padding string) string {
	title := "James Keveren"
	description := "Software Developer"

	return `<html style="font-family:monospace; background: black; color: white">
	<title>` + title + `</title>
	<meta name=viewport content=width=device-width,user-scalable=no />
	<meta name="title" content="` + title + `">
	<meta name="description" content="` + description + `">
	<meta property="og:type" content="website">
	<meta property="og:url" content="` + url + `">
	<meta property="og:title" content="` + title + `">
	<meta property="og:description" content="` + description + `">
	<meta property="twitter:url" content="` + url + `">
	<meta property="twitter:title" content="` + title + `">
	<meta property="twitter:description" content="` + description + `">
	<link rel="icon" type="image/png" href="data:image/png">
	<style>a {color: #ff0}</style>
	<!-- This is just a filler comment to consume a few bytes so browsers start rendering content as it arrives.

	Here's how many bytes it takes for each browser to start rendering HTML as each byte arrives:
	- Google Chrome 78.0.3904.70:	   2
	- Microsoft Edge 44.18362.387.0:	511 (probably different since Blink)
	- Mozilla Firefox 69.0.3:		   1023
	- Internet Explorer 11.418.18362.0: 4095

	Anyway here's some padding garbage for IE compatability: ` + padding + "h -->\n"
}

type handler struct{}

func (handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Print("----------------------Request----------------------")
	log.Printf("IP: %v", r.RemoteAddr)
	for name, values := range r.Header {
		for _, value := range values {
			log.Printf("%v: %v", name, value)
		}
	}

	flusher, ok := w.(http.Flusher)
	if !ok {
		w.WriteHeader(500)
		w.Write([]byte("server error"))
		log.Panic("Could not get flusher")
		return
	}

	headers := w.Header()
	headers.Add("content-type", "text/html")
	headers.Add("cache-control", "no-store")
	w.Write(base)

	for i, line := range lines {
		if i != 0 {
			time.Sleep(100 * time.Millisecond)
		}
		w.Write(line)
		flusher.Flush()
	}
}
