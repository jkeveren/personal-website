package main

import (
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
	images := map[string]string{
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
	for name := range images {
		f[name] = newMockFile([]byte(name))
	}
	imageCount := len(images)

	g, err := newGalleryHandler(f)
	if err != nil {
		t.Fatal(err)
	}

	t.Run("sortedImagesLength", func(t *testing.T) {
		got := len(g.sortedImages)
		want := len(images)
		if got != want {
			t.Fatalf("Want: %d, Got: %d", want, got)
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

	t.Run("redirectHF", func(t *testing.T) {
		request, err := http.NewRequest("GET", "/galleryFirst", nil)
		if err != nil {
			t.Fatal(err)
		}
		recorder := httptest.NewRecorder()
		g.redirectHF(recorder, request)

		t.Run("statusCode", func(t *testing.T) {
			want := 307
			got := recorder.Code
			if got != want {
				t.Fatalf("Want: %d, Got: %d", want, got)
			}
		})

		t.Run("locationHeader", func(t *testing.T) {
			want := "/gallery/" + g.sortedImages[0]
			got := recorder.HeaderMap.Get("Location")
			if got != want {
				t.Fatalf("Want: %s, Got: %s", want, got)
			}
		})
	})

	t.Run("pageHF", func(t *testing.T) {
		request, err := http.NewRequest("GET", "/gallery/"+g.sortedImages[0], nil)
		if err != nil {
			t.Fatal(err)
		}
		recorder := httptest.NewRecorder()
		g.pageHF(recorder, request)

		t.Run("Content", func(t *testing.T) {
			want := "<!-- gallery72yr98mj -->"
			got := string(recorder.Body.Bytes())
			if got[:len(want)] != want {
				t.Fatalf("Want: %s..., Got: %s", want, got[:len(want)])
			}
		})
	})

	t.Run("imageHF", func(t *testing.T) {
		request, err := http.NewRequest("GET", "/gallery/image/image.jpg", nil)
		if err != nil {
			t.Fatal(err)
		}
		recorder := httptest.NewRecorder()
		g.imageHF(recorder, request)

		t.Run("content", func(t *testing.T) {
			want := "image.jpg"
			got := string(recorder.Body.Bytes())
			if want != got {
				t.Fatalf("Want: %s, Got: %s", want, got)
			}
		})

		t.Run("headers", func(t *testing.T) {
			for name, mime := range images {
				request, err := http.NewRequest("GET", "/"+name, nil)
				if err != nil {
					t.Fatal(err)
				}
				recorder := httptest.NewRecorder()
				g.imageHF(recorder, request)

				t.Run("Content-Type/"+name, func(t *testing.T) {
					got := recorder.HeaderMap.Get("Content-Type")
					if got != mime {
						t.Fatalf("Want: %s, Got: %s", mime, got)
					}
				})

				t.Run("Last-Modified/"+name, func(t *testing.T) {
					file, err := f.Open(name)
					if err != nil {
						t.Fatal(err)
					}
					stat, err := file.Stat()
					if err != nil {
						t.Fatal(err)
					}
					// Last-Modified header is always in GMT
					l, err := time.LoadLocation("GMT")
					if err != nil {
						t.Fatal(err)
					}
					want := stat.ModTime().In(l).Format(time.RFC1123)
					got := recorder.HeaderMap.Get("Last-Modified")
					if want != got {
						t.Fatalf("Want: %s, Got: %s", want, got)
					}
				})

				var sortedIndex int
				for i, sortedImage := range g.sortedImages {
					if sortedImage == name {
						sortedIndex = i
						break
					}
				}
				if sortedIndex < imageCount-1 {
					t.Run("Next/"+name, func(t *testing.T) {
						got := recorder.HeaderMap.Get("Next")
						want := g.sortedImages[sortedIndex+1]
						if want != got {
							t.Fatalf("Want: %s, Got: %s", want, got)
						}
					})
				}
				if sortedIndex > 0 {
					t.Run("Previous/"+name, func(t *testing.T) {
						got := recorder.HeaderMap.Get("Previous")
						want := g.sortedImages[sortedIndex-1]
						if want != got {
							t.Fatalf("Want: %s, Got: %s", want, got)
						}
					})
				}
			}

			t.Run("Cache", func(t *testing.T) {
				want := "public, max-age=3600, no-transform"
				got := recorder.HeaderMap.Get("Cache-Control")
				if got != want {
					t.Fatalf("Want: %s, Got: %s", want, got)
				}
			})

			t.Run("not-found", func(t *testing.T) {
				request, err := http.NewRequest("GET", "/gallery/image/yn8g7", nil)
				if err != nil {
					t.Fatal(err)
				}
				recorder := httptest.NewRecorder()
				g.imageHF(recorder, request)

				want := 404
				got := recorder.Code
				if got != want {
					t.Fatalf("Want: %d, Got: %d", want, got)
				}
			})
		})
	})
}
