package main

import (
	"embed"
	"flag"
	"fmt"
	"net/http"
	"os"
)

//go:embed web
var web embed.FS

func main() {

	// Address
	address := flag.String("a", "0.0.0.0:8000", "Address to bind to")
	flag.Parse()

	// Routes
	http.Handle("/", newHomeHandler())
	http.Handle("/favicon.ico", http.NotFoundHandler())
	gh, err := newGalleryHandler(os.DirFS("/usr/share/personal-website/gallery"))
	if err != nil {
		if err == noImagesError {
			fmt.Println(err.Error())
		} else {
			panic(err)
		}
	} else {
		http.HandleFunc("/gallery/", gh.imageHF)
	}
	http.HandleFunc("/gallery", gh.indexHF)

	// HTTP server
	fmt.Println("Starting HTTP Server on address " + *address + ". Configure using the -a flag.")
	panic(http.ListenAndServe(*address, nil))
}
