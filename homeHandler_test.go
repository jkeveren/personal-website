package main

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
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
					t.Errorf("Want %d, got %d", paddingLength, diff)
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
						t.Errorf("Want %d, got %d", targetLength, length)
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
						t.Errorf("Want %d, got %d", targetLength, length)
					}
				})
			}
		})

		t.Run("html", func(t *testing.T) {
			if !bytes.Contains(h.makeHead(0), []byte("<html")) {
				// Testing for html tag is a simple way to ensure that actual content is sent not just padding.
				t.Error("HTML head does not contain html tag")
			}
		})
	})

	t.Run("HTTP", func(t *testing.T) {
		c, cancel := context.WithCancel(context.Background())
		request, err := http.NewRequestWithContext(c, "GET", "/", nil)
		if err != nil {
			t.Fatal(err)
		}
		recorder := newStreamRecorder()
		go func() {
			for {
				body := <-recorder.c
				if bytes.Contains(body, []byte("<html")) {
					cancel()
					return
				}
			}
			// Testing for html tag is a simple way to ensure that actual content is sent not just padding.
			t.Error("Response does not contain html tag")
		}()
		h.ServeHTTP(recorder, request)
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
