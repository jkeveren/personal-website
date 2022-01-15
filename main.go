package main

import (
	"fmt"
	"net/http"
	"os"
)

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
