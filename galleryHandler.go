package main

import (
	"net/http"
)

type galleryHandler struct {
	html []byte
}

func newGalleryHandler() galleryHandler {
	js, err := web.ReadFile("web/gallery.js")
	if err != nil {
		panic(err)
	}
	return galleryHandler{
		// contains identifier used in tests
		append(append([]byte("<!-- q98ny7g0sk --><!DOCTYPE html><script type=module>\n"), js...), "</script>"...),
	}
}

func (g galleryHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Write(g.html)
}
