package main

import (
	"embed"
	"flag"
	"fmt"
	"net/http"
	// "os"
)

//go:embed web
var web embed.FS

func main() {

	http.Handle("/", newHomeHandler())

	address := flag.String("a", "0.0.0.0:8000", "Address to bind to. Includes port.")
	flag.Parse()

	// HTTP server
	fmt.Println("Starting HTTP Server on address " + *address + ". Configure using \"a\" flag.")
	panic(http.ListenAndServe(*address, nil))
}
