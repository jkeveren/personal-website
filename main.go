package main

import (
	"embed"
	"flag"
	"fmt"
	// "io/fs"
	"net/http"
)

//go:embed web
var web embed.FS

func main() {

	// Address
	address := flag.String("a", "0.0.0.0:8000", "Address to bind to")
	flag.Parse()

	// Routes
	http.Handle("/", newHomeHandler())
	http.Handle("/gallery", newGalleryHandler())
	prefix := "/static/"
	http.Handle(prefix, http.StripPrefix(prefix, http.FileServer(http.Dir("web/static"))))

	// HTTP server
	fmt.Println("Starting HTTP Server on address " + *address + ". Configure using the -a flag.")
	panic(http.ListenAndServe(*address, nil))
}
