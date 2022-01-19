package main

import (
	"errors"
	"io"
	"io/fs"
	"net/http"
	"path"
	"sort"
)

type galleryHandler struct {
	html         []byte
	fs           fs.FS
	sortedImages []string
}

var noImagesError = errors.New("No images were found in fs")

func newGalleryHandler(f fs.FS) (galleryHandler, error) {
	var nie error
	nieString := ""

	g := galleryHandler{
		fs: f,
	}

	// sortedImages
	firstImage := ""
	dirEntries, err := fs.ReadDir(g.fs, ".")
	if err != nil || len(dirEntries) == 0 {
		nie = noImagesError
		nieString = nie.Error()
	} else {
		// sort images to inverse chronological order
		sort.Slice(dirEntries, func(i, j int) bool {
			iInfo, err := dirEntries[i].Info()
			if err != nil {
				panic(err)
			}
			jInfo, err := dirEntries[j].Info()
			if err != nil {
				panic(err)
			}
			return iInfo.ModTime().After(jInfo.ModTime())
		})
		for _, dirEntry := range dirEntries {
			g.sortedImages = append(g.sortedImages, dirEntry.Name())
		}
		firstImage = g.sortedImages[0]
	}

	// html
	js, err := web.ReadFile("web/gallery.js")
	if err != nil {
		panic(err)
	}
	start := "<!-- gallery72yr98mj --><!DOCTYPE html>" + nieString + "<script type=module data-first-image=\"" + firstImage + "\" >\n" // contains identifier used in tests
	end := "</script>"
	g.html = append(append([]byte(start), js...), end...)

	return g, nie
}

func (g galleryHandler) indexHF(w http.ResponseWriter, r *http.Request) {
	w.Write(g.html)
}

func (g galleryHandler) imageHF(w http.ResponseWriter, r *http.Request) {
	image := path.Base(r.URL.Path)
	file, err := g.fs.Open(image)
	if err != nil {
		switch err {
		case fs.ErrInvalid, fs.ErrNotExist:
			w.WriteHeader(http.StatusBadRequest)
		default:
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}
	stat, err := file.Stat()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	h := w.Header()
	h.Add("Cache-Control", "public, max-age=3600, no-transform")
	for i, sortedImage := range g.sortedImages {
		if sortedImage == image {
			if i < len(g.sortedImages)-1 {
				h.Add("Next", g.sortedImages[i+1])
			}
			if i > 0 {
				h.Add("Previous", g.sortedImages[i-1])
			}
			break
		}
	}

	http.ServeContent(w, r, r.URL.Path, stat.ModTime(), file.(io.ReadSeeker))
}
