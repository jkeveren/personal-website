package main

import (
	"bytes"
	"math"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"testing"
	"testing/fstest"
	"time"
)

func newMockFile(content []byte) *fstest.MapFile {
	// randomize the file mod times so sorting must occur
	modTime := time.UnixMilli(rand.Int63())
	return &fstest.MapFile{
		Data:    content,
		ModTime: modTime,
	}
}

func TestGalleryHandler(t *testing.T) {

	rand.Seed(time.Now().UnixNano())

	// generate mock files
	f := fstest.MapFS{}
	imagesWithMime := map[string]string{
		// "file name": "expected mime type",
		"image.jpg":  "image/jpeg",
		"image.JPG":  "image/jpeg",
		"image.jpeg": "image/jpeg",
		"image.JPEG": "image/jpeg",
		"image.png":  "image/png",
		"image.PNG":  "image/png",
		"image.mp4":  "video/mp4",
		"image.MP4":  "video/mp4",
	}
	for name := range imagesWithMime {
		f[name] = newMockFile([]byte(name))
	}

	g, err := newGalleryHandler(f)
	if err != nil {
		t.Fatal(err)
	}

	t.Run("sortedImagesLength", func(t *testing.T) {
		got := len(g.sortedImages)
		imageCount := len(imagesWithMime)
		if got != imageCount {
			t.Fatalf("Want %d, Got %d", imageCount, got)
		}
	})

	t.Run("sortedImagesOrder", func(t *testing.T) {
		lastTime := time.UnixMilli(math.MaxInt64)
		for _, imageName := range g.sortedImages {
			imageInfo, err := f.Stat(imageName)
			if err != nil {
				t.Fatal(err)
			}
			mt := imageInfo.ModTime()
			t.Log(mt.After(lastTime))
			if mt.After(lastTime) {
				for _, imageName := range g.sortedImages {
					imageInfo, err := f.Stat(imageName)
					if err != nil {
						t.Fatal(err)
					}
					t.Log(imageInfo.ModTime().Unix())
				}
				t.Fatal("Not in order of latest first")
			}
			lastTime = mt
		}
	})

	t.Run("indexHF", func(t *testing.T) {
		request, err := http.NewRequest("GET", "/gallery", nil)
		if err != nil {
			t.Fatal(err)
		}
		recorder := httptest.NewRecorder()
		g.indexHF(recorder, request)

		t.Run("HTMLContent", func(t *testing.T) {
			want := "<!-- gallery72yr98mj -->"
			got := string(recorder.Body.Bytes())
			if got[:len(want)] != want {
				t.Fatalf("Want: %s..., Got: %s", want, got[:len(want)])
			}
		})

		t.Run("firstImage", func(t *testing.T) {
			want := []byte("data-first-image=\"" + g.sortedImages[0] + "\"")
			got := recorder.Body.Bytes()
			if !bytes.Contains(got, want) {
				n := 150
				t.Fatalf("Could not find \"%s\" in body.\nFirst %d characters of body: \"%s\"", want, n, string(got[:n]))
			}
		})
	})

	t.Run("imageHF", func(t *testing.T) {
		request, err := http.NewRequest("GET", "/image.jpg", nil)
		if err != nil {
			t.Fatal(err)
		}
		recorder := httptest.NewRecorder()
		g.imageHF(recorder, request)

		t.Run("content", func(t *testing.T) {
			want := "image.jpg"
			got := string(recorder.Body.Bytes())
			if want != got {
				t.Fatalf("Want %s, Got %s", want, got)
			}
		})

		t.Run("headers", func(t *testing.T) {
			t.Run("Content-Type", func(t *testing.T) {
				for name, want := range imagesWithMime {
					t.Run(name, func(t *testing.T) {
						request, err := http.NewRequest("GET", "/"+name, nil)
						if err != nil {
							t.Fatal(err)
						}
						recorder := httptest.NewRecorder()
						g.imageHF(recorder, request)
						got := recorder.HeaderMap.Get("Content-Type")
						if got != want {
							t.Fatalf("Want %s, Got %s", want, got)
						}
					})
				}
			})

			t.Run("Cache", func(t *testing.T) {
				want := "public, max-age=3600, no-transform"
				got := recorder.HeaderMap.Get("Cache-Control")
				if got != want {
					t.Fatalf("Want %s, Got %s", want, got)
				}
			})
		})
	})
}
