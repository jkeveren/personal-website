package main

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"os"
	"strings"
	"time"
)

// Minimum number of bytes to send before text will render as it arrives.
var minBaseLength = 1023 // 1023 for Firefox, 511 for Edge, 2 for Chrome
var base []byte
var lines [][]byte

var content = `<pre style="line-height: 1; letter-spacing: -1">
  █████████████████████████
  ██ ▄▄▄▄▄ █▀▀ ▀▀█ ▄▄▄▄▄ ██
  ██ █   █ █▀ █  █ █   █ ██
  ██ █▄▄▄█ █▀▀ ▀▀█ █▄▄▄█ ██
  ██▄▄▄▄▄▄▄█▄█ █▄█▄▄▄▄▄▄▄██
  ██▄▄▄▄▄▄▄▄▄▄▀▀▄ ▀ ▀ ▀ ▀██
  ██▀ ██▄█▄▄ ▄█▄▄ ▄▄▀▀ ▄███
  ██▄█▄▄█▄▄█▀ ▀▄██ ▀  ▀█ ██
  ██ ▄▄▄▄▄ █▄█▄▀  ▄ ▀▀ ▄▄██
  ██ █   █ █ █▄▄▄▀▄▄ ▄▀▀███
  ██ █▄▄▄█ █  █▄▄ █▄▀▀ ████
  ██▄▄▄▄▄▄▄█▄█▄▄▄▄▄▄█▄█▄███
  ▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀
</pre>James Keveren
<a target=_blank href=mailto:james@keve.ren>james@keve.ren</a>

I write software and make things.

Links:
- <a target=_blank href=https://instagram.com/jameskeveren>Instagram</a>
- <a target=_blank href=https://www.youtube.com/channel/UCOK3qJpL0I8qbXZsGzwYKzQ>YouTube</a>
- <a target=_blank href=https://www.tiktok.com/@jameskeveren>TikTok</a>
- <a target=_blank href=https://github.com/jkeveren>GitHub</a>
- <a target=_blank href=https://www.thingiverse.com/jkeveren/designs>Thingiverse</a>

Services:
- iPerf3: iperf.keve.ren (Maximum 370Mbs s->c / 35Mbs c->s)
`

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
		port = "5000"
	}
	fmt.Println("Starting HTTP Server on Port " + port + ". Configure using PORT environment variable.")
	panic(http.ListenAndServe(":"+port, handler{}))
}

type handler struct{}

func (handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	dump, _ := httputil.DumpRequest(r, false)
	fmt.Printf("%s", dump)

	// forwarded := r.Header.Get("forwarded")
	// if len(forwarded) == 0 {
	// forwarded = r.Header.Get("x-forwarded-for")
	// }
	//
	// fmt.Printf(
	// `Request: %s %s %s
	// IP: %v
	// Forwarded: %v
	// Agent: %v
	// Referrer: %v
	// `,
	// r.Method,
	// r.URL.Path,
	// r.Proto,
	// r.RemoteAddr,
	// forwarded,
	// r.Header.Get("user-agent"),
	// r.Header.Get("referer"),
	// )

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
			time.Sleep(50 * time.Millisecond)
		}
		w.Write(line)
		flusher.Flush()
	}
}

func makeBase(padding string) string {
	title := "James Keveren"
	description := "Software Developer"

	return `<meta charset="utf-8">
<html style="font-family:sans-serif;font-weight:bold;font-size:16px;background:#000;color:#fff">
<title>` + title + `</title>
<meta name=viewport content=width=device-width,user-scalable=no />
<meta name=title content="` + title + `">
<meta name=description content="` + description + `">
<link rel=icon type=image/png href=data:image/png>
<style>a{color:#ff0}</style>
<script async src=https://www.googletagmanager.com/gtag/js?id=UA-107575308-2></script>
<!-- The following is a padding scream to ensure that browsers start rendering content as it arrives. Email me for more info. ` + padding + "h -->\n"
}
