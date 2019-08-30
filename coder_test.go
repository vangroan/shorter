package main


import (
	"testing"
)

func TestEncoder(t *testing.T) {
	num := 123
	short := EncodeNumber(num)

	if short != "abc" {
		t.Errorf("")
	}
}
