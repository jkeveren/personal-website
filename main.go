package main

import (
	"embed"
	"fmt"
	"net/http"
	"os"
)

//go:embed web
var web embed.FS

func main() {

	http.Handle("/", newHomeHandler())

	// port
	port := os.Getenv("PORT")
	if port == "" {
		port = "5000"
	}

	// HTTP server
	fmt.Println("Starting HTTP Server on Port " + port + ". Configure using PORT environment variable.")
	panic(http.ListenAndServe(":"+port, nil))
}
