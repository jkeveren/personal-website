package main

import (
	"fmt"
	"net/http"
	"os"
	"bytes"
	"time"
)

// Minimum number of bytes to send before text will render as it arrives.
var minBaseLength = 1023 // 1023 for Firefox, 511 for Edge, 2 for Chrome
var lineDelayMilliseconds = 20 // delay between sending each line

func main() {
	t := trickler{}

	// line delay
	t.lineDelay = time.Duration(lineDelayMilliseconds)
	content, err := os.ReadFile("./content.html")
	if err != nil {
		panic(err)
	}

	// HTML head
	emptyBaseLength := len(makeHead(""))
	paddingLength := minBaseLength - emptyBaseLength
	var padding string
	for i := 0; i < paddingLength; i++ {
		padding += "a"
	}
	t.head = []byte(makeHead(padding))

	// main content lines
	stringLines := bytes.SplitAfter(content, []byte("\n"))
	for _, stringLine := range stringLines {
		t.lines = append(t.lines, []byte(stringLine))
	}

	// port
	port := os.Getenv("PORT")
	if port == "" {
		port = "5000"
	}
	fmt.Println("Starting HTTP Server on Port " + port + ". Configure using PORT environment variable.")
	panic(http.ListenAndServe(":"+port, t))
}

// Returns HTML head with optional padding
func makeHead(padding string) string {
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
