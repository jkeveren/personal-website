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
var lineDelayMilliseconds = 20 // delay between sending each line
var lineDelay time.Duration
var base []byte
var lines [][]byte

var content = `<pre style="font-family: sans-serif; font-weight: bold"><span style="font-family: monospace;line-height: 1; letter-spacing: -1">
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
</span>
James Keveren
<a target=_blank href=mailto:james@keve.ren>james@keve.ren</a>

Software Engineer and Maker of things.
<span>Experience with Go, JavaScript, C#, and a little C++.</span>

Links:
- <a href=https://github.com/jkeveren>GitHub</a>
	<span>Mostly private but you can see some of my hobby projects.</span>
- <a href=gallery>Gallery</a>
	<span>This documents my non-software projects without depening on social media.</span>
- <a href=https://www.youtube.com/channel/UCOK3qJpL0I8qbXZsGzwYKzQ>YouTube</a>
	<span>Infrequent long form videos on projects that lend themselves to video.</span>
- <a href=https://www.tiktok.com/@jameskeveren>TikTok</a>, <a target=_blank href=https://instagram.com/jameskeveren>Instagram</a>
	<span>Messing around in my workshop.</span>
- <a href=https://www.thingiverse.com/jkeveren/designs>Thingiverse</a>
	<span>Old 3D printing projects.</span>

Live Services:
- <a href=https://cam.keve.ren>Garden Cam</a>: <span>Portal into my garden. Sometimes the sun hides behind the planet.</span>
- <a href=https://massdraw.keve.ren>MassDraw</a>: <span>Exremely simple multi-user whiteboard.</span>
- iPerf3: <span>iperf -c iperf.keve.ren (Maximum 350Mbps c->s / 35Mbps c<-s)</span>

Skills and Technology <span>(This is optimised for scraping but if you're human that's ok too)</span>:
- Software Development, Software Engineering<span>
	Golang, JavaScript, Nodejs, NPM, C#, .Net, Entity Framework, C++, SQL, MSSQL, NoSQL, MongoDB, Git, Mercurial, Magic++, TDD, Regexp, Mocha, Express, Gulp, Pug, HTML, CSS, Websockets, JSON, HTTP</span>
- System Administration<span>
	GCP, firebase, AWS, Linux (Arch, Debian, CentOS), Systemd, Fish, Bash, ssh, Haproxy, Nginx, Plesk, dm-crypt, Xinetd, Rsync, Cage, FFmpeg, LTO, Dell iDrac, Raspberry pi</span>
- Networking<span>
	Opnsense, Unifi, Poe, HTTP, DNS, TLS, SSL, Letsencrypt, Certbot, Fibre Channel</span>
- IT support<span>
	Microsoft 365, MailEnable, Desktop Hardware</span>
- Non-IT<span>
	SolidWorks, Fusion360, OpenScad, Blender, Cura, Chitubox, Kdenlive, </span>
</pre>
`

func main() {
	lineDelay = time.Duration(lineDelayMilliseconds)

	emptyBaseLength := len(makeBase(""))
	paddingLength := minBaseLength - emptyBaseLength
	var padding string
	for i := 0; i < paddingLength; i++ {
		padding += "a"
	}
	base = []byte(makeBase(padding))

	// content := strings.Replace(content, "\n", "<br>", -1)
	stringLines := strings.SplitAfter(content, "\n")
	// stringLines := strings.SplitAfter(content, "<br>")
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
