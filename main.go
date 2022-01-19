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
	galleryLocation := flag.String("g", "/usr/share/personal-website/gallery", "Gallery image location")
	flag.Parse()

	// Routes
	http.Handle("/", newHomeHandler())
	http.Handle("/favicon.ico", http.NotFoundHandler())
	gh, err := newGalleryHandler(os.DirFS(*galleryLocation))
	if err != nil {
		if err == noImagesError {
			fmt.Println(err.Error())
		} else {
			panic(err)
		}
	} else {
		http.HandleFunc("/galleryImage/", gh.imageHF)
	}
	http.HandleFunc("/gallery", gh.redirectHF)
	http.HandleFunc("/gallery/", gh.pageHF)

	// HTTP server
	fmt.Println("Starting HTTP Server on address " + *address + ". Configure using the -a flag.")
	panic(http.ListenAndServe(*address, nil))
}
