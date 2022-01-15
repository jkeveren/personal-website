package main

import (
	"strconv"
	"testing"
)

func TestHomeHandler(t *testing.T) {
	tr := homeHandler{}

	t.Run("makeParametrichead", func(t *testing.T) {
		shortLength := len(tr.makeParametricHead(0))
		tests := []int{
			0,
			50,
			1000,
		}
		for _, paddingLength := range tests {
			t.Run(strconv.Itoa(paddingLength), func(t *testing.T) {
				length := len(tr.makeParametricHead(paddingLength))
				diff := length - shortLength
				if diff != paddingLength {
					t.Errorf("got %d, want %d", diff, paddingLength)
				}
			})
		}
	})

	t.Run("makeHeadLong", func(t *testing.T) {
		tests := []int{
			// numbers must be MORE than than minimum head content length
			10000,
			11000,
			12000,
		}
		for _, targetLength := range tests {
			t.Run(strconv.Itoa(targetLength), func(t *testing.T) {
				length := len(tr.makeHead(targetLength))
				if length != targetLength {
					t.Errorf("Want %d, got %d", targetLength, length)
				}
			})
		}
	})

	t.Run("makeHeadShort", func(t *testing.T) {
		tests := []int{
			// numbers must be LESS than than minimum head content length
			0,
			20,
			50,
		}
		for _, targetLength := range tests {
			t.Run(strconv.Itoa(targetLength), func(t *testing.T) {
				length := len(tr.makeHead(targetLength))
				if length < targetLength {
					t.Errorf("Want %d, got %d", targetLength, length)
				}
			})
		}
	})
}
