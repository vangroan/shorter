package main

import (
	"testing"
)

type testRecord struct {
	url         string
	expectError bool
	actual      error
	message     string
}

func TestValidation(t *testing.T) {
	data := [5]testRecord{
		testRecord{
			url:         "javascript:alert(1);",
			expectError: true,
			actual:      nil,
			message:     "Validiation did not catch injected javascript.",
		},
		testRecord{
			url:         "data:text/html,<script>alert('hi');</script>",
			expectError: true,
			actual:      nil,
			message:     "Validiation did not catch injected javascript.",
		},
		testRecord{
			url:         "http://www.test.com",
			expectError: false,
			actual:      nil,
			message:     "Validiation blocked valid URL.",
		},
		testRecord{
			url:         "https://www.test.com",
			expectError: false,
			actual:      nil,
			message:     "Validiation blocked valid URL.",
		},
		testRecord{
			url:         "www.test.com",
			expectError: true,
			actual:      nil,
			message:     "Validiation did not catch URL without schema.",
		},
	}

	for i := 0; i < len(data); i++ {
		data[i].actual = validateURI(data[i].url)
	}

	for _, record := range data {
		if record.expectError {
			if record.actual == nil {
				t.Error(record.message, "Expected Error:", record.expectError, "Actual: nil")
			}
		} else {
			if record.actual != nil {
				t.Error(record.message, "Expected Error:", record.expectError, "Actual:", record.actual.Error())
			}
		}
	}
}
