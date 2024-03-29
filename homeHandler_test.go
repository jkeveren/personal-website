package main

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"
)

func TestHomeHandler(t *testing.T) {
	h := newHomeHandler()

	t.Run("makeParametrichead", func(t *testing.T) {
		shortLength := len(h.makeParametricHead(0))
		tests := []int{
			0,
			50,
			1000,
		}
		for _, paddingLength := range tests {
			t.Run(strconv.Itoa(paddingLength), func(t *testing.T) {
				length := len(h.makeParametricHead(paddingLength))
				diff := length - shortLength
				if diff != paddingLength {
					t.Fatalf("Want %d, got %d", paddingLength, diff)
				}
			})
		}
	})

	t.Run("makeHead", func(t *testing.T) {
		t.Run("long", func(t *testing.T) {
			tests := []int{
				// numbers must be MORE than than minimum head content length
				10000,
				11000,
				12000,
			}
			for _, targetLength := range tests {
				t.Run(strconv.Itoa(targetLength), func(t *testing.T) {
					length := len(h.makeHead(targetLength))
					if length != targetLength {
						t.Fatalf("Want %d, got %d", targetLength, length)
					}
				})
			}
		})

		t.Run("short", func(t *testing.T) {
			tests := []int{
				// numbers must be LESS than than minimum head content length
				0,
				20,
				50,
			}
			for _, targetLength := range tests {
				t.Run(strconv.Itoa(targetLength), func(t *testing.T) {
					length := len(h.makeHead(targetLength))
					if length < targetLength {
						t.Fatalf("Want %d, got %d", targetLength, length)
					}
				})
			}
		})

		t.Run("html", func(t *testing.T) {
			// Testing for html tag is a simple way to ensure that actual content is sent not just padding.
			if !bytes.Contains(h.makeHead(0), []byte("<html")) {
				t.Fatal("HTML head does not contain html tag")
			}
		})
	})

	t.Run("ServeHTTP", func(t *testing.T) {
		c, cancel := context.WithCancel(context.Background())
		t.Cleanup(func() {
			cancel()
		})
		request, err := http.NewRequestWithContext(c, "GET", "/", nil)
		if err != nil {
			t.Fatal(err)
		}
		recorder := newStreamRecorder()
		go h.ServeHTTP(recorder, request)

		t.Run("body", func(t *testing.T) {
			want := "<!-- homem98y2r8 -->"
			got := <-recorder.c
			if !bytes.Contains(got, []byte(want)) {
				t.Fatalf("Want %s, Got %s...", want, string(got[:100]))
				return
			}
		})

		t.Run("cacheHeader", func(t *testing.T) {
			want := "no-store"
			got := recorder.HeaderMap.Get("Cache-Control")
			if got != want {
				t.Fatalf("Want %s, Got %s", want, got)
			}
		})

		t.Run("redirect", func(t *testing.T) {
			request, err := http.NewRequest("GET", "/yn8ge8r7", nil)
			if err != nil {
				t.Fatal(err)
			}
			recorder := httptest.NewRecorder()
			h.ServeHTTP(recorder, request)

			t.Run("statusCode", func(t *testing.T) {
				want := 307
				got := recorder.Code
				if got != want {
					t.Fatalf("Want %d, Got %d", want, got)
				}
			})

			t.Run("locationHeader", func(t *testing.T) {
				want := "/"
				got := recorder.HeaderMap.Get("Location")
				if got != want {
					t.Fatalf("Want %s, Got %s", want, got)
				}
			})
		})
	})

	t.Run("newHomeHandler", func(t *testing.T) {
		tests := []struct {
			name string
			want interface{}
			got  interface{}
		}{
			{"headLength", 1023, len(h.head)},
			{"lineDelay", time.Duration(100) * time.Millisecond, h.lineDelay},
		}
		for _, test := range tests {
			t.Run(test.name, func(t *testing.T) {
				if test.got != test.want {
					t.Fatalf("Want %d, Got %d", test.want, test.got)
				}
			})
		}

		t.Run("line count", func(t *testing.T) {
			min := 10
			got := len(h.lines)
			if got < min {
				t.Fatalf("Want >%d, Got %d", min, got)
			}
		})
	})
}

type streamRecorder struct {
	*httptest.ResponseRecorder
	c chan []byte
}

func newStreamRecorder() streamRecorder {
	return streamRecorder{
		httptest.NewRecorder(),
		make(chan []byte),
	}
}

func (r streamRecorder) Write(buf []byte) (int, error) {
	r.c <- buf
	return len(buf), nil
}
