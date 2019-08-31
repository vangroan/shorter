package main

import (
	"testing"
)

type data struct {
	expectedNum uint64
	expectedStr string
	actualNum   uint64
	actualStr   string
}

func TestReverse(t *testing.T) {
	s := "1234567890"
	reversed := reverse(s)

	if reversed != "0987654321" {
		t.Error("String was not reversed properly. Expected:", s, "Actual", reversed)
	}
}

func TestEncoder(t *testing.T) {
	nums := [4]data{
		data{
			expectedNum: 0,
			expectedStr: "a",
			actualStr:   "",
		},
		data{
			expectedNum: 10,
			expectedStr: "k",
			actualStr:   "",
		},
		data{
			expectedNum: 64,
			expectedStr: "ba",
			actualStr:   "",
		},
		data{
			expectedNum: 4500,
			expectedStr: "bgu",
			actualStr:   "",
		},
	}

	for i := 0; i < len(nums); i++ {
		nums[i].actualStr = EncodeNumber(nums[i].expectedNum)
	}

	for _, num := range nums {
		if num.expectedStr != num.actualStr {
			t.Error("Expected:", num.expectedStr, "Actual:", num.actualStr)
		}
	}
}

func TestDecoder(t *testing.T) {
	nums := [4]data{
		data{
			expectedNum: 0,
			expectedStr: "a",
		},
		data{
			expectedNum: 10,
			expectedStr: "k",
		},
		data{
			expectedNum: 64,
			expectedStr: "ba",
		},
		data{
			expectedNum: 4500,
			expectedStr: "bgu",
		},
	}

	for i := 0; i < len(nums); i++ {
		nums[i].actualNum = DecodeNumber(nums[i].expectedStr)
	}

	for _, num := range nums {
		if num.expectedNum != num.actualNum {
			t.Error("Expected:", num.expectedNum, "Actual:", num.actualNum)
		}
	}
}
