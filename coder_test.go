package main

import (
	"testing"
)

type data struct {
	num      uint64
	expected string
	actual   string
}

func TestReverse(t *testing.T) {
	s := "1234567890"
	reversed := reverse(s)

	if reversed != "0987654321" {
		t.Error("String was not reversed properly. Expected:", s, "Actual", reversed)
	}
}

func TestEncoder(t *testing.T) {
	nums := [3]data{
		data{
			num:      10,
			expected: "k",
			actual:   "",
		},
		data{
			num:      64,
			expected: "ba",
			actual:   "",
		},
		data{
			num:      4500,
			expected: "bgu",
			actual:   "",
		},
	}

	for i := 0; i < len(nums); i++ {
		nums[i].actual = EncodeNumber(nums[i].num)
	}

	for _, num := range nums {
		if num.expected != num.actual {
			t.Error("Expected:", num.expected, "Actual:", num.actual)
		}
	}
}
