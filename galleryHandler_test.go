package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGallery(t *testing.T) {
	g := newGalleryHandler()
	t.Run("", func(t *testing.T) {
		request, err := http.NewRequest("GET", "/gallery", nil)
		if err != nil {
			t.Fatal(err)
		}
		recorder := httptest.NewRecorder()
		g.ServeHTTP(recorder, request)
		want := "<!-- q98ny7g0sk -->"
		got := string(recorder.Body.Bytes())
		if got[:len(want)] != want {
			t.Fatalf("Want: %s..., Got: %s", want, got)
		}
	})
}
